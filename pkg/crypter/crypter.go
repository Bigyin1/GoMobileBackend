package crypter

import (
	"io"
	"net/url"
	"path"
	"sync"

	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"github.com/palantir/stacktrace"
	uuid "github.com/satori/go.uuid"
)

type Service struct {
	fileRepository infrastructure.FileRepository
	fileUriPrefix  string
	keyGenerator   func() []byte
}

type InputFiles map[string]io.Reader

func NewCrypterService(fileRepository infrastructure.FileRepository, urlPrefix string, keyGen func() []byte) *Service {
	return &Service{fileRepository: fileRepository, fileUriPrefix: urlPrefix, keyGenerator: keyGen}
}

func (s *Service) combineFileURL(uuid, key string) string {
	u, _ := url.Parse(s.fileUriPrefix)
	u.Path = path.Join(u.Path, uuid)
	q := u.Query()
	q.Set("key", key)
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *Service) EncryptAndSaveFile(fileReader io.Reader, fileName string, mapping *Mapping) {
	key := s.keyGenerator()
	fid := uuid.NewV4().String()
	fileWriter, err := s.fileRepository.GetFileWriterByID(fid)
	if err != nil {
		mapping.AddError(fileName, "Error while creating file", fid)
		return
	}
	defer fileWriter.Close()
	err = encrypt(fileReader, fileWriter, key)
	if err != nil {
		mapping.AddError(fileName, "Error while encrypting file", fid)
		return
	}
	mapping.Add(fileName, s.combineFileURL(fid, string(key)), fid)
}

func (s *Service) EncryptAndSaveFiles(files InputFiles) Mapping {
	mapping := NewFilesMapping()
	wg := &sync.WaitGroup{}
	wg.Add(len(files))
	for fileName, fileData := range files {
		go func(fileData io.Reader, fileName string) {
			s.EncryptAndSaveFile(fileData, fileName, &mapping)
			wg.Done()
		}(fileData, fileName)
	}
	wg.Wait()
	return mapping
}

func (s *Service) DecryptFileAndGet(fileId, key string, dest io.Writer) error {
	fileReader, err := s.fileRepository.FindFileReaderByID(fileId)
	if err != nil {
		return stacktrace.Propagate(err, "failed to find file in repo")
	}
	defer fileReader.Close()
	err = decrypt(fileReader, dest, []byte(key))
	if err != nil {
		return stacktrace.PropagateWithCode(err, ErrWrongKey, "failed to decrypt file")
	}
	return nil
}
