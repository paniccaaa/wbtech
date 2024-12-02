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

	var num int64 = 66

	num = setBit(num, 4, true)
	fmt.Println(num)

	num = setBit(num, 4, false)
	fmt.Println(num)
}
