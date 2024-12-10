package solutions

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// DONE
// Реализовать постоянную запись данных в канал (главный поток).
// Реализовать набор из N воркеров, которые читают произвольные данные из канала и выводят в stdout.
// Необходима возможность выбора количества воркеров при старте.
// Программа должна завершаться по нажатию Ctrl+C.
// Выбрать и обосновать способ завершения работы всех воркеров

func Solve4() {
	var numWorkers = flag.Int("worker", 1, "num of workers")
	flag.Parse()

	var wg sync.WaitGroup

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)

	ch := make(chan int)
	i := -1212831

	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()

	wg.Add(*numWorkers)
	for w := range *numWorkers {
		go func() {
			defer wg.Done()
			for i := range ch {
				fmt.Printf("hello from %d worker receive i=%d\n", w, i)
			}
		}()
	}

	wg.Wait()
}
