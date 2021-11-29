package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("incorrect offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", -100, 0, "byte")
		assert.Equal(t, ErrIncorrectOffset, err)
	})
	t.Run("incorrect limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, -100, "byte")
		assert.Equal(t, ErrIncorrectLimit, err)
	})
	t.Run("incorrect mode", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, 0, "bite")
		assert.Equal(t, ErrIncorrectMode, err)
	})
	t.Run("offset > length", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 50_000, 0, "byte")
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	})
	t.Run("limit > length", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, 50_000, "byte")
		filePtr, _ := os.Open("out.txt")
		r := io.Reader(filePtr)
		result, _ := ioutil.ReadAll(r)
		assert.Equal(t, nil, err)
		assert.Equal(t, []byte{}, result)
	})
}
