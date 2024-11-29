package solutions

import (
	"fmt"
)

// DONE
// Разработать конвейер чисел.
// Даны два канала:
// в первый пишутся числа (x) из массива,
// во второй — результат операции x*2, после чего данные из второго канала должны выводиться в stdout.

func Solve9() {
	arr := []int{1, 12, 123, 4, 5, 8, 12, 43, 78, 89}

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for _, v := range arr {
			ch1 <- v
		}
		close(ch1)
	}()

	go func() {
		for v := range ch1 {
			ch2 <- v * 2
		}
		close(ch2)
	}()

	for v := range ch2 {
		fmt.Println("hello from ch2 channel", v)
	}
}
