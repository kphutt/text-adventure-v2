package game

import (
	"strings"
	"testing"
)

func TestMovement(t *testing.T) {
	game := createSimpleLayout()
	if game.Player.Location.Name != "Room B" {
		t.Errorf("Expected player to start in Room B, but in %s", game.Player.Location.Name)
	}

	// Test valid movement (go west)
	game.HandleCommand("go west")
	if game.Player.Location.Name != "Room A" {
		t.Errorf("Expected player to be in Room A, but in %s", game.Player.Location.Name)
	}

	// Test moving back (go east)
	game.HandleCommand("go east")
	if game.Player.Location.Name != "Room B" {
		t.Errorf("Expected player to move back to Room B, but in %s", game.Player.Location.Name)
	}

	// Test invalid typed movement
	msg, _ := game.HandleCommand("go north")
	if !strings.Contains(msg, "You can't go that way.") {
		t.Errorf("Expected 'You can't go that way.', but got '%s'", msg)
	}
	if game.Player.Location.Name != "Room B" {
		t.Errorf("Expected player to still be in Room B, but in %s", game.Player.Location.Name)
	}

	// Test WASD movement
	game.HandleCommand("d") // 'd' is east to Room C
	if game.Player.Location.Name != "Room C" {
		t.Errorf("Expected player to be in Room C, but in %s", game.Player.Location.Name)
	}

	// Test invalid single-letter command
	msg, _ = game.HandleCommand("x")
	if !strings.Contains(msg, "I don't understand that command.") {
		t.Errorf("Expected 'I don't understand that command.', but got '%s'", msg)
	}
	if game.Player.Location.Name != "Room C" {
		t.Errorf("Expected player to still be in Room C, but in %s", game.Player.Location.Name)
	}
}

func TestLook(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"] // Move player to the room with the item
	lookResult := game.Look()
	if !strings.Contains(lookResult, "This is Room A.") {
		t.Error("Look command did not return room description.")
	}
	if !strings.Contains(lookResult, "- test_item") {
		t.Error("Look command did not list items.")
	}
	if !strings.Contains(lookResult, "Exits:") {
		t.Error("Look command did not list exits.")
	}
}

func TestTakeAndDrop(t *testing.T) {
	game := createLayoutWithItems()

	// Test taking from an empty room using the 'e' shortcut
	msg, _ := game.HandleCommand("e")
	if !strings.Contains(msg, "There is nothing to take.") {
		t.Errorf("Expected 'There is nothing to take', but got '%s'", msg)
	}

	// Move player to the room with the item
	game.Player.Location = game.AllRooms["Room A"] 

	// Test taking an item that exists
	msg, _ = game.HandleCommand("take test_item")
	if !strings.Contains(msg, "You took the test_item.") {
		t.Errorf("Expected 'You took the test_item', but got '%s'", msg)
	}
	if len(game.Player.Inventory) != 1 || game.Player.Inventory[0].Name != "test_item" {
		t.Error("Player inventory should contain the test_item.")
	}
	if len(game.Player.Location.Items) != 0 {
		t.Error("Room should not contain the item after taking.")
	}

	// Test taking a non-existent item
	msg, _ = game.HandleCommand("take shield")
	if !strings.Contains(msg, "You don't see that here.") {
		t.Errorf("Expected 'You don't see that here', but got '%s'", msg)
	}

	// Test dropping an un-held item
	msg, _ = game.HandleCommand("drop key")
	if !strings.Contains(msg, "You don't have that.") {
		t.Errorf("Expected 'You don't have that', but got '%s'", msg)
	}

	// Test dropping a held item
	msg, _ = game.HandleCommand("drop test_item")
	if !strings.Contains(msg, "You dropped the test_item.") {
		t.Errorf("Expected 'You dropped the test_item', but got '%s'", msg)
	}
	if len(game.Player.Inventory) != 0 {
		t.Error("Player inventory should be empty after dropping.")
	}
	if len(game.Player.Location.Items) != 1 || game.Player.Location.Items[0].Name != "test_item" {
		t.Error("Room should contain the item after dropping.")
	}
}

func TestInventory(t *testing.T) {
	game := createLayoutWithItems()
	msg := game.Inventory()
	if !strings.Contains(msg, "You are not carrying anything.") {
		t.Errorf("Expected 'You are not carrying anything', but got '%s'", msg)
	}

	game.Player.Location = game.AllRooms["Room A"] // Move to room with item
	game.HandleCommand("take test_item")
	msg = game.Inventory()
	if !strings.Contains(msg, "You have the following items:") || !strings.Contains(msg, "- test_item") {
		t.Errorf("Inventory did not list the test_item correctly. Got: %s", msg)
	}
}

func TestUnlock(t *testing.T) {
	// Test unlocking in a room with no locked doors
	gameNoLock := createSimpleLayout()
	msg, _ := gameNoLock.HandleCommand("unlock")
	if !strings.Contains(msg, "There is nothing to unlock here.") {
		t.Errorf("Expected 'There is nothing to unlock here.', but got '%s'", msg)
	}

	// Test the full unlock sequence
	game := createLayoutWithLock()

	// Try to unlock without key
	msg, _ = game.HandleCommand("u")
	if !strings.Contains(msg, "You don't have the key.") {
		t.Errorf("Expected 'You don't have the key', but got '%s'", msg)
	}
	// Verify the door is still locked
	if !game.Player.Location.Exits["east"].Locked {
		t.Error("Door should still be locked.")
	}

	// Go get the key
	game.HandleCommand("go west")  // to Room A
	game.HandleCommand("take key")
	if len(game.Player.Inventory) != 1 {
		t.Fatal("Player should have the key now.")
	}

	// Go back to unlock
	game.HandleCommand("go east")  // back to Room B
	msg, shouldExit := game.HandleCommand("u")
	if shouldExit {
		t.Error("Game should not exit when unlocking a regular door.")
	}
	if !strings.Contains(msg, "You unlocked the door.") {
		t.Errorf("Expected 'You unlocked the door.', but got '%s'", msg)
	}
	if game.Player.Location.Exits["east"].Locked {
		t.Error("Door should be unlocked now.")
	}

	// Test moving through the now unlocked door
	game.HandleCommand("go east")
	if game.Player.Location.Name != "Room C" {
		t.Errorf("Expected to be in Room C, but in %s", game.Player.Location.Name)
	}
}
