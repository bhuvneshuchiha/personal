package main

import (
	"fmt"
	"strconv"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered safely from panic")
			fmt.Println("Reason for panic :-", r)
		}
	}()
	fmt.Println("Please enter a number")
	var inputString string
	fmt.Scanln(&inputString)
	finalInt, _ := strconv.Atoi(inputString)

	fmt.Println("Performing division")
	if finalInt == 0 {
		panic("Division by zero error")
	}
	output := 100/finalInt
	fmt.Println(output)

}
