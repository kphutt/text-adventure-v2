package game

import (
	"strings"
	"testing"
	"text-adventure-v2/world"
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

// --- GetAllRooms tests ---

func TestGetAllRooms_FindsAllConnectedRooms(t *testing.T) {
	// Build a 4-room graph: A--B--C, B--D
	a := createRoomWithExits("A")
	b := createRoomWithExits("B")
	c := createRoomWithExits("C")
	d := createRoomWithExits("D")
	linkRooms(a, "east", b, "west")
	linkRooms(b, "east", c, "west")
	linkRooms(b, "south", d, "north")

	rooms := make(map[string]*world.Room)
	GetAllRooms(a, rooms)

	if len(rooms) != 4 {
		t.Errorf("Expected 4 rooms, got %d", len(rooms))
	}
	for _, name := range []string{"A", "B", "C", "D"} {
		if _, ok := rooms[name]; !ok {
			t.Errorf("Missing room %s", name)
		}
	}
}

func TestGetAllRooms_SingleRoom(t *testing.T) {
	a := createRoomWithExits("A")
	rooms := make(map[string]*world.Room)
	GetAllRooms(a, rooms)

	if len(rooms) != 1 {
		t.Errorf("Expected 1 room, got %d", len(rooms))
	}
}

func TestGetAllRooms_HandlesLoops(t *testing.T) {
	// A--B--C--A (circular)
	a := createRoomWithExits("A")
	b := createRoomWithExits("B")
	c := createRoomWithExits("C")
	linkRooms(a, "east", b, "west")
	linkRooms(b, "east", c, "west")
	linkRooms(c, "east", a, "west") // creates a loop

	rooms := make(map[string]*world.Room)
	GetAllRooms(a, rooms)

	if len(rooms) != 3 {
		t.Errorf("Expected 3 rooms (no infinite loop), got %d", len(rooms))
	}
}

// --- Look output tests ---

func TestLook_ShowsDescriptionAndExits(t *testing.T) {
	game := createSimpleLayout()
	msg, _ := game.HandleCommand("look")

	// Room B has exits west and east
	if !strings.Contains(msg, "This is Room B.") {
		t.Errorf("Look should show room description, got: %s", msg)
	}
	if !strings.Contains(msg, "Exits:") {
		t.Errorf("Look should show exits header, got: %s", msg)
	}
	if !strings.Contains(msg, "- east") {
		t.Errorf("Look should list east exit, got: %s", msg)
	}
	if !strings.Contains(msg, "- west") {
		t.Errorf("Look should list west exit, got: %s", msg)
	}
}

func TestLook_ShowsItems(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]

	msg, _ := game.HandleCommand("l") // shortcut

	if !strings.Contains(msg, "You see the following items:") {
		t.Errorf("Look should list items when present, got: %s", msg)
	}
	if !strings.Contains(msg, "- test_item") {
		t.Errorf("Look should show item name, got: %s", msg)
	}
}

func TestLook_ExitsSorted(t *testing.T) {
	game := createSimpleLayout()
	msg := game.Look()

	eastIdx := strings.Index(msg, "- east")
	westIdx := strings.Index(msg, "- west")
	if eastIdx < 0 || westIdx < 0 {
		t.Fatalf("Expected both east and west in output, got: %s", msg)
	}
	if eastIdx > westIdx {
		t.Error("Exits should be sorted alphabetically (east before west)")
	}
}

// --- Inventory output tests ---

func TestInventory_EmptyMessage(t *testing.T) {
	game := createSimpleLayout()
	msg, _ := game.HandleCommand("i")

	if msg != "You are not carrying anything." {
		t.Errorf("Expected empty inventory message, got: %s", msg)
	}
}

func TestInventory_ListsItems(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]
	game.HandleCommand("take test_item")

	msg, _ := game.HandleCommand("inventory")

	if !strings.Contains(msg, "You have the following items:") {
		t.Errorf("Expected inventory header, got: %s", msg)
	}
	if !strings.Contains(msg, "- test_item") {
		t.Errorf("Expected item in inventory list, got: %s", msg)
	}
}

// --- Take behavior tests ---

func TestTake_AutoPicksSingleItem(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"] // has exactly one item

	msg, success := game.Take("")
	if !success {
		t.Error("Take with empty name should auto-pick the single item")
	}
	if !strings.Contains(msg, "You took the test_item.") {
		t.Errorf("Expected auto-pick message, got: %s", msg)
	}
	if len(game.Player.Inventory) != 1 || game.Player.Inventory[0].Name != "test_item" {
		t.Error("Item should be in player's inventory after auto-pick")
	}
	if len(game.Player.Location.Items) != 0 {
		t.Error("Item should be removed from room after auto-pick")
	}
}

func TestTake_AsksWhenMultipleItems(t *testing.T) {
	game := createSimpleLayout()
	game.Player.Location.Items = []*world.Item{
		{Name: "sword", Description: "A sword."},
		{Name: "shield", Description: "A shield."},
	}

	msg, success := game.Take("")
	if success {
		t.Error("Take with empty name and multiple items should fail")
	}
	if msg != "What do you want to take?" {
		t.Errorf("Expected disambiguation prompt, got: %s", msg)
	}
}

func TestTake_CaseInsensitive(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]

	msg, success := game.Take("TEST_ITEM")
	if !success {
		t.Error("Take should be case-insensitive")
	}
	if !strings.Contains(msg, "You took the test_item.") {
		t.Errorf("Expected success message, got: %s", msg)
	}
}

