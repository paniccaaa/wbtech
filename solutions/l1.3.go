package solutions

import (
	"fmt"
	"sync"
)

// Дана последовательность чисел: 2,4,6,8,10.
// Найти сумму их квадратов(22+32+42….) с использованием конкурентных вычислений.

// 4 + 16 + 36 + 64 + 100 == 220
func Solve3() {
	var wg sync.WaitGroup
	var arr = []int{2, 4, 6, 8, 10}

	squaresSum := 0
	squareChan := make(chan int, len(arr))
	wg.Add(len(arr))
	for _, v := range arr {
		go func(v int) {
			defer wg.Done()
			squareChan <- v * v
		}(v)
	}

	go func() {
		wg.Wait()
		close(squareChan)
	}()

	for v := range squareChan {
		squaresSum += v
	}

	fmt.Println(squaresSum)
}
