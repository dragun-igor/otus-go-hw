package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

const (
	NULL     byte = 0x00 // Терминальный ноль
	LINEFEED byte = 0x0a // Перенос строки
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
	env := Environment{}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fileStat := range files {
		fileName := fileStat.Name()
		fileInfo, err := fileStat.Info()
		if err != nil {
			return nil, err
		}
		if strings.Contains(fileName, "=") || fileStat.IsDir() {
			continue
		}
		newEnvValue := EnvValue{}
		if fileInfo.Size() == 0 {
			newEnvValue.NeedRemove = true
		}
		fileDir := path.Join(dir, "/", fileName)
		file, err := os.Open(fileDir)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(file)
		byteLine, err := reader.ReadBytes(0x0a)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		byteLine = bytes.ReplaceAll(byteLine, []byte{NULL}, []byte{LINEFEED})
		byteLine = bytes.TrimRight(byteLine, " \t\n")
		newEnvValue.Value = string(byteLine)
		env[fileName] = newEnvValue
	}
	return env, nil
}
