package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	t.Run("incorrect offset", func(t *testing.T) {
		f, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(f.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", -100, 0)
		assert.Equal(t, ErrIncorrectOffset, err)
		f.Close()
	})
	t.Run("incorrect limit", func(t *testing.T) {
		f, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(f.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, -100)
		assert.Equal(t, ErrIncorrectLimit, err)
		f.Close()
	})
	t.Run("unsupported file", func(t *testing.T) {
		f, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(f.Name())
		err := Copy("/dev/urandom", "testdata/out.txt", 0, 0)
		assert.Equal(t, ErrUnsupportedFile, err)
		f.Close()
	})
	t.Run("offset > length", func(t *testing.T) {
		f, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(f.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", 50_000, 0)
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
		f.Close()
	})
	t.Run("limit > length", func(t *testing.T) {
		f, _ := os.CreateTemp("testdata", "out.txt")
		defer os.Remove(f.Name())
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, 50_000)
		filePtr, _ := os.Open("out.txt")
		r := io.Reader(filePtr)
		result, _ := io.ReadAll(r)
		assert.Equal(t, nil, err)
		assert.Equal(t, []byte{}, result)
		f.Close()
	})
}
