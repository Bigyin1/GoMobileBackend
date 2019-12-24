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
	testFiles := []string{"./testdata/encrTest.txt", "./testdata/gopher.png"}
	for _, file := range testFiles {
		f, err := os.Open(file)
		assert.Nil(t, err)
		initData, err := ioutil.ReadAll(f)
		initDataReader := bytes.NewReader(initData)
		var encrData bytes.Buffer
		var decrData bytes.Buffer
		assert.Nil(t, err)

		err = encrypt(initDataReader, &encrData, key)
		assert.Nil(t, err)
		err = decrypt(&encrData, &decrData, key)
		assert.Nil(t, err)
		assert.True(t, bytes.Equal(initData, decrData.Bytes()))
	}
}
