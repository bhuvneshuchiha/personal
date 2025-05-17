package main

import (
	"fmt"
	"sync"
)

func main() {

	var mu sync.Mutex
	var wg sync.WaitGroup

	counter := 0

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}
