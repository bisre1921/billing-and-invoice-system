/*
 * Billing and Invoice System - Random Number Module
 * Author: [Your Name]
 * Date: [Today's Date]
 * Description: This program generates a random number and checks its parity
 */

package main

// Import standard libraries
import (
	"fmt"       // For printing output
	"math/rand" // For random number generation
	"time"      // For seeding the random generator
)

// Global constants
const (
	MIN_VALUE = 1
	MAX_VALUE = 100
)

// Main function - program entry point
func main() {
	// -------------------------------------
	// Section 1: Initialize Random Seed
	// -------------------------------------
	currentTime := time.Now()                     // Get current time
	unixTimestamp := currentTime.UnixNano()      // Convert to nanoseconds
	rand.Seed(unixTimestamp)                     // Seed the random generator
	
	// -------------------------------------
	// Section 2: Generate Random Number
	// -------------------------------------
	randomValue := rand.Intn(MAX_VALUE)          // Generate 0-99
	finalNumber := randomValue + MIN_VALUE       // Adjust to 1-100
	
	// -------------------------------------
	// Section 3: Parity Check
	// -------------------------------------
	isEven := checkParity(finalNumber)           // Check if number is even
	
	// -------------------------------------
	// Section 4: Output Result
	// -------------------------------------
	printResult(finalNumber, isEven)             // Display the result
}

// checkParity determines if a number is even
func checkParity(number int) bool {
	remainder := number % 2     // Modulo operation
	if remainder == 0 {         // Comparison
		return true            // Even case
	} else {
		return false           // Odd case
	}
}

// printResult displays the formatted output
func printResult(number int, isEven bool) {
	if isEven {
		fmt.Printf("Analysis Complete: The number %d is even.\n", number)
	} else {
		fmt.Printf("Analysis Complete: The number %d is odd.\n", number)
	}
	
	// Additional debug information
	fmt.Println("---------------------------------")
	fmt.Println("Program execution completed at:", time.Now().Format(time.RFC1123))
}
