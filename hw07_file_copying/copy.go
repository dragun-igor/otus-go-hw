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
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrIncorrectOffset
	}
	if limit < 0 {
		return ErrIncorrectLimit
	}
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o740)
	if err != nil {
		return err
	}
	defer fromFile.Close()

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
	if limit == 0 || limit+offset > stat.Size() {
		limit = stat.Size() - offset
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()
	fromFile.Seek(offset, 0)
	reader := bufio.NewReader(fromFile)
	writer := bufio.NewWriter(toFile)
	bar := pb.StartNew(int(limit))
	for i := 0; int64(i) < limit; i++ {
		b, err := reader.ReadByte()
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if errors.Is(err, io.EOF) {
			break
		}
		err = writer.WriteByte(b) // Записываем байт
		if err != nil {
			return err
		}
		bar.Add(1)
		time.Sleep(time.Millisecond)
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	bar.Finish()
	return nil
}
