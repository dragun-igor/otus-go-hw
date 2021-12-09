package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Отключил линтер на строку, потому что жалуется что аргументов может не хватать, но проверка есть в main.go
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	// Задаём и удаляем переменные окружения в соответствии с информацией из мапы env
	for key, val := range env {
		envVal, ok := os.LookupEnv(key) // Проверяем наличие переменной
		switch {
		case ok && val.NeedRemove: // Если переменная есть и нужно удалить - ансетим её
			os.Unsetenv(key)
		case !ok && val.NeedRemove: // Если переменной нет и нужно удалить - пропускаем цикл
			continue
		case ok && !val.NeedRemove: // Если переменная есть и не надо удалять и значения разные, ансетим и сетим
			if envVal != val.Value {
				os.Unsetenv(key)
				os.Setenv(key, val.Value)
			}
		default: // Иначе просто сетим
			os.Setenv(key, val.Value)
		}
	}

	// Пробрасываем стандартные потоки в вызываемую программу
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Запускаем, если была ошибка возвращаем код ошибки
	var ee *exec.ExitError
	if err := command.Run(); err != nil {
		if errors.As(err, &ee) {
			returnCode = ee.ExitCode()
			log.Println("exit code error:", returnCode)
		} else {
			log.Printf("general error: %v", err)
		}
	} else {
		log.Println("success!")
	}
	return
}
