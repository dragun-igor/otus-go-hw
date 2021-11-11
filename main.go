package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func digit(r rune, prevr rune) (string, error) {
	if unicode.IsLetter(prevr) || prevr == 10 {
		num, _ := strconv.Atoi(string(r))
		result := strings.Repeat(string(prevr), num)
		return result, nil
	}
	return "", ErrInvalidString
}

func notdigit(r rune, prevr rune) (string, error) {
	if unicode.IsLetter(prevr) {
		return string(prevr), nil
	}
	if unicode.IsDigit(prevr) {
		return "", nil
	}
	if prevr == 10 {
		return "\n", nil
	}
	return "", ErrInvalidString
}

func test(instr string) (string, error) {
	var result strings.Builder
	var prevr rune
	var err error
	var str string
	runes := []rune(instr)
	for i, r := range runes {
		switch i {
		case 0:
			if !unicode.IsLetter(r) && r != 10 {
				err = ErrInvalidString
			}
		default:
			if unicode.IsDigit(r) {
				str, err = digit(r, prevr)
			}
			if unicode.IsLetter(r) || r == 10 {
				str, err = notdigit(r, prevr)
			}
		case len(runes) - 1:
			if unicode.IsDigit(r) {
				str, err = digit(r, prevr)
			}
			if unicode.IsLetter(r) || r == 10 {
				str, err = notdigit(r, prevr)
				str += string(r)
			}
		}
		if err != nil {
			return "", err
		}
		result.WriteString(str)
		prevr = r
	}
	return result.String(), nil
}

func main() {

	fmt.Println(test(45))

}
