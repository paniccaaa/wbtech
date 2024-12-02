package solutions

import (
	"fmt"
	"slices"
)

// DONE
// Реализовать бинарный поиск встроенными методами языка.

func Solve17() {
	binarySearch := func(arr []int, target int) int {
		low := 0
		high := len(arr) - 1

		for low <= high {
			middle := (low + high) / 2

			if target == arr[middle] {
				return middle
			}

			if arr[middle] > target {
				high = middle - 1
			}

			if arr[middle] < target {
				low = middle + 1
			}
		}

		return -1
	}

	arr := []int{121, 56, 23, 78, 2, 43, 90}
	slices.Sort(arr)

	result := binarySearch(arr, 2)
	fmt.Println(result)
}
