package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// This is just a wrapper to direct users to the tests package
	t.Log("Main tests are now located in the 'tests' package.")
	t.Log("Please run 'go test ./tests/...' instead.")
}
