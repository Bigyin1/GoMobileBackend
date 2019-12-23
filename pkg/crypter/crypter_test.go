package crypter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"github.com/palantir/stacktrace"
	"github.com/stretchr/testify/assert"
)

func TestCrypterGetFile(t *testing.T) {
	fileRepo := infrastructure.NewInFsFileStorage("./testdata")
	service := NewCrypterService(fileRepo, "", GetRandomEncrKey)

	testFilesMap := map[string][]string{"./testdata/encrTest.txt": {"encrTest.txt.encr", "QdCk6jOBjdUsus5Z"},
		"./testdata/gopher.png": {"gopher.png.encr", "NUzwzHGFKubFxL0a"}}

	for initPath, encrPath := range testFilesMap {
		var decrypted bytes.Buffer
		err := service.DecryptFileAndGet(encrPath[0], encrPath[1], &decrypted)
		assert.Nil(t, err)
		exp, err := ioutil.ReadFile(initPath)
		assert.Nil(t, err)
		fmt.Println(len(decrypted.Bytes()), len(exp))
		assert.True(t, bytes.Equal(exp, decrypted.Bytes()))
	}
}

func TestCrypterGetFileWrongKey(t *testing.T) {
	fileRepo := infrastructure.NewInFsFileStorage("./testdata")
	service := NewCrypterService(fileRepo, "", GetRandomEncrKey)

	encrFiles := [][]string{{"encrTest.txt.encr", "QdCk6jOBjdUsus5y"},
		{"gopher.png.encr", "NUzwzHhFKubFxL0a"}}

	for _, encrFile := range encrFiles {
		var decr bytes.Buffer
		err := service.DecryptFileAndGet(encrFile[0], encrFile[1], &decr)
		assert.NotNil(t, err)
		assert.EqualValues(t, ErrWrongKey, stacktrace.GetCode(err))
	}
}

func TestCrypterSaveFile(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "encr_test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpdir)
	fileRepo := infrastructure.NewInFsFileStorage(tmpdir)
	testKey := "NUzwzHhFKubFxL0a"
	service := NewCrypterService(fileRepo, "", func() []byte {
		return []byte(testKey)
	})

	inputs := []string{"./testdata/encrTest.txt", "./testdata/gopher.png"}
	inputFilesMap := make(InputFiles)
	for _, fPath := range inputs {
		f, err := os.Open(fPath)
		assert.Nil(t, err)
		defer f.Close()
		inputFilesMap[fPath] = f
	}

	mapping := service.EncryptAndSaveFiles(inputFilesMap)
	for _, storedFile := range mapping {
		expected, err := ioutil.ReadAll(inputFilesMap[storedFile.Name])
		assert.Nil(t, err)
		var decr bytes.Buffer
		err = service.DecryptFileAndGet(storedFile.GetFid(), testKey, &decr)
		assert.Nil(t, err)

		assert.True(t, bytes.Equal(expected, decr.Bytes()))
	}
}
