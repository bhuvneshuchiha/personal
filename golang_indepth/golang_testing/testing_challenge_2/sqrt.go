package testing_challenge_2
// package main

import (
	"errors"
	"fmt"
)

func Sqrt(a int) (int, error) {

	if a < 0 {
		return 0, errors.New("Cannot compute square root of a negative number")
	}else {
		// Binary search
		start := 0
		end := a
		mid := (start + end)/2

		for i := start; i < end; i++ {
			if i * i == a {
				return i, nil
			}else if i * i < a {
				start = mid
			}else {
				end = mid
			}
			mid = (start + end)/2
		}

		return start, nil
	}
}


func main() {

	sqrtAns, err := Sqrt(-144)
	if err != nil {
		fmt.Println(err)
		return
	}else {
		fmt.Println(sqrtAns)
	}
}