func TestTake_ItemNotFound(t *testing.T) {
	game := createSimpleLayout()
	msg, success := game.Take("nonexistent")
	if success {
		t.Error("Take should fail for nonexistent item")
	}
	if msg != "You don't see that here." {
		t.Errorf("Expected not-found message, got: %s", msg)
	}
}

// --- HandleCommand quit/help/e tests ---

func TestHandleCommand_Quit(t *testing.T) {
	game := createSimpleLayout()
	msg, shouldExit := game.HandleCommand("quit")

	if !shouldExit {
		t.Error("Quit should signal exit")
	}
	if msg != "Goodbye!" {
		t.Errorf("Expected goodbye message, got: %s", msg)
	}
	if game.Turns != 0 {
		t.Error("Quit should not increment turns")
	}
}

func TestHandleCommand_QuitShortcut(t *testing.T) {
	game := createSimpleLayout()
	_, shouldExit := game.HandleCommand("q")

	if !shouldExit {
		t.Error("'q' shortcut should signal exit")
	}
}

func TestHandleCommand_Help(t *testing.T) {
	game := createSimpleLayout()
	msg, shouldExit := game.HandleCommand("help")

	if shouldExit {
		t.Error("Help should not signal exit")
	}
	// Verify help contains the key commands players need to know
	for _, cmd := range []string{"w,a,s,d", "take", "drop", "unlock", "score", "quit"} {
		if !strings.Contains(msg, cmd) {
			t.Errorf("Help should mention '%s', got: %s", cmd, msg)
		}
	}
	if game.Turns != 0 {
		t.Error("Help should not increment turns")
	}
}

func TestHandleCommand_HelpShortcut(t *testing.T) {
	game := createSimpleLayout()
	msg, _ := game.HandleCommand("h")

	if !strings.Contains(msg, "take") {
		t.Errorf("'h' shortcut should show help, got: %s", msg)
	}
}

func TestHandleCommand_E_TakesFirstItem(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]

	msg, _ := game.HandleCommand("e")

	if !strings.Contains(msg, "You took the test_item.") {
		t.Errorf("'e' should take the first item in the room, got: %s", msg)
	}
	if len(game.Player.Inventory) != 1 {
		t.Error("Item should be in inventory after 'e' command")
	}
}

func TestHandleCommand_E_NothingToTake(t *testing.T) {
	game := createSimpleLayout() // no items in any room
	msg, _ := game.HandleCommand("e")

	if msg != "There is nothing to take." {
		t.Errorf("'e' in empty room should say nothing to take, got: %s", msg)
	}
}

// --- Move tests for locked door message ---

func TestMove_LockedDoorMessage(t *testing.T) {
	game := createLayoutWithLock()
	msg, success := game.Move("east")

	if success {
		t.Error("Moving through locked door should fail")
	}
	if msg != "The door is locked." {
		t.Errorf("Expected locked door message, got: %s", msg)
	}
}

// --- Drop behavior tests ---

func TestDrop_EmptyName(t *testing.T) {
	game := createSimpleLayout()
	msg, success := game.Drop("")
	if success {
		t.Error("Drop with empty name should fail")
	}
	if msg != "What do you want to drop?" {
		t.Errorf("Expected prompt, got: %s", msg)
	}
}

func TestDrop_CaseInsensitive(t *testing.T) {
	game := createLayoutWithItems()
	game.Player.Location = game.AllRooms["Room A"]
	game.HandleCommand("take test_item")

	msg, success := game.Drop("TEST_ITEM")
	if !success {
		t.Error("Drop should be case-insensitive")
	}
	if !strings.Contains(msg, "You dropped the test_item.") {
		t.Errorf("Expected drop message, got: %s", msg)
	}
}

func TestDrop_NotHeld(t *testing.T) {
	game := createSimpleLayout()
	msg, success := game.Drop("sword")
	if success {
		t.Error("Drop should fail for items not in inventory")
	}
	if msg != "You don't have that." {
		t.Errorf("Expected not-held message, got: %s", msg)
	}
}

// --- Unlock message tests ---

func TestUnlock_NoLockedDoorMessage(t *testing.T) {
	game := createSimpleLayout()
	msg, _, _ := game.Unlock()
	if msg != "There is nothing to unlock here." {
		t.Errorf("Expected no-locked-door message, got: %s", msg)
	}
}

func TestUnlock_NoKeyMessage(t *testing.T) {
	game := createLayoutWithLock()
	msg, _, _ := game.Unlock()
	if msg != "You don't have the key." {
		t.Errorf("Expected no-key message, got: %s", msg)
	}
}

func TestUnlock_SuccessMessage(t *testing.T) {
	game := createLayoutWithLock()
	game.HandleCommand("go west")
	game.HandleCommand("take key")
	game.HandleCommand("go east")

	msg, success, shouldExit := game.Unlock()
	if !success {
		t.Error("Unlock with key should succeed")
	}
	if shouldExit {
		t.Error("Unlock of non-treasure room should not exit")
	}
	if msg != "You unlocked the door." {
		t.Errorf("Expected unlock message, got: %s", msg)
	}
}

// --- Helper functions for these tests ---

func createRoomWithExits(name string) *world.Room {
	return &world.Room{Name: name, Description: "This is " + name + ".", Exits: make(map[string]*world.Exit)}
}

func linkRooms(a *world.Room, dirAB string, b *world.Room, dirBA string) {
	a.Exits[dirAB] = &world.Exit{Room: b}
	b.Exits[dirBA] = &world.Exit{Room: a}
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
