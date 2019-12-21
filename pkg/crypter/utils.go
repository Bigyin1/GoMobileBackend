package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/palantir/stacktrace"
	"io"
	mrand "math/rand"
	"strings"
	"time"
)

func encrypt(text []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "problem with enc key")
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "")
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(text []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, stacktrace.PropagateWithCode(err, ErrWrongKey, "key is too small")
	}
	if len(text) < aes.BlockSize {
		return nil, stacktrace.PropagateWithCode(err, ErrUnexpected, "problem with enc key")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, stacktrace.NewMessageWithCode(ErrWrongKey, "Failed to decrypt file with key: %s", string(key))
	}
	return data, nil
}

func GetRandomEncrKey() []byte {
	mrand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[mrand.Intn(len(chars))])
	}
	return []byte(b.String())
}
