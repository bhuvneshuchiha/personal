package main

import "fmt"

func mayPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from:", r)
		}
	}()
	panic("unexpected crash!")
}

func main() {
	fmt.Println("Before panic")
	mayPanic()
	fmt.Println("After panic")
}

