package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func sendChanDemo(ch chan int) {
	defer wg.Done()
	for i := range 3 {
		ch <- i
	}
	close(ch)
}

func receiveChanDemo(ch chan int) {
	defer wg.Done()
	for val := range ch {
		fmt.Println(val)
	}
}

func main() {
	ch := make(chan int, 3)
	wg.Add(1)
	go sendChanDemo(ch)
	wg.Add(1)
	go receiveChanDemo(ch)
	wg.Wait()

}

