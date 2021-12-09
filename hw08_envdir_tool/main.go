package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("incorrect input")
	}
	// Аргументы [util name, dir, command, arg...]
	// Нам необходимо извлечь dir и передать в функцию ReadDir
	// И извлечь [command, arg...] и передать в функцию RunCmd
	dir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	// Закрытие утилиты с кодом ошибки
	os.Exit(RunCmd(cmd, env))
}
