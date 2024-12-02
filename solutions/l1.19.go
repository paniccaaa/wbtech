package solutions

import "fmt"

// Разработать программу, которая переворачивает подаваемую на ход строку
// (например: «главрыба — абырвалг»).
// Символы могут быть unicode.

func Solve19() {
	reverseString := func(str string) string {
		r := []rune(str)
		low, high := 0, len(r)-1

		for low <= high {
			r[low], r[high] = r[high], r[low]
			low++
			high--
		}

		return string(r)
	}

	var str = "главрыба"
	fmt.Println(reverseString(str))
}
