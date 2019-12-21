package crypter

import (
	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"github.com/satori/go.uuid"
	"net/url"
	"path"
)

type Service struct {
	fileRepository infrastructure.FileRepository
	fileUriPrefix  string
}

func NewCrypterService(fileRepository infrastructure.FileRepository, prefix string) *Service {
	return &Service{fileRepository: fileRepository, fileUriPrefix: prefix}
}

func (s *Service) combineFileURL(uuid, key string) string {
	u, _ := url.Parse(s.fileUriPrefix)
	u.Path = path.Join(u.Path, uuid)
	q := u.Query()
	q.Set("key", key)
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *Service) encryptAndSaveFile(data []byte, key []byte) {

}

func (s *Service) EncryptAndSaveFiles(files map[string][]byte) Mapping {
	key := getRandomEncrKey()
	mapping := newFilesMapping()
	for fileName, fileData := range files {
		fid := uuid.NewV4().String()
		encryptedFileData, err := encrypt(fileData, key)
		if err != nil {
			mapping.Add(fileName, err.Error())
			continue
		}
		err = s.fileRepository.StoreFileByID(fid, encryptedFileData)
		if err != nil {
			mapping.Add(fileName, err.Error())
			continue
		}
		mapping.Add(fileName, s.combineFileURL(fid, string(key)))
	}
	return mapping
}

func (s *Service) DecryptFile(fileId, key string) ([]byte, error) {
	f, err := s.fileRepository.FindFileByID(fileId)
	if err != nil {
		return nil, err
	}
	res, err := decrypt(f, []byte(key))
	if err != nil {
		return nil, err
	}
	return res, nil
}
