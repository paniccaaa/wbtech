package solutions

import "fmt"

// DONE
// Поменять местами два числа без создания временной переменной.

func Solve13() {
	a := 5
	b := 9
	a, b = b, a
	fmt.Println(a, b)
}
