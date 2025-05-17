package testing_challenge

import "errors"

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("Cannot divide by Zero")
	}else {
		return a/b, nil
	}
}

func Subtract(a, b int) int {
	if a > b {
		return a - b
	}else {
		return b - a
	}
}
