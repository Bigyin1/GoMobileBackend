package infrastructure

import (
	"errors"
	"fmt"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileRepository interface {
	FindFileByID(fid string) ([]byte, error)
	StoreFileByID(fid string, f []byte) error
}

type InFsFileStorage struct {
	storagePath string
}

func NewInFsFileStorage(storagePath string) *InFsFileStorage {
	return &InFsFileStorage{storagePath: storagePath}
}

func (s *InFsFileStorage) FindFileByID(fid string) ([]byte, error) {
	file, err := os.Open(filepath.Join(s.storagePath, fid))
	defer file.Close()
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File with id: %s not exists", fid))
	}
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *InFsFileStorage) StoreFileByID(fid string, f []byte) error {
	path := filepath.Join(s.storagePath, fid)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return err
	}
	newFile, _ := os.Create(path)
	defer newFile.Close()
	_, err := newFile.Write(f)
	if err != nil {
		return err
	}
	return nil
}
