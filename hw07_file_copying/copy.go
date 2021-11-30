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
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIncorrectOffset       = errors.New("offset less then zero")
	ErrIncorrectLimit        = errors.New("limit less then zero")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Если смещение меньше 0 - ошибка
	if offset < 0 {
		return ErrIncorrectOffset
	}
	// Если ограничение меньше 0 - ошибка
	if limit < 0 {
		return ErrIncorrectLimit
	}
	// Открываем файл, возвращаем ошибку, перед возвратом закрываем файл, режим READ_ONLY
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o740)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	// Если смещение больше длины файла - ошибка
	// Если ограничение ноль, приравниваем длине файла
	stat, err := fromFile.Stat()
	if err != nil {
		return err
	}
	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit+offset > stat.Size() {
		limit = stat.Size() - offset
	}
	// Создаём файл
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()
	// Устанавливаем смещение
	fromFile.Seek(offset, 0)
	// Создаём буферезированные reader и writer
	reader := bufio.NewReader(fromFile)
	writer := bufio.NewWriter(toFile)
	// Создаём прогресс бар
	bar := pb.StartNew(int(limit))
	// Побайтово копируем данные
	// Можно сделать, чтобы переносилось большими частями, но тогда прогресс бар будет не видно
	for i := 0; int64(i) < limit; i++ {
		b, err := reader.ReadByte()
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if errors.Is(err, io.EOF) {
			break
		}
		err = writer.WriteByte(b)
		if err != nil {
				return err
			}
			bar.Add(1)
			time.Sleep(time.Millisecond) // Для наглядности програсс бара
		}
	// Переносим данные из буффера
	if err := writer.Flush(); err != nil {
		return err
	}
	bar.Finish()
	return nil
}
