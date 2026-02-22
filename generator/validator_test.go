package generator

import (
	"testing"
	"text-adventure-v2/world"
)

// buildValidWorld creates a simple 3-room world: Start -> Middle -[locked]-> Treasure
// with a key in Start. This is the minimal valid world.
func buildValidWorld() (*world.Room, map[string]*world.Room) {
	start := &world.Room{Name: "Start", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "key", Description: "A key."}}}
	middle := &world.Room{Name: "Middle", Exits: make(map[string]*world.Exit)}
	treasure := &world.Room{Name: "Treasure", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "treasure", Description: "Gold!"}}}

	start.Exits["east"] = &world.Exit{Room: middle}
	middle.Exits["west"] = &world.Exit{Room: start}
	middle.Exits["east"] = &world.Exit{Room: treasure, Locked: true}
	treasure.Exits["west"] = &world.Exit{Room: middle}

	allRooms := map[string]*world.Room{
		"Start": start, "Middle": middle, "Treasure": treasure,
	}
	return start, allRooms
}

func TestValidateWorld_Valid(t *testing.T) {
	start, allRooms := buildValidWorld()
	if err := validateWorld(start, allRooms); err != nil {
		t.Errorf("Expected valid world, got error: %v", err)
	}
}

func TestValidateWorld_NoKey(t *testing.T) {
	start, allRooms := buildValidWorld()
	// Remove the key
	start.Items = nil

	err := validateWorld(start, allRooms)
	if err == nil {
		t.Fatal("Expected error for missing key, got nil")
	}
}

func TestValidateWorld_NoTreasure(t *testing.T) {
	start, allRooms := buildValidWorld()
	// Remove the treasure
	allRooms["Treasure"].Items = nil

	err := validateWorld(start, allRooms)
	if err == nil {
		t.Fatal("Expected error for missing treasure, got nil")
	}
}

func TestValidateWorld_TreasureReachableWithoutKey(t *testing.T) {
	// No locked door — treasure is directly reachable
	start := &world.Room{Name: "Start", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "key", Description: "A key."}}}
	treasure := &world.Room{Name: "Treasure", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "treasure", Description: "Gold!"}}}

	start.Exits["east"] = &world.Exit{Room: treasure}
	treasure.Exits["west"] = &world.Exit{Room: start}

	allRooms := map[string]*world.Room{"Start": start, "Treasure": treasure}

	err := validateWorld(start, allRooms)
	if err == nil {
		t.Fatal("Expected error when treasure is reachable without needing the key")
	}
}

func TestValidateWorld_TreasureUnreachable(t *testing.T) {
	// Treasure is completely disconnected (no path even ignoring locks)
	start := &world.Room{Name: "Start", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "key", Description: "A key."}}}
	treasure := &world.Room{Name: "Treasure", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "treasure", Description: "Gold!"}}}

	// No exits connecting them at all
	allRooms := map[string]*world.Room{"Start": start, "Treasure": treasure}

	err := validateWorld(start, allRooms)
	if err == nil {
		t.Fatal("Expected error when treasure is unreachable")
	}
}

func TestValidateWorld_KeyUnreachable(t *testing.T) {
	// Key is behind a locked door
	start := &world.Room{Name: "Start", Exits: make(map[string]*world.Exit)}
	keyRoom := &world.Room{Name: "KeyRoom", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "key", Description: "A key."}}}
	treasure := &world.Room{Name: "Treasure", Exits: make(map[string]*world.Exit),
		Items: []*world.Item{{Name: "treasure", Description: "Gold!"}}}

	start.Exits["east"] = &world.Exit{Room: keyRoom, Locked: true}
	keyRoom.Exits["west"] = &world.Exit{Room: start}
	keyRoom.Exits["east"] = &world.Exit{Room: treasure, Locked: true}
	treasure.Exits["west"] = &world.Exit{Room: keyRoom}

	allRooms := map[string]*world.Room{"Start": start, "KeyRoom": keyRoom, "Treasure": treasure}

	err := validateWorld(start, allRooms)
	if err == nil {
		t.Fatal("Expected error when key is unreachable")
	}
}

// --- validatorBfs tests ---

func TestValidatorBfs_DirectPath(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	a.Exits["east"] = &world.Exit{Room: b}
	b.Exits["west"] = &world.Exit{Room: a}

	path, err := validatorBfs(a, b, false)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	if len(path) != 2 {
		t.Errorf("Expected path of length 2, got %d", len(path))
	}
}

func TestValidatorBfs_BlockedByLock(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	a.Exits["east"] = &world.Exit{Room: b, Locked: true}
	b.Exits["west"] = &world.Exit{Room: a}

	_, err := validatorBfs(a, b, false)
	if err == nil {
		t.Fatal("Expected error when path is blocked by locked door")
	}
}

func TestValidatorBfs_IgnoreLocks(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	a.Exits["east"] = &world.Exit{Room: b, Locked: true}
	b.Exits["west"] = &world.Exit{Room: a}

	path, err := validatorBfs(a, b, true)
	if err != nil {
		t.Fatalf("Expected path when ignoring locks, got error: %v", err)
	}
	if len(path) != 2 {
		t.Errorf("Expected path of length 2, got %d", len(path))
	}
}

func TestValidatorBfs_NoPath(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	// No exits connecting them

	_, err := validatorBfs(a, b, false)
	if err == nil {
		t.Fatal("Expected error when no path exists")
	}
}

func TestValidatorBfs_SameRoom(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}

	path, err := validatorBfs(a, a, false)
	if err != nil {
		t.Fatalf("Expected path from room to itself, got error: %v", err)
	}
	if len(path) != 1 {
		t.Errorf("Expected path of length 1, got %d", len(path))
	}
}

func TestValidatorBfs_MultiHopPath(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	c := &world.Room{Name: "C", Exits: make(map[string]*world.Exit)}
	d := &world.Room{Name: "D", Exits: make(map[string]*world.Exit)}

	a.Exits["east"] = &world.Exit{Room: b}
	b.Exits["west"] = &world.Exit{Room: a}
	b.Exits["east"] = &world.Exit{Room: c}
	c.Exits["west"] = &world.Exit{Room: b}
	c.Exits["east"] = &world.Exit{Room: d}
	d.Exits["west"] = &world.Exit{Room: c}

	path, err := validatorBfs(a, d, false)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	if len(path) != 4 {
		t.Errorf("Expected path of length 4, got %d", len(path))
	}
}
