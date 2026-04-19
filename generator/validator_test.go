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

// --- validateGeometry tests ---

func TestValidateGeometry_ValidDirections(t *testing.T) {
	cases := []struct {
		name   string
		dir    string
		dx, dy int
		back   string
	}{
		{"east", "east", 1, 0, "west"},
		{"west", "west", -1, 0, "east"},
		{"south", "south", 0, 1, "north"},
		{"north", "north", 0, -1, "south"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
			b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: c.dx, Y: c.dy}
			a.Exits[c.dir] = &world.Exit{Room: b}
			b.Exits[c.back] = &world.Exit{Room: a}

			allRooms := map[string]*world.Room{"A": a, "B": b}

			if err := validateGeometry(allRooms); err != nil {
				t.Errorf("Expected valid %s exit, got error: %v", c.dir, err)
			}
		})
	}
}

func TestValidateGeometry_WrongCoordsPerDirection(t *testing.T) {
	cases := []struct {
		name   string
		dir    string
		bx, by int // incorrect coords for B
	}{
		{"east points two steps away", "east", 2, 0},
		{"west points two steps away", "west", -2, 0},
		{"south points two steps away", "south", 0, 2},
		{"north points two steps away", "north", 0, -2},
		{"east has drifted on Y axis", "east", 1, 1},
		{"south has drifted on X axis", "south", 1, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
			b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: c.bx, Y: c.by}
			a.Exits[c.dir] = &world.Exit{Room: b}

			allRooms := map[string]*world.Room{"A": a, "B": b}

			if err := validateGeometry(allRooms); err == nil {
				t.Fatalf("Expected error for %s exit pointing to (%d,%d), got nil", c.dir, c.bx, c.by)
			}
		})
	}
}

func TestValidateGeometry_UnknownDirection(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	a.Exits["up"] = &world.Exit{Room: b}

	allRooms := map[string]*world.Room{"A": a, "B": b}

	if err := validateGeometry(allRooms); err == nil {
		t.Fatal("Expected error for unknown direction, got nil")
	}
}

func TestValidateGeometry_MultiRoomChain(t *testing.T) {
	// 4 rooms in a line: A(0,0) <-> B(1,0) <-> C(2,0) <-> D(3,0)
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: 1, Y: 0}
	c := &world.Room{Name: "C", Exits: make(map[string]*world.Exit), X: 2, Y: 0}
	d := &world.Room{Name: "D", Exits: make(map[string]*world.Exit), X: 3, Y: 0}

	a.Exits["east"] = &world.Exit{Room: b}
	b.Exits["west"] = &world.Exit{Room: a}
	b.Exits["east"] = &world.Exit{Room: c}
	c.Exits["west"] = &world.Exit{Room: b}
	c.Exits["east"] = &world.Exit{Room: d}
	d.Exits["west"] = &world.Exit{Room: c}

	allRooms := map[string]*world.Room{"A": a, "B": b, "C": c, "D": d}

	if err := validateGeometry(allRooms); err != nil {
		t.Errorf("Expected valid multi-room chain, got error: %v", err)
	}
}

func TestValidateGeometry_SameCoordinates(t *testing.T) {
	// Two rooms at (0,0) with an east exit between them — impossible geometry.
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	a.Exits["east"] = &world.Exit{Room: b}

	allRooms := map[string]*world.Room{"A": a, "B": b}

	if err := validateGeometry(allRooms); err == nil {
		t.Fatal("Expected error for two rooms at same coordinates with a directional exit, got nil")
	}
}

func TestValidateGeometry_PartiallyInvalid(t *testing.T) {
	// 3 rooms, first two correctly wired, third at bad coords.
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: 1, Y: 0}
	bad := &world.Room{Name: "Bad", Exits: make(map[string]*world.Exit), X: 5, Y: 5}

	a.Exits["east"] = &world.Exit{Room: b}
	b.Exits["west"] = &world.Exit{Room: a}
	b.Exits["east"] = &world.Exit{Room: bad} // expected (2,0), got (5,5)

	allRooms := map[string]*world.Room{"A": a, "B": b, "Bad": bad}

	if err := validateGeometry(allRooms); err == nil {
		t.Fatal("Expected error when one of three rooms has bad coordinates, got nil")
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
