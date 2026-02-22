package game

import (
	"strings"
	"testing"
)

func TestWinCondition(t *testing.T) {
	// Setup a game world specifically designed to test the win condition
	game := createLayoutWithWinCondition()

	// 1. Verify initial state
	if game.IsWon {
		t.Fatal("Game should not start in a won state.")
	}

	// 2. Go to the room with the key
	_, shouldExit := game.HandleCommand("go west")
	if shouldExit {
		t.Fatal("Moving to get the key should not end the game.")
	}

	// 3. Take the key
	_, shouldExit = game.HandleCommand("take key")
	if shouldExit {
		t.Fatal("Taking the key should not end the game.")
	}

	// 4. Go back to the locked door
	_, shouldExit = game.HandleCommand("go east")
	if shouldExit {
		t.Fatal("Moving back to the door should not end the game.")
	}

	// 5. Unlock the final door
	msg, shouldExit := game.HandleCommand("u")

	// 6. Assert the win state
	if !shouldExit {
		t.Error("Unlocking the final door should trigger the game to exit.")
	}
	if !game.IsWon {
		t.Error("IsWon flag should be true after unlocking the final door.")
	}
	if !strings.Contains(msg, "You win!") {
		t.Errorf("Expected win message, but got '%s'", msg)
	}
}
