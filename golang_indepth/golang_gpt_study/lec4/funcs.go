package main

import (
	"fmt"
)

type testUser struct {
	name string
	age int
}

//Method with value reciever
func (t testUser) incrementAge(i int) testUser{
	return testUser {
		name: t.name,
		age: t.age + i,
	}
}

//Method with pointer reciever
func(t *testUser) incrementPointerAge(i int){
	t.age = t.age + i
}



//Function accepting n number of arguments
func mult(nums ...int) int {
	total := 1
	for _, v := range nums {
		total *= v
	}
	return total
}

//Closure
func demoClose(a, b int) func(a, b int) int{
	sum := 0
	return func(a, b int) int{
		sum = a + b
		return sum
	}
}

func main() {
	tot := mult(1,2,3,4)
	fmt.Println(tot)

	//Struct init to demonstrate value and ptr receiver methods.
	user := testUser {
		name: "Bhuvnesh",
		age: 20,
	}

	newUser := user.incrementAge(5)
	fmt.Println("New user is", newUser)
	fmt.Println("Old user is", user)

	fmt.Println("*******************")
	user.incrementPointerAge(10)
	fmt.Println("Old user now becomes", user)
}
