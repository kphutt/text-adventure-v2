package game

import (
	"strings"
	"testing"
)

func TestMovement(t *testing.T) {
	game := NewGame()
	if game.Player.Location.Name != "Hall" {
		t.Errorf("Expected player to start in Hall, but in %s", game.Player.Location.Name)
	}

	// Test valid movement
	game.HandleCommand("go north")
	if game.Player.Location.Name != "Dungeon" {
		t.Errorf("Expected player to be in Dungeon, but in %s", game.Player.Location.Name)
	}

	// Test invalid movement
	msg, _ := game.HandleCommand("go north")
	if !strings.Contains(msg, "The door is locked.") {
		t.Errorf("Expected 'The door is locked', but got '%s'", msg)
	}
	if game.Player.Location.Name != "Dungeon" {
		t.Errorf("Expected player to still be in Dungeon, but in %s", game.Player.Location.Name)
	}

	// Test WASD movement
	game.HandleCommand("s")
	if game.Player.Location.Name != "Hall" {
		t.Errorf("Expected player to be in Hall, but in %s", game.Player.Location.Name)
	}
}

func TestLook(t *testing.T) {
	game := NewGame()
	lookResult := game.Look()
	if !strings.Contains(lookResult, "You are in a long, dark hall.") {
		t.Error("Look command did not return room description.")
	}
	if !strings.Contains(lookResult, "- sword") {
		t.Error("Look command did not list items.")
	}
	if !strings.Contains(lookResult, "Exits:") {
		t.Error("Look command did not list exits.")
	}
}

func TestTakeAndDrop(t *testing.T) {
	game := NewGame()

	// Test taking an item that exists
	msg, _ := game.HandleCommand("take sword")
	if !strings.Contains(msg, "You took the sword.") {
		t.Errorf("Expected 'You took the sword', but got '%s'", msg)
	}
	if len(game.Player.Inventory) != 1 || game.Player.Inventory[0].Name != "sword" {
		t.Error("Player inventory should contain the sword.")
	}
	if len(game.Player.Location.Items) != 0 {
		t.Error("Room should not contain the sword after taking.")
	}

	// Test taking a non-existent item
	msg, _ = game.HandleCommand("take shield")
	if !strings.Contains(msg, "You don't see that here.") {
		t.Errorf("Expected 'You don't see that here', but got '%s'", msg)
	}

	// Test dropping an item
	msg, _ = game.HandleCommand("drop sword")
	if !strings.Contains(msg, "You dropped the sword.") {
		t.Errorf("Expected 'You dropped the sword', but got '%s'", msg)
	}
	if len(game.Player.Inventory) != 0 {
		t.Error("Player inventory should be empty after dropping.")
	}
	if len(game.Player.Location.Items) != 1 || game.Player.Location.Items[0].Name != "sword" {
		t.Error("Room should contain the sword after dropping.")
	}
}

func TestInventory(t *testing.T) {
	game := NewGame()
	msg := game.Inventory()
	if !strings.Contains(msg, "You are not carrying anything.") {
		t.Errorf("Expected 'You are not carrying anything', but got '%s'", msg)
	}

	game.HandleCommand("take sword")
	msg = game.Inventory()
	if !strings.Contains(msg, "You have the following items:") || !strings.Contains(msg, "- sword") {
		t.Errorf("Inventory did not list the sword correctly. Got: %s", msg)
	}
}

func TestUnlock(t *testing.T) {
	game := NewGame()

	// Try to unlock without key
	game.HandleCommand("go north") // to Dungeon
	msg, _ := game.HandleCommand("unlock")
	if game.IsWon {
		t.Error("Game should not be won when unlocking fails.")
	}
	if !strings.Contains(msg, "You don't have the key.") {
		t.Errorf("Expected 'You don't have the key', but got '%s'", msg)
	}

	// Go get the key
	game.HandleCommand("go south") // back to Hall
	game.HandleCommand("go east")  // to Closet
	game.HandleCommand("take key")

	// Go back to unlock
	game.HandleCommand("go west")  // back to Hall
	game.HandleCommand("go north") // to Dungeon
	msg, _ = game.HandleCommand("unlock")
	if !game.IsWon {
		t.Error("Expected the game to be won")
	}
	if !strings.Contains(msg, "You unlocked the door!") {
		t.Errorf("Expected 'You unlocked the door!', but got '%s'", msg)
	}
}
