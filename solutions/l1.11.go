package solutions

import "fmt"

// DONE
// Реализовать пересечение двух неупорядоченных множеств.

func Solve11() {
	set1 := map[int]bool{9: true, 2: true, 4: true}
	set2 := map[int]bool{2: true, 10: true, 8: true}

	intersec := make(map[int]bool)

	for k := range set1 {
		if set2[k] {
			intersec[k] = true
		}
	}

	// output: map[2:true 4:true]
	fmt.Println(intersec)
}
