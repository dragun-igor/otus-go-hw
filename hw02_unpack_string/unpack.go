package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func digit(r rune, prevr rune) (string, error) {
	if unicode.IsLetter(prevr) || prevr == 10 { // r == 10 => "\n"
		num, _ := strconv.Atoi(string(r))
		result := strings.Repeat(string(prevr), num)
		return result, nil
	}
	return "", ErrInvalidString
}

func notdigit(prevr rune) (string, error) {
	if unicode.IsLetter(prevr) {
		return string(prevr), nil
	}
	if unicode.IsDigit(prevr) {
		return "", nil
	}
	if prevr == 10 { // r == 10 => "\n"
		return "\n", nil
	}
	return "", ErrInvalidString
}

func Unpack(in string) (string, error) {
	var result strings.Builder
	var prevr rune
	var err error
	var str string
	runes := []rune(in)
	for i, r := range runes {
		switch i {
		case 0:
			if !unicode.IsLetter(r) && r != 10 { // r == 10 => "\n"
				err = ErrInvalidString
			}
		case len(runes) - 1:
			if unicode.IsDigit(r) {
				str, err = digit(r, prevr)
			}
			if unicode.IsLetter(r) || r == 10 { // r == 10 => "\n"
				str, err = notdigit(prevr)
				str += string(r)
			}
		default:
			if unicode.IsDigit(r) {
				str, err = digit(r, prevr)
			}
			if unicode.IsLetter(r) || r == 10 { // r == 10 => "\n"
				str, err = notdigit(prevr)
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
