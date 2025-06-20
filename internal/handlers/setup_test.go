package handlers

import (
	"testing"
)

func init() {
	// Enable test mode for all handler tests
	IsTestMode = true
}

// TestHandlersSetup ensures basic test environment setup
// This test ensures that the handler package is initialized correctly for tests
func TestHandlersSetupBasic(t *testing.T) {
	if !IsTestMode {
		t.Error("Test mode not properly set")
	}
}
