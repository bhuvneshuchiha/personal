package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func mimicGoRoutine(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	i := 0
	for {
		<- time.After(1 * time.Second)
		fmt.Println(i)
		i++
		select {
		case <- ctx.Done():
			fmt.Println("STOP")
			return
		}
	}
}


func main() {

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go mimicGoRoutine(&wg, ctx)
	<- time.After(3 * time.Second)
	cancel()
	wg.Wait()

}












