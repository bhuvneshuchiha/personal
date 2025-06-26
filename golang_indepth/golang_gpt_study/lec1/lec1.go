package main

import (
	"fmt"
	"log/slog"
)

var a int
var b bool
var c string
const pi = 3.13

func greetTime(name string) {
	fmt.Println("Hello", name)
	slog.Info("Hi")
}

func addNums(a, b int) int {
	return a + b
}

func main() {
	greetTime("bhuvnesh")
	sum := addNums(1,2)
	fmt.Println(sum)
	fmt.Println("a:",a,"b:",b,"c:",c)
}
