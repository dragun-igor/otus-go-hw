package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	t.Run("incorrect offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", -100, 0)
		assert.Equal(t, ErrIncorrectOffset, err)
	})
	t.Run("incorrect limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, -100)
		assert.Equal(t, ErrIncorrectLimit, err)
	})
	t.Run("offset > length", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 50_000, 0)
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	})
	t.Run("limit > length", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out.txt", 0, 50_000)
		filePtr, _ := os.Open("out.txt")
		r := io.Reader(filePtr)
		result, _ := ioutil.ReadAll(r)
		assert.Equal(t, nil, err)
		assert.Equal(t, []byte{}, result)
	})
}
