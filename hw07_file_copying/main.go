package main

import (
	"flag"
	"log"
)

var (
	from, to, mode string
	limit, offset  int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.StringVar(&mode, "mode", "byte", "byte or rune mode")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	// Добавлен флаг для работы с рунами, в test.sh добавил тесты для русского языка
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	err := Copy(from, to, offset, limit, mode)
	if err != nil {
		log.Fatal(err)
	}
}
