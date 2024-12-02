package solutions

import "fmt"

// DONE
// Имеется последовательность строк - (cat, cat, dog, cat, tree)
//  создать для нее собственное множество.

func Solve12() {
	arr := []string{"cat", "cat", "dog", "cat", "tree"}

	set := make(map[string]bool)

	for _, v := range arr {
		set[v] = true
	}

	fmt.Println(set)
}
