package solutions

import "fmt"

// DONE
// Реализовать быструю сортировку массива (quicksort) встроенными методами языка.

func Solve16() {
	arr := []int{45, 1, 32, 4, 67, 5, 9, 44}

	arr = quicksort(arr, 0, len(arr)-1)
	fmt.Println(arr)
}

func partition(arr []int, low, high int) ([]int, int) {
	pivot := arr[high]
	i := low

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[i], arr[high] = arr[high], arr[i]

	return arr, i
}

func quicksort(arr []int, low, high int) []int {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quicksort(arr, low, p-1)
		arr = quicksort(arr, p+1, high)
	}

	return arr
}
