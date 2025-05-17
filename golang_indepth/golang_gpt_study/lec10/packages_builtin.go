package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	input := " 12 , 7,  15,abc, 9 , -3 ,   "
	s1 := strings.Split(input, ",")
	total := 0
	for _, v := range s1 {
		v = strings.TrimSpace(v)
		v, er := strconv.Atoi(v)
		if er != nil {
			fmt.Println("Inside error")
			fmt.Println(fmt.Errorf("Cannot convert this one %s", er))
			continue
		}
		total += v
	}
	fmt.Println(total)
}
