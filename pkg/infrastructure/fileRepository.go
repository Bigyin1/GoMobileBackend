package infrastructure

import (
	"fmt"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"log"
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
	err := createFileStorageDirectory(storagePath)
	if err != nil {
		log.Fatalln("Failed to create storage path", err)
	}
	return &InFsFileStorage{storagePath: storagePath}
}

func (s *InFsFileStorage) FindFileByID(fid string) ([]byte, error) {
	file, err := os.Open(filepath.Join(s.storagePath, fid))
	if os.IsNotExist(err) {
		return nil, stacktrace.NewMessageWithCode(ErrFileNotFound, fmt.Sprintf("File with id: %s not exists", fid))
	}
	if err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "")
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "")
	}
	return b, nil
}

func (s *InFsFileStorage) StoreFileByID(fid string, f []byte) error {
	path := filepath.Join(s.storagePath, fid)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return stacktrace.PropagateWithCode(err, ErrUnexpected, "File with uuid already exists")
	}
	newFile, _ := os.Create(path)
	defer newFile.Close()
	_, err := newFile.Write(f)
	if err != nil {
		return stacktrace.PropagateWithCode(err, ErrUnexpected, "")
	}
	return nil
}
