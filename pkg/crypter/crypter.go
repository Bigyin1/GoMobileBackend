package crypter

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"github.com/palantir/stacktrace"
	"github.com/satori/go.uuid"
	"net/url"
	"path"
)

type Service struct {
	fileRepository infrastructure.FileRepository
	fileUriPrefix  string
	keyGenerator func() []byte
}

type InputFiles map[string][]byte

func NewCrypterService(fileRepository infrastructure.FileRepository, urlPrefix string, keyGen func() []byte) *Service {
	return &Service{fileRepository: fileRepository, fileUriPrefix: urlPrefix, keyGenerator:keyGen}
}

func (s *Service) combineFileURL(uuid, key string) string {
	u, _ := url.Parse(s.fileUriPrefix)
	u.Path = path.Join(u.Path, uuid)
	q := u.Query()
	q.Set("key", key)
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *Service) encryptAndSaveFile(data []byte, key []byte, fid string) {

}

func (s *Service) EncryptAndSaveFiles(files InputFiles) Mapping {
	mapping := newFilesMapping()
	for fileName, fileData := range files {
		key := s.keyGenerator()
		fid := uuid.NewV4().String()
		encryptedFileData, err := encrypt(fileData, key)
		if err != nil {
			mapping.AddError(fileName, stacktrace.RootCause(err).Error(), fid)
			continue
		}
		err = s.fileRepository.StoreFileByID(fid, encryptedFileData)
		if err != nil {
			mapping.AddError(fileName, stacktrace.RootCause(err).Error(), fid)
			continue
		}
		mapping.Add(fileName, s.combineFileURL(fid, string(key)), fid)
	}
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
