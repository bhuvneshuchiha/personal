package main

import (
	"context"
	"fmt"
	"sync"
)
type ctxKey string

func mimicGo(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	id := ctx.Value(ctxKey("id"))
	fmt.Println("Processing started for the id", id)
}

func main() {
	var wg sync.WaitGroup

	ctx := context.WithValue(context.Background(), ctxKey("id"), 42)

	wg.Add(1)
	go mimicGo(&wg, ctx)

	wg.Wait()
}
