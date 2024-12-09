package solutions

import (
	"fmt"
	"sync"
)

// DONE
// Реализовать конкурентную запись данных в map.

type cmap struct {
	mu sync.Mutex
	m  map[int]int
}

func (c *cmap) set(k, v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k] = v
	fmt.Println("key =", k, "v =", v)
}

func Solve7() {
	var wg sync.WaitGroup
	example := cmap{m: make(map[int]int)}

	wg.Add(100)
	for i := range 100 {
		go func() {
			defer wg.Done()
			example.set(i, i+10)
		}()
	}
	wg.Wait()
}

// worker pool
func Solve7_1() {
	var wg sync.WaitGroup

	example := cmap{m: make(map[int]int)}
	ch := make(chan int)

	go func() {
		for i := range 100 {
			ch <- i
		}
		close(ch)
	}()

	wg.Add(5)
	for range 5 {
		go func() {
			defer wg.Done()
			for v := range ch {
				example.set(v, v+20)
			}
		}()
	}
	wg.Wait()

}
