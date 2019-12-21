package crypter

import (
	"bytes"
	"github.com/Bigyin1/GoMobileBackend/pkg/infrastructure"
	"github.com/palantir/stacktrace"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestCrypterGetFile(t *testing.T) {
	fileRepo := infrastructure.NewInFsFileStorage("./testdata")
	service := NewCrypterService(fileRepo, "", GetRandomEncrKey)

	testFilesMap := map[string][]string{"./testdata/encrTest.txt": {"encrTest.txt.encr","QdCk6jOBjdUsus5Z"} ,
		"./testdata/gopher.png": {"gopher.png.encr", "NUzwzHGFKubFxL0a"}}

	for initPath, encrPath := range testFilesMap {
		decrypted, err := service.DecryptFileAndGet(encrPath[0], encrPath[1])
		assert.Nil(t, err)
		exp, err := ioutil.ReadFile(initPath)
		assert.Nil(t, err)
		assert.True(t, bytes.Equal(exp, decrypted))
	}
}

func TestCrypterGetFileWrongKey(t *testing.T) {
	fileRepo := infrastructure.NewInFsFileStorage("./testdata")
	service := NewCrypterService(fileRepo, "", GetRandomEncrKey)

	encrFiles := [][]string{{"encrTest.txt.encr","QdCk6jOBjdUsus5y"},
		{"gopher.png.encr", "NUzwzHhFKubFxL0a"}}

	for _, encrFile := range encrFiles {
		_, err := service.DecryptFileAndGet(encrFile[0], encrFile[1])
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
	inputFilesMap := make(map[string][]byte)
	for _, fPath := range inputs {
		data, err := ioutil.ReadFile(fPath)
		assert.Nil(t, err)
		inputFilesMap[fPath] = data
	}

	mapping := service.EncryptAndSaveFiles(inputFilesMap)
	for _, storedFile := range mapping {
		expected := inputFilesMap[storedFile.Name]
		got, err := service.DecryptFileAndGet(storedFile.GetFid(), testKey)
		assert.Nil(t, err)
		assert.True(t, bytes.Equal(expected, got))
	}
}
