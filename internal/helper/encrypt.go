package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/Luzifer/go-openssl"
	"io"
)

func EncryptMessage(inputMessage, inputKey string) (string, error) {
	byteMsg := []byte(inputMessage)
	key := []byte(inputKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(message, inputKey string) (string, error) {
	key := []byte(inputKey)
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}

func EncryptClientMessage(message, secret string) (string, error) {
	o := openssl.New()

	enc, err := o.EncryptBytes(secret, []byte(message))
	if err != nil {
		return "", err
	}

	return string(enc), nil
}

func DecryptClientMessage(encrypted, secret string) (string, error) {
	o := openssl.New()

	dec, err := o.DecryptBytes(secret, []byte(encrypted))
	if err != nil {
		return "", err
	}

	return string(dec), nil
}
