package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)


func mimicGo(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for {
		select {
		case <- time.After(500 * time.Millisecond):
				fmt.Println(i)
				i++
		case <- ctx.Done():
			fmt.Println("Ctx.Done() triggered")
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, _ := context.WithTimeout(context.Background(), 2 * time.Second)
	wg.Add(1)
	go mimicGo(ctx, &wg)
	wg.Wait()
}
