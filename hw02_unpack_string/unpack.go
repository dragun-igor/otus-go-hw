package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrFirstChar     = errors.New("the first char must be letter, rune(92) or rune(10)")
	ErrInvalidChar   = errors.New("invalid char")
	ErrTwoDigit      = errors.New("two digit in a row")
	ErrProtLetter    = errors.New("protected letter char")
	ErrLastBackslash = errors.New("last rune is not protected backslash")
	ErrInvalidString = errors.New("invalid string")
)

const (
	BACKSLASH = 92 // руна бэкслэша
	LINEFEED  = 10 // руна переноса строки
)

type token struct {
	r         rune   // текущая руна
	prev      *token // предыдущий элемент
	protected bool   // экранирован или нет
	digit     bool   // цифра
	letter    bool   // буква
	backslash bool   // обратный слэш
	linefeed  bool   // перенос строки
	first     bool   // первый элемент
	last      bool   // последний элемент
	out       string // конечная строка
	err       error  // ошибка
}

func (token *token) Digit() (string, error) { // Функции для обработки цифровой руны
	if token.first {
		return "", ErrFirstChar
	}
	if token.prev.digit && !token.prev.protected {
		return "", ErrTwoDigit
	}
	num, errAtoi := strconv.Atoi(string(token.r))
	if errAtoi != nil {
		return "", errAtoi
	}

	if token.protected {
		return string(token.r), nil
	}
	if !(token.prev.backslash && !token.prev.protected) {
		if num > 0 {
			return strings.Repeat(string(token.prev.r), num-1), nil
		}
		if num == 0 {
			token.prev.out = ""
			return "", nil
		}
	}
	return "", ErrInvalidString
}

func (token *token) NotDigit() (string, error) { // Метод обработки нецифровой руны (буквы, бэкслэш, перенос строки)
	if token.letter && token.protected {
		return "", ErrProtLetter
	}
	if token.backslash && !token.protected && token.last {
		return "", ErrLastBackslash
	}
	if token.backslash && !token.protected {
		return "", nil
	}
	return string(token.r), nil
}

func Unpack(in string) (string, error) {
	if in == "" {
		return "", nil
	}
	var result strings.Builder
	runes := []rune(in)
	tokens := make([]token, len(runes))
	tokens[0].first = true
	tokens[len(tokens)-1].last = true
	for i, r := range runes {
		tokens[i].r = r
		tokens[i].digit = unicode.IsDigit(r)
		tokens[i].letter = unicode.IsLetter(r)
		tokens[i].backslash = r == BACKSLASH
		tokens[i].linefeed = r == LINEFEED
		if !tokens[i].digit && !tokens[i].letter && !tokens[i].backslash && !tokens[i].linefeed {
			return "", ErrInvalidChar
		}
		if i > 0 {
			tokens[i].prev = &tokens[i-1]
			tokens[i].protected = tokens[i].prev.r == BACKSLASH && !tokens[i].prev.protected
		}
		if tokens[i].digit {
			tokens[i].out, tokens[i].err = tokens[i].Digit()
		} else {
			tokens[i].out, tokens[i].err = tokens[i].NotDigit()
		}
		if tokens[i].err != nil {
			return "", tokens[i].err
		}
	}
	for _, t := range tokens {
		result.WriteString(t.out)
	}
	return result.String(), nil
}
