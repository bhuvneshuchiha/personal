package main

import (
	"fmt"
)

func doLogic(num int) {
	if num%2 != 0 {
		fmt.Println("Number is odd")
	}else {
		fmt.Println("Number is even")
	}
	if num%3 == 0 {
		fmt.Println("Number is divisible by 3")
	}else {
		fmt.Println("Number is not divisible by 3")
	}

	i := 1
	for {
		fmt.Println(i)
		i++
		if i > num {
			break
		}
	}

	switch num {
	case 1:
		fmt.Println("Number is 1")
	case 2:
		fmt.Println("Number is 2")
	case 3:
		fmt.Println("Number is 3")
	case 4:
		fmt.Println("Number is 4")
	case 5:
		fmt.Println("Number is 5")
	default:
		fmt.Println("Number is too big")
	}
}

func main() {
	doLogic(5)
}
