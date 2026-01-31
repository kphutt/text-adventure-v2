package renderer

import (
	"testing"
	"text-adventure-v2/world"
)

func TestRenderMap_SimpleRow(t *testing.T) {
	// Create a simple layout: three rooms in a row
	roomA := &world.Room{Name: "A", X: 0, Y: 0}
	roomB := &world.Room{Name: "B", X: 1, Y: 0}
	roomC := &world.Room{Name: "C", X: 2, Y: 0}

	allRooms := map[string]*world.Room{
		"A": roomA,
		"B": roomB,
		"C": roomC,
	}

	// Case 1: Player is in the middle room
	view := MapView{
		AllRooms:       allRooms,
		PlayerLocation: roomB,
	}

	expected := "[ ][@][ ]\n"
	actual := RenderMap(view)

	if actual != expected {
		t.Errorf("RenderMap SimpleRow failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}

	// Case 2: Player is in the first room
	view.PlayerLocation = roomA
	expected = "[@][ ][ ]\n"
	actual = RenderMap(view)
		if actual != expected {
			t.Errorf("RenderMap SimpleRow failed for player in first room.\nExpected:\n%s\nGot:\n%s", expected, actual)
		}
	}
	
	func TestRenderMap_MultiLine(t *testing.T) {
		// Create an L-shaped layout
		roomA := &world.Room{Name: "A", X: 0, Y: 0}
		roomB := &world.Room{Name: "B", X: 0, Y: 1}
		roomC := &world.Room{Name: "C", X: 1, Y: 1}
	
		allRooms := map[string]*world.Room{
			"A": roomA,
			"B": roomB,
			"C": roomC,
		}
	
		view := MapView{
			AllRooms:       allRooms,
			PlayerLocation: roomC,
		}
	
		// Note: The extra spaces at the end of the lines are important.
		expected := "[ ]   \n" +
			        "[ ][@]\n"
	
		actual := RenderMap(view)
	
		if actual != expected {
			t.Errorf("RenderMap MultiLine failed.\nExpected:\n%s\nGot:\n%s", expected, actual)
		}
	}
	