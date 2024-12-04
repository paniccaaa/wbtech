package solutions

import (
	"fmt"
	"time"
)

// DONE
// Реализовать собственную функцию sleep.

func Solve25() {
	sleep := func(d time.Duration) {
		<-time.After(d)
	}

	sleep1 := func(d time.Duration) {
		start := time.Now()
		for time.Since(start) < d {
		}
	}

	fmt.Println(time.Now())
	sleep(time.Second)
	fmt.Println(time.Now())

	fmt.Println(time.Now())
	sleep1(time.Second)
	fmt.Println(time.Now())
}
