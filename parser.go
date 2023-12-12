package envs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func parseEnvFile(envFilePath string) (map[string]string, error) {
	envMap := make(map[string]string)

	envFile, err := os.Open(envFilePath)
	if err != nil {
		return nil, err
	}
	defer envFile.Close()

	envFileScanner := bufio.NewScanner(envFile)
	for envFileScanner.Scan() {
		line := strings.TrimSpace(envFileScanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("error: malformed line in env file: %s", strings.Join(tokens, " "))
		}

		envKey := strings.TrimSpace(tokens[0])
		envVal := strings.TrimSpace(tokens[1])

		envMap[envKey] = envVal
	}

	if err := envFileScanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}

func parseEnvSFile(envSFilePath string) (map[string]string, error) {
	envMap := make(map[string]string)

	envSFileContent, err := os.ReadFile(envSFilePath)
	if err != nil {
		return nil, err
	}
	
	envFileBytes, err := DecryptBytesAES([]byte("testkeyaaaaabaaa"), envSFileContent)
	if err != nil {
		return nil, err
	}

	envFileBytesBuf := bytes.NewBuffer(envFileBytes)

	envFileScanner := bufio.NewScanner(envFileBytesBuf)
	lineNum := 0 // To keep track of the line number

	for envFileScanner.Scan() {
		lineNum++
		line := strings.TrimSpace(envFileScanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			if utf8.ValidString(strings.Join(tokens, " ")) {
				return nil, fmt.Errorf("error: malformed line at line number: %d in env file: %s", lineNum, line)
			}
			return nil, errors.New("error: could not decrypt envs file")
		}

		envKey := strings.TrimSpace(tokens[0])
		envVal := strings.TrimSpace(tokens[1])

		// Handle escape characters if needed

		envMap[envKey] = envVal
	}

	if err := envFileScanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}
