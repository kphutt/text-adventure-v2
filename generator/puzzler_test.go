package generator

import (
	"testing"
	"text-adventure-v2/world"
)

// buildLinearRooms creates a chain: R0 -- R1 -- R2 -- ... -- Rn
func buildLinearRooms(n int) (*world.Room, map[string]*world.Room) {
	rooms := make([]*world.Room, n)
	allRooms := make(map[string]*world.Room)

	for i := 0; i < n; i++ {
		name := string(rune('A' + i))
		rooms[i] = &world.Room{Name: name, Exits: make(map[string]*world.Exit), X: i, Y: 0}
		allRooms[name] = rooms[i]
	}
	for i := 0; i < n-1; i++ {
		rooms[i].Exits["east"] = &world.Exit{Room: rooms[i+1]}
		rooms[i+1].Exits["west"] = &world.Exit{Room: rooms[i]}
	}
	return rooms[0], allRooms
}

// --- bfs tests ---

func TestBfs_Adjacent(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	a.Exits["east"] = &world.Exit{Room: b}
	b.Exits["west"] = &world.Exit{Room: a}

	path, err := bfs(a, b)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	if len(path) != 2 || path[0] != a || path[1] != b {
		t.Errorf("Expected path [A, B], got %v", path)
	}
}

func TestBfs_SameRoom(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}

	path, err := bfs(a, a)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	if len(path) != 1 {
		t.Errorf("Expected path of length 1, got %d", len(path))
	}
}

func TestBfs_NoPath(t *testing.T) {
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}

	_, err := bfs(a, b)
	if err == nil {
		t.Fatal("Expected error when no path exists")
	}
}

func TestBfs_MultiHop(t *testing.T) {
	start, allRooms := buildLinearRooms(5) // A--B--C--D--E
	end := allRooms["E"]

	path, err := bfs(start, end)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	if len(path) != 5 {
		t.Errorf("Expected path of length 5, got %d", len(path))
	}
	if path[0] != start || path[len(path)-1] != end {
		t.Error("Path should start at A and end at E")
	}
}

func TestBfs_TraversesLockedDoors(t *testing.T) {
	// bfs (unlike validatorBfs) always traverses all exits regardless of lock status
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit)}
	a.Exits["east"] = &world.Exit{Room: b, Locked: true}
	b.Exits["west"] = &world.Exit{Room: a}

	path, err := bfs(a, b)
	if err != nil {
		t.Fatalf("bfs should traverse locked doors, got error: %v", err)
	}
	if len(path) != 2 {
		t.Errorf("Expected path of length 2, got %d", len(path))
	}
}

// --- findLongestPath tests ---

func TestFindLongestPath_Linear(t *testing.T) {
	start, allRooms := buildLinearRooms(5) // A--B--C--D--E

	path, err := findLongestPath(start, allRooms)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	// Longest path from A should be to E (5 rooms)
	if len(path) != 5 {
		t.Errorf("Expected longest path of length 5, got %d", len(path))
	}
	if path[0] != start {
		t.Error("Longest path should start at the start room")
	}
}

func TestFindLongestPath_SingleRoom(t *testing.T) {
	start := &world.Room{Name: "A", Exits: make(map[string]*world.Exit)}
	allRooms := map[string]*world.Room{"A": start}

	_, err := findLongestPath(start, allRooms)
	if err == nil {
		t.Fatal("Expected error for single room (no paths to other rooms)")
	}
}

func TestFindLongestPath_BranchedMap(t *testing.T) {
	// A -- B -- C
	//      |
	//      D -- E -- F
	a := &world.Room{Name: "A", Exits: make(map[string]*world.Exit), X: 0, Y: 0}
	b := &world.Room{Name: "B", Exits: make(map[string]*world.Exit), X: 1, Y: 0}
	c := &world.Room{Name: "C", Exits: make(map[string]*world.Exit), X: 2, Y: 0}
	d := &world.Room{Name: "D", Exits: make(map[string]*world.Exit), X: 1, Y: 1}
	e := &world.Room{Name: "E", Exits: make(map[string]*world.Exit), X: 1, Y: 2}
	f := &world.Room{Name: "F", Exits: make(map[string]*world.Exit), X: 1, Y: 3}

	a.Exits["east"] = &world.Exit{Room: b}
	b.Exits["west"] = &world.Exit{Room: a}
	b.Exits["east"] = &world.Exit{Room: c}
	c.Exits["west"] = &world.Exit{Room: b}
	b.Exits["south"] = &world.Exit{Room: d}
	d.Exits["north"] = &world.Exit{Room: b}
	d.Exits["south"] = &world.Exit{Room: e}
	e.Exits["north"] = &world.Exit{Room: d}
	e.Exits["south"] = &world.Exit{Room: f}
	f.Exits["north"] = &world.Exit{Room: e}

	allRooms := map[string]*world.Room{"A": a, "B": b, "C": c, "D": d, "E": e, "F": f}

	path, err := findLongestPath(a, allRooms)
	if err != nil {
		t.Fatalf("Expected path, got error: %v", err)
	}
	// Longest BFS shortest-path from A is to F: A-B-D-E-F = 5
	if len(path) != 5 {
		t.Errorf("Expected longest path of length 5 (A to F), got %d", len(path))
	}
}
