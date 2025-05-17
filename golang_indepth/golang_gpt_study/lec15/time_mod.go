package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ch := make(chan int , 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 5{
			<- time.After(1 * time.Second)
			ch <- i
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case val, ok := <-ch:
				if !ok {
					return
				}
				fmt.Println("received from channel", val)
				case <- time.After(2 * time.Second):
					fmt.Println("timeout")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for tick := range time.Tick(2 * time.Second) {
			fmt.Println("Tick at", tick)
			i++
			if i == 5 {
				break
			}
		}
	}()
	wg.Wait()
}
