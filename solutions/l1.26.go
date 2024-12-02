package solutions

import (
	"fmt"
	"strings"
)

// DONE
// Разработать программу, которая проверяет, что все символы в строке уникальные
// (true — если уникальные, false etc). Функция проверки должна быть регистронезависимой.

// Например:
// abcd — true
// abCdefAaf — false
// aabcd — false

func Solve26() {
	strIsUnique := func(str string) bool {
		str = strings.ToLower(str)
		m := make(map[rune]int)

		for _, v := range str {
			m[v] += 1

			if val := m[v]; val > 1 {
				return false
			}
		}

		return true
	}

	arr := []string{"abcd", "bCdeAaf", "aabsd"}

	for _, v := range arr {
		fmt.Println(strIsUnique(v))
	}
}
