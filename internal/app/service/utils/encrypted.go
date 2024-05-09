package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func EncryptAES256CBC(plaintext string, key string, iv string) (string, error) {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), aes.BlockSize)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func DecryptAES256CBC(ciphertext string, key string, iv string) (string, error) {
	bKey := []byte(key)
	bIV := []byte(iv)
	bCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks(bCiphertext, bCiphertext)

	unpaddedText, err := PKCS5Unpadding(bCiphertext)
	if err != nil {
		return "", err
	}

	return string(unpaddedText), nil
}

func PKCS5Unpadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	padding := int(data[len(data)-1])
	if padding >= len(data) || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}
