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
	// Создаём мапу
	env := Environment{}
	// Открываем директорию, получаем список файлов
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fileStat := range files {
		fileName := fileStat.Name()
		// Если в имени файла = или это является директорией, пропускаем
		if strings.Contains(fileName, "=") || fileStat.IsDir() {
			continue
		}
		newEnvValue := EnvValue{}
		// Если размер 0, то переменную необходимо будет удалить
		if fileStat.Size() == 0 {
			newEnvValue.NeedRemove = true
		}
		fileDir := path.Join(dir, "/", fileName)
		// Открываем файл, читаем из него байты до переноса строки
		file, err := os.Open(fileDir)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(file)
		byteLine, err := reader.ReadBytes(0x0a)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		// Заменяем терминальный ноль на перенос строки, справа удаляем пробелы табуляцию
		// Перенос появляется после ReadBytes(), так как разделитель остаётся в массиве байт
		byteLine = bytes.ReplaceAll(byteLine, []byte{NULL}, []byte{LINEFEED})
		byteLine = bytes.TrimRight(byteLine, " \t\n")
		// Забиваем значение, помещаем структуру в мапу по ключу-названию файла
		newEnvValue.Value = string(byteLine)
		env[fileName] = newEnvValue
	}
	return env, nil
}
