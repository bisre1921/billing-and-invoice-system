package tests

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

var (
	short = flag.Bool("short", false, "skip long-running tests")
)

func TestMain(m *testing.M) {
	// Parse flags to ensure testing.Short() works correctly
	flag.Parse()

	// Setup before running tests
	fmt.Println("Setting up tests...")

	// Load environment variables from .env file
	// Try multiple paths to find the .env file
	envPaths := []string{
		".env",
		"../.env",
		"../../.env",
		filepath.Join("..", ".env"),
	}

	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			fmt.Printf("Loaded environment from: %s\n", path)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		log.Println("Warning: Could not load .env file. Using environment variables directly.")
	}
		// No database connection needed for unit tests

	// Run tests
	exitCode := m.Run()
	// Cleanup
	fmt.Println("Cleaning up after tests...")

	os.Exit(exitCode)
}
