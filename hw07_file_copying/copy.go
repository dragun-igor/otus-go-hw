package main

import (
	"bufio"
	"errors"
	"fmt"
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
	// Если длина 0 - файл не поддерживается
	// Если ограничение ноль, приравниваем длине файла
	stat, err := fromFile.Stat()
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		return ErrUnsupportedFile
	}
	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	fmt.Println(stat.Size())
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
	// Создаём буферизированные reader и writer
	reader := bufio.NewReader(fromFile)
	writer := bufio.NewWriter(toFile)
	// Создаём прогресс бар
	bar := pb.StartNew(int(limit))
	// Побайтово копируем данные
	// Я использовал ReadByte() для большей наглядности прогресс бара (он так более плавный)
	// Намерено не использовал метод CopyN
	// Знаю, что можно использовать io.reader.Read() и []Byte и ещё кучу вариантов
	for i := 0; int64(i) < limit; i++ {
		b, err := reader.ReadByte()                // Читаем байт
		if err != nil && !errors.Is(err, io.EOF) { // Если ошибка и ошибка не End Of File, то возвращаем ошибку
			return err
		}
		if errors.Is(err, io.EOF) { // Если ошибка End Of File прерываем цикл, не записывая байт
			break
		}
		err = writer.WriteByte(b) // Записываем байт
		if err != nil {
			return err
		}
		bar.Add(1)                   // Добавляем единицу к прогресс бару
		time.Sleep(time.Millisecond) // Для наглядности прогресс бара
	}
	// Переносим данные из буффера
	if err := writer.Flush(); err != nil {
		return err
	}
	bar.Finish()
	_, err = io.CopyN(writer, reader, 12)
	return nil
}
