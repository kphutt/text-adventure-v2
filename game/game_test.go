package game

import (
	"strings"
	"testing"
)

func TestMovement(t *testing.T) {
	game := createSimpleLayout()
	if game.Turns != 0 {
		t.Errorf("Game should start with 0 turns, but has %d", game.Turns)
	}

	// Test valid movement (go west)
	game.HandleCommand("go west")
	if game.Player.Location.Name != "Room A" {
		t.Errorf("Expected player to be in Room A, but in %s", game.Player.Location.Name)
	}
	if game.Turns != 1 {
		t.Errorf("Valid move should increment turns. Expected 1, got %d", game.Turns)
	}

	// Test moving back (go east)
	game.HandleCommand("go east")
	if game.Player.Location.Name != "Room B" {
		t.Errorf("Expected player to move back to Room B, but in %s", game.Player.Location.Name)
	}
	if game.Turns != 2 {
		t.Errorf("Valid move should increment turns. Expected 2, got %d", game.Turns)
	}

	// Test invalid typed movement
	msg, _ := game.HandleCommand("go north")
	if !strings.Contains(msg, "You can't go that way.") {
		t.Errorf("Expected 'You can't go that way.', but got '%s'", msg)
	}
	if game.Turns != 2 {
		t.Errorf("Invalid move should not increment turns. Expected 2, got %d", game.Turns)
	}

	// Test WASD movement
	game.HandleCommand("d") // 'd' is east to Room C
	if game.Player.Location.Name != "Room C" {
		t.Errorf("Expected player to be in Room C, but in %s", game.Player.Location.Name)
	}
	if game.Turns != 3 {
		t.Errorf("Valid WASD move should increment turns. Expected 3, got %d", game.Turns)
	}

	// Test invalid single-letter command
	msg, _ = game.HandleCommand("x")
	if !strings.Contains(msg, "I don't understand that command.") {
		t.Errorf("Expected 'I don't understand that command.', but got '%s'", msg)
	}
	if game.Turns != 3 {
		t.Errorf("Unknown command should not increment turns. Expected 3, got %d", game.Turns)
	}
}

func TestLook(t *testing.T) {
	game := createLayoutWithItems()
	game.HandleCommand("look")
	if game.Turns != 0 {
		t.Errorf("Look command should not increment turns, but got %d", game.Turns)
	}
}

func TestTakeAndDrop(t *testing.T) {
	game := createLayoutWithItems()

	// Test taking from an empty room using the 'e' shortcut (should not increment turn)
	game.HandleCommand("e")
	if game.Turns != 0 {
		t.Errorf("Failed 'e' command should not increment turns, but got %d", game.Turns)
	}

	// Move player to the room with the item
	game.Player.Location = game.AllRooms["Room A"]

	// Test taking an item that exists (should increment turn)
	game.HandleCommand("take test_item")
	if game.Turns != 1 {
		t.Errorf("Successful take should increment turns. Expected 1, got %d", game.Turns)
	}

	// Test taking a non-existent item (should not increment turn)
	game.HandleCommand("take shield")
	if game.Turns != 1 {
		t.Errorf("Failed take should not increment turns. Expected 1, got %d", game.Turns)
	}

	// Test dropping an un-held item (should not increment turn)
	game.HandleCommand("drop key")
	if game.Turns != 1 {
		t.Errorf("Failed drop should not increment turns. Expected 1, got %d", game.Turns)
	}

	// Test dropping a held item (should increment turn)
	game.HandleCommand("drop test_item")
	if game.Turns != 2 {
		t.Errorf("Successful drop should increment turns. Expected 2, got %d", game.Turns)
	}
}

func TestInventory(t *testing.T) {
	game := createLayoutWithItems()
	game.HandleCommand("inventory")
	if game.Turns != 0 {
		t.Errorf("Inventory command should not increment turns, but got %d", game.Turns)
	}

	game.Player.Location = game.AllRooms["Room A"]
	game.HandleCommand("take test_item")
	if game.Turns != 1 {
		t.Errorf("Turns should be 1 after taking item.")
	}
	game.HandleCommand("inventory")
	if game.Turns != 1 {
		t.Errorf("Inventory command should not increment turns after taking. Expected 1, got %d", game.Turns)
	}
}

