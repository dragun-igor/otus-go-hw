package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expected := Environment{
		"BAR": EnvValue{
			Value:      "bar",
			NeedRemove: false,
		},
		"EMPTY": EnvValue{
			Value:      "",
			NeedRemove: false,
		},
		"FOO": EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		},
		"HELLO": EnvValue{
			Value:      "\"hello\"",
			NeedRemove: false,
		},
		"UNSET": EnvValue{
			Value:      "",
			NeedRemove: true,
		},
	}
	// Проверка работы функции
	t.Run("without error", func(t *testing.T) {
		actual, err := ReadDir("./testdata/env")
		require.Equal(t, expected, actual, "actual and expected are not same")
		require.Nil(t, err, "function returns with error")
	})
	// При добавлении файла со знаком равно, функция должна игнорировать этот файл и не добавлять в мапу
	t.Run("= in file name", func(t *testing.T) {
		// Создаём файл с = в названии перед окончанием функции файл удаляем
		_, _ = os.Create("./testdata/env/F=OO")
		defer os.Remove("./testdata/env/F=OO")
		actual, err := ReadDir("./testdata/env")
		require.Equal(t, expected, actual, "actual and expected are not same")
		require.Nil(t, err, "function returns with error")
	})
}
