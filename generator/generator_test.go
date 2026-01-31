package generator

import (
	"testing"
)

// TestGenerate_Success ensures that the generator can produce a valid world without errors.
func TestGenerate_Success(t *testing.T) {
	for i := 0; i < 20; i++ { // Run multiple times for confidence
		config := DefaultConfig()
		startRoom, err := Generate(config)

		if err != nil {
			t.Fatalf("Generate() failed on iteration %d: %v", i, err)
		}
		if startRoom == nil {
			t.Fatalf("Generate() returned a nil startRoom on iteration %d", i)
		}
	}
}

// TestGenerate_ConfigValidation ensures the generator fails gracefully with bad config.
func TestGenerate_ConfigValidation(t *testing.T) {
	config := DefaultConfig()
	// Create an impossible configuration
	config.NumberOfRooms = len(config.RoomNamePool) + 1

	_, err := Generate(config)
	if err == nil {
		t.Fatal("Generate() should have failed with an impossible config, but it did not.")
	}
}
