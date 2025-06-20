package handlers

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Run all tests
	result := m.Run()
	
	// Exit with the test result code
	os.Exit(result)
}
