package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)


type ctxKey string

func mimicGo(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	id := ctx.Value(ctxKey("id"))
	for {
		select {
		case <- ctx.Done():
			return
		default:
			fmt.Println(id)
			time.Sleep(500 * time.Millisecond)

		}
	}
}




func main() {
	ctxInit, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	ctx := context.WithValue(ctxInit, ctxKey("id"), 101)

	var wg sync.WaitGroup
	wg.Add(1)
	go mimicGo(&wg, ctx)
	wg.Wait()

}
