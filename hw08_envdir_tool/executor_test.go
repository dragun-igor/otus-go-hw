package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
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
	// Даём команду слип, и проверяем что функция выполняется от 5 до 6 секунд, проверяем код ошибки
	t.Run("sleep and error code", func(t *testing.T) {
		start := time.Now()
		errorCode := RunCmd([]string{"sleep", "5"}, env)
		elapsedTime := time.Since(start)
		require.Equal(t, 0, errorCode)
		require.GreaterOrEqual(t, elapsedTime, time.Second*5)
		require.LessOrEqual(t, elapsedTime, time.Second*6)
	})

	// Параллельно eventually запускаем функцию слип
	go func() {
		_ = RunCmd([]string{"sleep", "5"}, env)
	}()

	// Проверяет, что переменные окружения были заданы и удалены
	require.Eventually(t, func() bool {
		val, ok := os.LookupEnv("HELLO")
		if !ok || val != "\"hello\"" {
			return false
		}
		_, ok = os.LookupEnv("UNSET")
		return !ok
	}, time.Second*5, time.Second)

	// Проверяем что файла с названием isnotexist нету, создаём файл, проверяем, удаляем файл, проверяем
	t.Run("create and delete file", func(t *testing.T) {
		_, err := os.Stat("./isnotexist")
		require.True(t, os.IsNotExist(err))
		errorCode := RunCmd([]string{"touch", "isnotexist"}, Environment{})
		_, err = os.Stat("./isnotexist")
		require.Nil(t, err)
		require.Equal(t, 0, errorCode)
		errorCode = RunCmd([]string{"rm", "./isnotexist"}, Environment{})
		_, err = os.Stat("./isnotexist")
		require.True(t, os.IsNotExist(err))
		require.Equal(t, 0, errorCode)
	})
}
