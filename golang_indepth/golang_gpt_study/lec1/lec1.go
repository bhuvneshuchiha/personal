package main

import "fmt"

var a int
var b bool
var c string
const pi = 3.13

func greetTime(name string) {
	fmt.Println("Hello", name)
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
