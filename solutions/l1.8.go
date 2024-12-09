package solutions

import "fmt"

// TODO
// Дана переменная int64. Разработать программу которая устанавливает i-й бит в 1 или 0.

func Solve8() {
	setBit := func(num int64, i int, setToOne bool) int64 {
		if setToOne {
			return num | (1 << i)
		}

		return num & ^(1 << i)
	}

	// 65        = 1000001
	// mask      = 00010000
	// 65 | mask = 1010001

	num := setBit(65, 4, true)
	fmt.Println(num) // 81

	// 65          = 1000001
	// mask        = 00010000
	// invert      = 11101111
	// 65 | invert = 1000001

	num2 := setBit(65, 4, false)
	fmt.Println(num2) // 65
}
