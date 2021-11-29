package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIncorrectOffset       = errors.New("offset less then zero")
	ErrIncorrectLimit        = errors.New("limit less then zero")
	ErrIncorrectMode         = errors.New("incorrect mode value")
)

func Copy(fromPath, toPath string, offset, limit int64, mode string) error {
	// Если режим ни байт, ни руна - ошибка
	if mode != "byte" && mode != "rune" {
		return ErrIncorrectMode
	}
	// Если смещение меньше 0 - ошибка
	if offset < 0 {
		return ErrIncorrectOffset
	}
	// Если ограничение меньше 0 - ошибка
	if limit < 0 {
		return ErrIncorrectLimit
	}
	// Открываем файл, возвращаем ошибку, перед возвратом закрываем файл
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0740)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	// Если смещение больше длины файла - ошибка
	// Если ограничение ноль, приравниваем длине файла
	if stat, err := fromFile.Stat(); err != nil {
		return err
	} else {
		if stat.Size() < offset {
			return ErrOffsetExceedsFileSize
		}
		if limit == 0 || limit+offset > stat.Size() {
			limit = stat.Size() - offset
		}
	}
	// Создаём файл
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()
	// Создаём буферезированные reader и writer
	reader := bufio.NewReader(fromFile)
	writer := bufio.NewWriter(toFile)
	// Создаём прогресс бар
	bar := pb.StartNew(int(limit))
	// В зависимости от режима копируем побайтово (порунно) копируем данные
	// Можно сделать, чтобы переносилось большими частями, но тогда прогресс бар будет не видно
	// В рунном режиме прогресс бар немного поломан, так как я не знаю как посчитать количество рун в файле
	switch mode {
	case "rune":
		for i := 0; int64(i) < limit+offset; i++ {
			r, _, err := reader.ReadRune()
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				break
			}
			if i >= int(offset) {
				_, err = writer.WriteRune(r)
				bar.Add(1)
				time.Sleep(time.Millisecond)
			}
		}
	case "byte":
		for i := 0; int64(i) < limit+offset; i++ {
			b, err := reader.ReadByte()
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				break
			}
			if i >= int(offset) {
				err = writer.WriteByte(b)
				bar.Add(1)
				time.Sleep(time.Millisecond)
			}
		}
	}
	// Переносим данные из буффера
	if err := writer.Flush(); err != nil {
		return err
	}
	bar.Finish()
	return nil
}
