package crypter

import (
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

type InputFiles map[string][]byte

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

func (s *Service) encryptAndSaveFile(fileData []byte, fileName string, mapping *Mapping, wg *sync.WaitGroup) {
	defer wg.Done()
	key := s.keyGenerator()
	fid := uuid.NewV4().String()
	encryptedFileData, err := encrypt(fileData, key)
	if err != nil {
		mapping.AddError(fileName, "Error while encrypting file", fid)
		return
	}
	err = s.fileRepository.StoreFileByID(fid, encryptedFileData)
	if err != nil {
		mapping.AddError(fileName, "Error while storing file", fid)
		return
	}
	mapping.Add(fileName, s.combineFileURL(fid, string(key)), fid)
}

func (s *Service) EncryptAndSaveFiles(files InputFiles) Mapping {
	mapping := newFilesMapping()
	wg := sync.WaitGroup{}
	wg.Add(len(files))
	for fileName, fileData := range files {
		go s.encryptAndSaveFile(fileData, fileName, &mapping, &wg)
	}
	wg.Wait()
	return mapping
}

func (s *Service) DecryptFileAndGet(fileId, key string) ([]byte, error) {
	f, err := s.fileRepository.FindFileByID(fileId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find file in repo")
	}
	res, err := decrypt(f, []byte(key))
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to decrypt file")
	}
	return res, nil
}
