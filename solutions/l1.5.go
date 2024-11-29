package solutions

import (
	"fmt"
	"sync"
	"time"
)

// DONE
// Разработать программу, которая будет последовательно отправлять значения в канал,
// а с другой стороны канала — читать.
// По истечению N секунд программа должна завершаться.

func Solve5() {
	var wg sync.WaitGroup

	ch := make(chan int)

	duration := 10 * time.Second
	done := time.After(duration)

	go writer(ch, &wg, done)

	wg.Add(1)
	go reader(ch, &wg)

	wg.Wait()

	fmt.Println("end")
}

func writer(ch chan int, wg *sync.WaitGroup, done <-chan time.Time) {
	defer wg.Done()
	i := 1
	for {
		select {
		case <-done:
			close(ch)
			return
		default:
			ch <- i
			i++
		}
	}
}

func reader(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Println("read from ch:", v)
	}
}
