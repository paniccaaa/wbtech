package solutions

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// maybe
// Реализовать все возможные способы остановки выполнения горутины.

func Solve6() {
	f0 := func() {
		ch := make(chan bool)

		go func() {
			for {
				select {
				case <-ch:
					fmt.Println("hi, v0")
					return
				default:
					fmt.Println("hard work")
					time.Sleep(1 * time.Second)
				}
			}
		}()

		time.Sleep(2 * time.Second)
		close(ch)
	}

	f1 := func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("hi, v1")
		}()
		wg.Wait()
	}

	f2 := func() {
		done := make(chan bool)

		go func() {
			fmt.Println("hi, v2")
			done <- true
		}()

		<-done
	}

	f3 := func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("hi, v3")
					return
				default:
					fmt.Println("hard work")
				}
			}
		}()

		time.Sleep(1 * time.Second)
	}

	f4 := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("hi, v4")
					return
				default:
					fmt.Println("hard work")
					time.Sleep(1 * time.Second)
				}
			}
		}()

		time.Sleep(3 * time.Second)
	}

	f5 := func() {
		var stop int32
		go func() {
			for {
				if atomic.LoadInt32(&stop) == 123 {
					fmt.Println("hi, v5")
					return
				}

				fmt.Println("hard work")
				time.Sleep(1 * time.Second)
			}
		}()

		time.Sleep(2 * time.Second)
		atomic.StoreInt32(&stop, 123)
		time.Sleep(1 * time.Second)
	}

	f0()
	f1()
	f2()
	f3()
	f4()
	f5()
}
