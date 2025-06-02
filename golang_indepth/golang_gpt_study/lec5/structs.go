package main

import (
	"fmt"
)

type Address struct {
	Street  string
	City    string
	ZipCode string
}

type Person struct {
	Name    string
	Age     int
	Address
}

func (p Person) PrintSummary() {
	fmt.Printf("Name: %s\nAge: %d\nStreet: %s\nCity: %s\nZipCode: %s\n",
		p.Name, p.Age, p.Address.Street, p.Address.City, p.Address.ZipCode)
}

func (a *Address) UpdateZip(zip string) {
	a.ZipCode = zip
}

func main() {
	p := Person{
		Name: "Bhuvnesh",
		Age:  25,
		Address: Address{
			Street:  "123 Main St",
			City:    "Delhi",
			ZipCode: "110001",
		},
	}
	p.PrintSummary()

	p.UpdateZip("110002")
	p.PrintSummary()

}
