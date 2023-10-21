package envs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

/*
Generates 32 byte key using scrypt. Saves salt as salt.data in current dir.
returns key as byte slice on success.
*/
func genKey(pass []byte) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	saltFile, err := os.Create("salt.data")
	if err != nil {
		return nil, err
	}
	defer saltFile.Close()
	saltFile.Write(salt)

	key, err := scrypt.Key(pass, salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return key, nil

}

// Returns 32 byte key from password. Reads salt from salt file then generates a key from the pass and salt.
func getKeyFromSaltFile(pass []byte) ([]byte, error) {
	saltFile, err := os.Open("salt.data")
	if err != nil {
		return nil, err
	}
	defer saltFile.Close()

	salt, err := io.ReadAll(saltFile)
	if err != nil {
		return nil, err
	}

	return scrypt.Key(pass, salt, 32768, 8, 1, 32)
}

// Symmetrical AES encryption based off key.
func EncryptBytesAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//allocate enough space on buffer for iv vector + provided plaintext
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	//Encrypt the byte stream
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return cipherText, nil
}

func DecryptBytesAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short (less than 16 bytes)")
	}

	//Pull apart iv from ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	//Decrypt stream
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
