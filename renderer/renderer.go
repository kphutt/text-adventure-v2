package renderer

import (
	"strings"
	"text-adventure-v2/world"
)

// MapView contains all the necessary data for the renderer to draw the map.
// It acts as a decoupled view model for the rendering engine.
type MapView struct {
	AllRooms       map[string]*world.Room
	PlayerLocation *world.Room
}

// RenderMap takes a MapView and produces an ASCII string representation of the map.
func RenderMap(view MapView) string {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, room := range view.AllRooms {
		if room.X < minX {
			minX = room.X
		}
		if room.Y < minY {
			minY = room.Y
		}
		if room.X > maxX {
			maxX = room.X
		}
		if room.Y > maxY {
			maxY = room.Y
		}
	}

	var b strings.Builder
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			isRoom := false
			for _, room := range view.AllRooms {
				if room.X == x && room.Y == y {
					if view.PlayerLocation == room {
						b.WriteString("[@]")
					} else {
						b.WriteString("[ ]")
					}
					isRoom = true
					break
				}
			}
			if !isRoom {
				b.WriteString("   ")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}