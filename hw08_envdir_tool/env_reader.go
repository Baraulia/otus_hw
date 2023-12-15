package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return Environment{}, nil
	}

	env := make(Environment)

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}
		if fileInfo.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
		} else {
			filePath := fmt.Sprintf("%s/%s", dir, fileInfo.Name())
			value, err := readFile(filePath)
			if err != nil {
				return nil, err
			}

			value = strings.ReplaceAll(strings.TrimRight(value, " \t\n"), "\x00", "\n")

			env[file.Name()] = EnvValue{
				Value:      value,
				NeedRemove: true,
			}
		}
	}

	return env, nil
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error while openning file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return line, nil
		}
		return "", fmt.Errorf("error while reading file %s: %w", filePath, err)
	}

	return line, nil
}
