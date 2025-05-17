package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func mimicGo(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	i := 0
	for {
		select {
			case <- time.After(500 * time.Millisecond):
					fmt.Println(i)
					i++
			case <- ctx.Done():
					fmt.Println("Ctx done")
					return
		}
	}
}

func main() {
	var wg sync.WaitGroup

	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	wg.Add(1)
	go mimicGo(&wg, ctx)
	wg.Wait()

}
