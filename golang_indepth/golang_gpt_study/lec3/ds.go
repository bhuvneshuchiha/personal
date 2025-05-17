package main

import "fmt"

func main() {
	//Declare the slice
	sl1 := []int{1,2,3,4,5}

	//Append the elements
	sl1 = append(sl1, 6)
	sl1 = append(sl1, 7)

	//Print out the results
	fmt.Println("The length of the slice is", len(sl1))
	fmt.Println("The capacity of the slice is", cap(sl1))
	fmt.Println("The slice is :-",sl1)

	//Initialize the map
	mp := make(map[string]float64)
	//Add values to a map
	mp["Bob"] = 98.2
	mp["Alice"] = 97.0
	mp["Teej"] = 88.0

	//Update the key
	mp["Teej"] = 80.0
	//Delete the key
	delete(mp, "Teej")

	//Print the map
	for k, v := range mp {
		fmt.Println(k,v)
	}

}

