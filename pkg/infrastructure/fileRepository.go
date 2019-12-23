package infrastructure

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/palantir/stacktrace"
)

type FileRepository interface {
	FindFileReaderByID(fid string) (io.ReadCloser, error)
	GetFileWriterByID(fid string) (io.WriteCloser, error)
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

func (s *InFsFileStorage) FindFileReaderByID(fid string) (io.ReadCloser, error) {
	file, err := os.Open(filepath.Join(s.storagePath, fid))
	if os.IsNotExist(err) {
		return nil, stacktrace.NewMessageWithCode(ErrFileNotFound, fmt.Sprintf("File with id: %s not exists", fid))
	}
	tst, _ := file.Stat()
	fmt.Println(tst.Size())
	if err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "")
	}
	return file, nil
}

func (s *InFsFileStorage) GetFileWriterByID(fid string) (io.WriteCloser, error) {
	path := filepath.Join(s.storagePath, fid)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "File with uuid already exists")
	}
	newFile, err := os.Create(path)
	if err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "Failed to create file")
	}
	return newFile, nil
}
