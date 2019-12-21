package crypter

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte(GetRandomEncrKey())
	testFiles := []string{ "./testdata/encrTest.txt", "./testdata/gopher.png"}
	for _, file := range testFiles {
		f, err := os.Open(file)
		assert.Nil(t, err)
		initData, err := ioutil.ReadAll(f)
		assert.Nil(t, err)

		encrData, err := encrypt(initData, key)
		assert.Nil(t, err)
		decrData, err := decrypt(encrData, key)
		assert.Nil(t, err)
		assert.True(t, bytes.Equal(initData, decrData))
	}
}
