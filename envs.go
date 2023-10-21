package envs

import (
	"errors"
	"io/fs"
	"log"
	"os"
)

func init() {
	if _, err := os.Stat("envs"); err == nil {
		return
	} else if !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("Unknown error occurred: %v", err)
	}

	envFileContent, err := os.ReadFile(".env")
	if err != nil {
		log.Fatalf("Could not open .env file: %v", err)
	}

	if err := setupENVS(envFileContent); err != nil {
		log.Fatalf("Error setting up .envs: %v", err)
	}
}

func setupENVS(originalEnvFileContent []byte) error {
	if err := os.Mkdir("envs", 0744); err != nil {
		return err
	}

	envsFile, err := os.Create("envs/.envs")
	if err != nil {
		return err
	}
	defer envsFile.Close()

	envsFileContent, err := EncryptBytesAES([]byte("testkeyaaaaaaaaa"), originalEnvFileContent)
	if err != nil {
		return err
	}

	if _, err := envsFile.Write(envsFileContent); err != nil {
		return err
	}

	if err := os.Remove(".env"); err != nil {
		log.Printf("Warning: Could not remove .env file: %v", err)
	}

	return nil
}
