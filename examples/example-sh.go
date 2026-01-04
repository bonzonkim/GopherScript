package main

import (
	"fmt"
	"strconv"
)

// AddNumbers takes two integers as strings and returns their sum as an integer.
// It also returns an error if the input strings cannot be converted to integers.
func AddNumbers(a, b string) (int, error) {
	num1, err := strconv.Atoi(a)
	if err != nil {
		return 0, fmt.Errorf("invalid input for a: %w", err)
	}
	num2, err := strconv.Atoi(b)
	if err != nil {
		return 0, fmt.Errorf("invalid input for b: %w", err)
	}
	return num1 + num2, nil
}

func main() {
	// Call AddNumbers with string representations of 2 and 3.
	result, err := AddNumbers("2", "3")

	// Handle potential errors.
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result if there are no errors.
	fmt.Println(result)
}
