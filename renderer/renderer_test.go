package renderer

import (
	"testing"
	"text-adventure-v2/world"
)

// helper to wire bidirectional exits between two rooms
func link(a *world.Room, dirAB string, b *world.Room, dirBA string, locked bool) {
	a.Exits[dirAB] = &world.Exit{Room: b, Locked: locked}
	b.Exits[dirBA] = &world.Exit{Room: a, Locked: locked}
}

func TestRenderMap_SingleRoom(t *testing.T) {
	room := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": room},
		PlayerLocation: room,
		VisitedRooms:   map[string]bool{"A": true},
	}

	expected := "┌───┐\n│ @ │\n└───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("SingleRoom failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_HorizontalConnection(t *testing.T) {
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 1, Y: 0, Exits: map[string]*world.Exit{}}
	link(roomA, "east", roomB, "west", false)

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB},
		PlayerLocation: roomA,
		VisitedRooms:   map[string]bool{"A": true, "B": true},
	}

	expected := "┌───┐       ┌───┐\n" +
		"│ @ ├───────┤ . │\n" +
		"└───┘       └───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("HorizontalConnection failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_VerticalConnection(t *testing.T) {
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 0, Y: 1, Exits: map[string]*world.Exit{}}
	link(roomA, "south", roomB, "north", false)

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB},
		PlayerLocation: roomA,
		VisitedRooms:   map[string]bool{"A": true, "B": true},
	}

	expected := "┌───┐\n" +
		"│ @ │\n" +
		"└─┬─┘\n" +
		"  │\n" +
		"┌─┴─┐\n" +
		"│ . │\n" +
		"└───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("VerticalConnection failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_LockedHorizontalDoor(t *testing.T) {
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 1, Y: 0, Exits: map[string]*world.Exit{}}
	link(roomA, "east", roomB, "west", true)

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB},
		PlayerLocation: roomA,
		VisitedRooms:   map[string]bool{"A": true},
	}

	expected := "┌───┐       ┌───┐\n" +
		"│ @ ╠═══════╣   │\n" +
		"└───┘       └───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("LockedHorizontalDoor failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_LockedVerticalDoor(t *testing.T) {
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 0, Y: 1, Exits: map[string]*world.Exit{}}
	link(roomA, "south", roomB, "north", true)

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB},
		PlayerLocation: roomA,
		VisitedRooms:   map[string]bool{"A": true},
	}

	expected := "┌───┐\n" +
		"│ @ │\n" +
		"└─╦─┘\n" +
		"  ║\n" +
		"┌─╩─┐\n" +
		"│   │\n" +
		"└───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("LockedVerticalDoor failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_OneSidedLock(t *testing.T) {
	// Lock only set on one direction (matches how the generator works)
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 0, Y: 1, Exits: map[string]*world.Exit{}}
	// A's south exit is locked, but B's north exit is not
	roomA.Exits["south"] = &world.Exit{Room: roomB, Locked: true}
	roomB.Exits["north"] = &world.Exit{Room: roomA, Locked: false}

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB},
		PlayerLocation: roomA,
		VisitedRooms:   map[string]bool{"A": true},
	}

	expected := "┌───┐\n" +
		"│ @ │\n" +
		"└─╦─┘\n" +
		"  ║\n" +
		"┌─╩─┐\n" +
		"│   │\n" +
		"└───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("OneSidedLock failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_LShape(t *testing.T) {
	// A at (0,0), B at (0,1), C at (1,1) — an L shape
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 0, Y: 1, Exits: map[string]*world.Exit{}}
	roomC := &world.Room{Name: "C", X: 1, Y: 1, Exits: map[string]*world.Exit{}}
	link(roomA, "south", roomB, "north", false)
	link(roomB, "east", roomC, "west", false)

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB, "C": roomC},
		PlayerLocation: roomB,
		VisitedRooms:   map[string]bool{"A": true, "B": true, "C": true},
	}

	expected := "┌───┐\n" +
		"│ . │\n" +
		"└─┬─┘\n" +
		"  │\n" +
		"┌─┴─┐       ┌───┐\n" +
		"│ @ ├───────┤ . │\n" +
		"└───┘       └───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("LShape failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderMap_FogOfWar(t *testing.T) {
	// Three rooms: player in A, B visited, C unvisited
	roomA := &world.Room{Name: "A", X: 0, Y: 0, Exits: map[string]*world.Exit{}}
	roomB := &world.Room{Name: "B", X: 1, Y: 0, Exits: map[string]*world.Exit{}}
	roomC := &world.Room{Name: "C", X: 2, Y: 0, Exits: map[string]*world.Exit{}}
	link(roomA, "east", roomB, "west", false)
	link(roomB, "east", roomC, "west", false)

	view := MapView{
		AllRooms:       map[string]*world.Room{"A": roomA, "B": roomB, "C": roomC},
		PlayerLocation: roomA,
		VisitedRooms:   map[string]bool{"A": true, "B": true},
	}

	expected := "┌───┐       ┌───┐       ┌───┐\n" +
		"│ @ ├───────┤ . ├───────┤   │\n" +
		"└───┘       └───┘       └───┘"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("FogOfWar failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderHUD(t *testing.T) {
	view := MapView{
		CurrentLocationName: "Test Room",
		TurnsTaken:          42,
		Score:               75,
	}

	actual := RenderHUD(view)
	expected := "Location: Test Room\n" +
		"Turns: 42\n" +
		"Score: 75\n" +
		"--------------------------------------------------"

	if actual != expected {
		t.Errorf("RenderHUD failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestRenderHUD_ZeroScore(t *testing.T) {
	view := MapView{
		CurrentLocationName: "Start",
		TurnsTaken:          0,
		Score:               5,
	}

	actual := RenderHUD(view)
	expected := "Location: Start\n" +
		"Turns: 0\n" +
		"Score: 5\n" +
		"--------------------------------------------------"

	if actual != expected {
		t.Errorf("RenderHUD_ZeroScore failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}
