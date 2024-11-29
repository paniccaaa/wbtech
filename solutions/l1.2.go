package solutions

import (
	"fmt"
	"sync"
)

// DONE
// Написать программу, которая конкурентно рассчитает значение квадратов чисел
// взятых из массива (2,4,6,8,10) и выведет их квадраты в stdout.

func Solve2() {

	arr := [5]int{2, 4, 6, 8, 10}
	var wg sync.WaitGroup

	wg.Add(len(arr))

	for _, v := range arr {
		go printSquare(v, &wg)
	}

	wg.Wait()
}

func printSquare(v int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(v * v)
}
