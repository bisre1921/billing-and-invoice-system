package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 1 and 100
	randomNum := rand.Intn(100) + 1

	// Check if the number is even or odd
	if randomNum%2 == 0 {
		fmt.Printf("%d is even!\n", randomNum)
	} else {
		fmt.Printf("%d is odd!\n", randomNum)
	}
}
