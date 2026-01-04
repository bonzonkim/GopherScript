package main

import (
	"fmt"
	"log"
)

// AddNumbers adds two integers and returns the result.
func AddNumbers(a, b int) int {
	return a + b
}

func main() {
	// Call the AddNumbers function.
	result := AddNumbers(2, 3)

	// Print the result.
	fmt.Println(result)

	//Example of handling errors, but AddNumbers doesn't error, so just an example
	_, err := trySomethingThatErrors()
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}

// trySomethingThatErrors is an example function that returns an error
func trySomethingThatErrors() (int, error) {
	return 0, fmt.Errorf("this is a simulated error")
}
