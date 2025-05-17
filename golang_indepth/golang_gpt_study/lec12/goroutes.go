package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func fn1() {
	defer wg.Done()
	fmt.Println("Message from 1st goroutine")
}

func fn2() {
	defer wg.Done()
	fmt.Println("Message from 2nd goroutine")
}

func fn3() {
	defer wg.Done()
	fmt.Println("Message from 3rd goroutine")
}

func main() {
	wg.Add(1)
	go fn1()
	wg.Add(1)
	go fn2()
	wg.Add(1)
	go fn3()

	wg.Wait()
}
