package dev02

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func Unpack(s string) (string, error) {
	var builder strings.Builder
	var prevRune rune
	escaped := false
	multiplier := ""

	for _, r := range s {
		if escaped {
			builder.WriteRune(r)
			prevRune = r
			escaped = false
			continue
		}

		if r == '\\' {
			if escaped {

				builder.WriteRune(r)
				prevRune = 0
			}

			escaped = true
			if multiplier != "" {
				count, err := strconv.Atoi(multiplier)
				if err != nil {
					return string(r), errors.New("некорректная строка")
				}

				builder.WriteString(strings.Repeat(string(prevRune), count-1))
				multiplier = ""
			}
			continue
		}

		if unicode.IsDigit(r) {
			if prevRune == 0 {
				return "", errors.New("incorrect string")
			}

			multiplier += string(r)
		} else {
			if multiplier != "" {
				count, err := strconv.Atoi(multiplier)
				if err != nil {
					return string(r), errors.New("incorrect string")
				}

				builder.WriteString(strings.Repeat(string(prevRune), count-1))
				multiplier = ""
			}

			builder.WriteRune(r)
			prevRune = r
		}
	}

	if escaped {
		return "", errors.New("incorrect string")
	}

	if multiplier != "" {
		count, err := strconv.Atoi(multiplier)
		if err != nil {
			return multiplier, errors.New("incorrect string")
		}
		builder.WriteString(strings.Repeat(string(prevRune), count-1))
	}

	return builder.String(), nil
}
