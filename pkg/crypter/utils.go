package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
)

func decrypt(source io.Reader, out io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: source}
	// Copy the input to the output stream, decrypting as we go.
	if _, err := io.Copy(out, reader); err != nil {
		return stacktrace.Propagate(err, "decrypt err")
	}
	return nil
}

func encrypt(source io.Reader, out io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	writer := &cipher.StreamWriter{S: stream, W: out}
	// Copy the input to the output buffer, encrypting as we go.
	if _, err := io.Copy(writer, source); err != nil {
		return stacktrace.Propagate(err, "encrypt err")
	}
	return nil
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
