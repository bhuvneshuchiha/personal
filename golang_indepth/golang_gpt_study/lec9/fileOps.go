package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Please enter 5 user names")
	file, err := os.OpenFile("prac.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Could not open the file: %s\n", err)
		return
	}
	for range 5 {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		_, err := file.WriteString(text)
		if err != nil {
			fmt.Printf("Error writing to file: %s\n", err)
			return
		}
	}
	file.Close() // close after writing

	// Reopen the file for reading
	file, err = os.Open("prac.txt")
	if err != nil {
		fmt.Printf("Could not open the file for reading: %s\n", err)
		return
	}
	defer file.Close()

	// Reading from the file
	fmt.Println("Reading from the file:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Print(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
	}
}
