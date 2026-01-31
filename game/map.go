package game

import "strings"

// GetMapString generates a string representation of the world map.
func (g *Game) GetMapString() string {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, room := range g.AllRooms {
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
			for _, room := range g.AllRooms {
				if room.X == x && room.Y == y {
					if g.Player.Location == room {
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