func TestUnlock(t *testing.T) {
	// Test unlocking in a room with no locked doors (should not increment turn)
	gameNoLock := createSimpleLayout()
	gameNoLock.HandleCommand("unlock")
	if gameNoLock.Turns != 0 {
		t.Errorf("Failed unlock should not increment turns, but got %d", gameNoLock.Turns)
	}

	// Test the full unlock sequence
	game := createLayoutWithLock()

	// Try to unlock without key (should not increment turn)
	game.HandleCommand("u")
	if game.Turns != 0 {
		t.Errorf("Failed unlock (no key) should not increment turns, but got %d", game.Turns)
	}

	// Go get the key
	game.HandleCommand("go west")  // Turn 1
	game.HandleCommand("take key") // Turn 2
	if game.Turns != 2 {
		t.Errorf("Expected 2 turns after moving and taking key, got %d", game.Turns)
	}

	// Go back to unlock
	game.HandleCommand("go east") // Turn 3
	game.HandleCommand("u")       // Turn 4 (successful unlock)
	if game.Turns != 4 {
		t.Errorf("Expected 4 turns after successful unlock, got %d", game.Turns)
	}

	// Test moving through the now unlocked door
	game.HandleCommand("go east") // Turn 5
	if game.Player.Location.Name != "Room C" {
		t.Errorf("Expected to be in Room C, but in %s", game.Player.Location.Name)
	}
	if game.Turns != 5 {
		t.Errorf("Expected 5 turns after final move, got %d", game.Turns)
	}
}

func TestScore_EmptyInventoryOneRoom(t *testing.T) {
	game := createSimpleLayout()
	score := game.Score()
	if score != 5 {
		t.Errorf("Expected score 5 (one room visited), got %d", score)
	}
}

func TestScore_AfterVisitingMultipleRooms(t *testing.T) {
	game := createSimpleLayout()
	// Start in Room B (1 room visited = 5)
	if game.Score() != 5 {
		t.Errorf("Expected score 5 at start, got %d", game.Score())
	}
	// Move to Room A (2 rooms visited = 10)
	game.HandleCommand("go west")
	if game.Score() != 10 {
		t.Errorf("Expected score 10 after visiting 2 rooms, got %d", game.Score())
	}
	// Move back to Room B (still 2 unique rooms = 10)
	game.HandleCommand("go east")
	if game.Score() != 10 {
		t.Errorf("Expected score 10 after revisiting Room B, got %d", game.Score())
	}
	// Move to Room C (3 rooms visited = 15)
	game.HandleCommand("go east")
	if game.Score() != 15 {
		t.Errorf("Expected score 15 after visiting 3 rooms, got %d", game.Score())
	}
}

func TestScore_Command(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]
	game.HandleCommand("take test_item")

	msg, shouldExit := game.HandleCommand("score")
	if shouldExit {
		t.Error("Score command should not exit the game")
	}
	if !strings.Contains(msg, "15") {
		t.Errorf("Expected score message to contain '15', got '%s'", msg)
	}
	// Score command should not increment turns
	if game.Turns != 1 {
		t.Errorf("Score command should not increment turns. Expected 1, got %d", game.Turns)
	}
}

func TestScore_WithInventoryItems(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]
	game.HandleCommand("take test_item")
	score := game.Score()
	// 10 points for 1 item + 5 points for 1 room visited = 15
	if score != 15 {
		t.Errorf("Expected score 15 (1 item + 1 room), got %d", score)
	}
}

func TestUnlock_WinCondition(t *testing.T) {
	game := createLayoutWithWinCondition()

	// Go get the key
	game.HandleCommand("go west") // to Room A
	game.HandleCommand("take key")

	// Go back to unlock the final door
	game.HandleCommand("go east") // back to Room B
	msg, shouldExit := game.HandleCommand("u")

	if !shouldExit {
		t.Error("Game should exit when unlocking the treasure room door.")
	}
	if !game.IsWon {
		t.Error("Game IsWon flag should be true after winning.")
	}
	if !strings.Contains(msg, "You win!") {
		t.Errorf("Expected win message, but got '%s'", msg)
	}
}
