package solutions

import (
	"fmt"
	"sync"
)

// DONE
// Реализовать структуру-счетчик, которая будет инкрементироваться в конкурентной среде.
// По завершению программа должна выводить итоговое значение счетчика.

type counter struct {
	c  int
	mu sync.Mutex
}

func (c *counter) inc(i int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c += i
}

func Solve18() {
	var wg sync.WaitGroup
	c := &counter{c: 12}

	wg.Add(100)
	for i := range 100 {
		go func(i int) {
			defer wg.Done()
			c.inc(i)
		}(i)
	}
	wg.Wait()

	fmt.Println(c.c)
}
