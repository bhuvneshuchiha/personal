package main

import (
	"fmt"
	"strings"
)


func main() {
	str := "madam1"
	str_1 := []string{}
	temp := ""

	for i := len(str)-1; i >= 0; i-- {
		str_1 = append(str_1, string(str[i]))
	}

	temp = strings.Join(str_1, "")

	if str == temp {
		fmt.Println("valid palindrom")
	}else {
		fmt.Println("not a valid palindrom")
	}
}
