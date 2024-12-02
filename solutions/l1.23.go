package solutions

import "fmt"

// DONE
// Удалить i-ый элемент из слайса.

func Solve23() {
	removeElem := func(arr []int, index int) []int {
		// dst = [3, 6, 7]; src = [6, 7] => [32 1 6 7 7]
		copy(arr[index:], arr[index+1:])

		// cut last duplicate
		arr = arr[:len(arr)-1]

		return arr
	}

	arr := []int{32, 1, 3, 6, 7}
	fmt.Println(removeElem(arr, 2))
}
