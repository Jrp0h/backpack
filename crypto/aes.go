package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/Jrp0h/backpack/utils"
)

var validTypes = map[string]func(string) ([]byte, error){
	"base64": keyFromBase64,
	"file":   keyFromFile,
	"raw":    keyFromRaw,
	"hex":    keyFromHex,
}

type AESCrypo struct {
	key []byte
}

func (crypter *AESCrypo) Encrypt(path string) (outErr error) {
	defer func() {
		recErr := recover()
		if recErr != nil {
			outErr = fmt.Errorf("crypto/aes: Couldn't encrypt '%s'\n%s", path, recErr)
		}
	}()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("crypto/aes: Couldn't read file '%s'\n%s", path, err.Error())
	}

	block, err := aes.NewCipher(crypter.key)
	if err != nil {
		return fmt.Errorf("crypto/aes: Couldn't create AES cipher.\n%s", err.Error())
	}

	data, err = pkcs7Pad(data, block.BlockSize())
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("crypto/aes: Couldn't initialize iv.\n%s", err.Error())
	}

	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], data)

	if err = ioutil.WriteFile(path, ciphertext, 0644); err != nil {
		return fmt.Errorf("crypto/aes: Couldn't write file '%s'\n%s", path, err.Error())
	}

	return nil
}

func (crypter *AESCrypo) Decrypt(path string) (outErr error) {
	defer func() {
		recErr := recover()
		if recErr != nil {
			outErr = fmt.Errorf("crypto/aes: Couldn't decrypt '%s'\n%s", path, recErr)
		}
	}()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("crypto/aes: Couldn't read file '%s'\n%s", path, err.Error())
	}

	block, err := aes.NewCipher(crypter.key)
	if err != nil {
		return fmt.Errorf("crypto/aes: Couldn't create AES cipher.\n%s", err.Error())
	}

	decrypted := make([]byte, len(data)-aes.BlockSize)
	iv := data[:aes.BlockSize]

	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(decrypted, data[aes.BlockSize:])

	paddingToRemove := decrypted[len(decrypted)-1]
	decrypted = decrypted[:len(decrypted)-int(paddingToRemove)]

	if err = ioutil.WriteFile(path, decrypted, 0644); err != nil {
		return fmt.Errorf("crypto/aes: Couldn't write file '%s'\n%s", path, err.Error())
	}

	return nil
}

func loadAES(data *map[string]string) (Crypto, error) {
	keyFromConfig, err := utils.ValueOrErrorString(data, "key", "crypto/aes")

	if err != nil {
		return nil, err
	}

	var key []byte

	keyFound := false
	for k, v := range validTypes {
		if strings.HasPrefix(keyFromConfig, k+":") {
			keyFound = true

			key, err = v(strings.Replace(keyFromConfig, k+":", "", 1))
			if err != nil {
				return nil, err
			}
		}
	}

	if !keyFound {
		return nil, fmt.Errorf("crypto/aes: Invalid key type '%s'", strings.Split(keyFromConfig, ":")[0])
	}

	if len(key) != 32 && len(key) != 24 && len(key) != 16 {
		return nil, fmt.Errorf("crypto/aes: Invalid key size '%d'", len(key))
	}
	utils.Log.Debug("crypto/aes: Key size is %dbit", len(key))

	return &AESCrypo{
		key,
	}, nil
}

func keyFromBase64(key string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: Invalid base64 key")
	}

	return decoded, nil
}

func keyFromFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: Couldn't open key file.\n'%s'", path)
	}

	return data, nil
}

func keyFromRaw(key string) ([]byte, error) {
	return []byte(key), nil
}

func keyFromHex(key string) ([]byte, error) {
	decoded, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("crypto/aes: Invalid hex key")
	}

	return decoded, nil
}

// Function taken from https://stackoverflow.com/questions/66221371/encrypt-aes-string-with-go-and-decrypt-with-crypto-js
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, fmt.Errorf("crypto/aes: Invalid blocksize")
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("crypto/aes: Invalid PKCS7 data (empty or not padded)")
	}

	n := blocksize - (len(b) % blocksize) // size
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}
