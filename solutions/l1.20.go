package solutions

import (
	"fmt"
	"strings"
)

// DONE
// Разработать программу, которая переворачивает слова в строке.
// Пример: «snow dog sun — sun dog snow».

func Solve20() {
	reverseWords := func(s string) string {
		arr := strings.Fields(s)
		low, high := 0, len(arr)-1

		for low <= high {
			arr[low], arr[high] = arr[high], arr[low]

			low++
			high--
		}

		return strings.Join(arr, " ")
	}

	var words = "snow dog sun"

	fmt.Println(reverseWords(words))
}
