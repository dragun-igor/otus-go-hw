package hw02_unpack_string

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func digit(r rune, prevr rune, protected bool) (string, error) { // Функция вызывается при цифровой руне
	// Если предыдущая руна - буква, перенос или экранированные слэш или цифра, то выполняем размножение и возвращаем строку
	if unicode.IsLetter(prevr) || prevr == 10 || ((prevr == 92 || unicode.IsDigit(prevr)) && protected) { // r == 10 => "\n", r == 92 => "\"
		num, _ := strconv.Atoi(string(r))
		result := strings.Repeat(string(prevr), num)
		return result, nil
	}
	// Если предыдущая руна - неэкранированный слэш, то возвращаем пустую строку
	if prevr == 92 && !protected {
		return "", nil
	}
	// В остальных случаях возвращаем ошибку
	return "", ErrInvalidString
}

func notdigit(r rune, prevr rune, protected bool) (string, error) { // Функция вызывается при нецифровой руне (буква, перенос, слэш)
	// Если предыдущая руна - буква, то возвращаем букву
	if unicode.IsLetter(prevr) {
		return string(prevr), nil
	}
	// Если предыдущая руна - неэкранированная цифра, то возвращаем пустую строку, если - экранированная цифра, то возвращаем цифру
	if unicode.IsDigit(prevr) {
		if !protected {
			return "", nil
		} else {
			return string(prevr), nil
		}
	}
	// Если предыдущая руна - перенос строки, то возвращаем перенос строки
	if prevr == 10 {
		return "\n", nil
	}
	// Если предыдущая руна экранированный слэш, то возвращаем слэш, если неэкранированый - пустую строку
	if prevr == 92 {
		if protected {
			return "\\", nil
		}
		if !protected && r == 92 {
			return "", nil
		}
	}
	// В остальных случаях возвращаем ошибку
	return "", ErrInvalidString
}

func Unpack(in string) (string, error) {
	var result strings.Builder
	var prevr rune
	var err error
	var str string
	runes := []rune(in)
	// Слайс значений экранирования
	protected := make([]bool, len(runes))
	for i, r := range runes {
		// Если первая руна ни буква, ни перенос, ни слэш, то возвращаем ошибку
		if !unicode.IsLetter(r) && r != 10 && r != 92 && i == 0 {
			err = ErrInvalidString
		}
		// Если предыдущая руна - неэкранированный слэш - взводим флаг экранированирования
		if prevr == 92 && !protected[i-1] {
			protected[i] = true
		}
		if i != 0 {
			if unicode.IsLetter(r) || r == 10 || r == 92 {
				str, err = notdigit(r, prevr, protected[i-1])
			}
			if unicode.IsDigit(r) && i != 0 {
				str, err = digit(r, prevr, protected[i-1])
			}
		}
		// Если руна - последняя и она - не неэкранированная цифра, то добавляем её к окончательному результату
		if i == len(runes)-1 && !(!protected[i] && unicode.IsDigit(r)) {
			str += string(r)
		}
		// Если руна - последняя и она - неэкранированный слэш, возвращаем ошибку
		if i == len(runes)-1 && !protected[i] && r == 92 {
			err = ErrInvalidString
		}
		if err != nil {
			return "", err
		}
		result.WriteString(str)
		prevr = r
	}
	return result.String(), nil
}
