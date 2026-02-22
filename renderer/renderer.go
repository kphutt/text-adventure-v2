package renderer

import (
	"fmt"
	"strings"
	"text-adventure-v2/world"
)

// MapView contains all the necessary data for the renderer to draw the map.
// It acts as a decoupled view model for the rendering engine.
type MapView struct {
	AllRooms            map[string]*world.Room
	PlayerLocation      *world.Room
	CurrentLocationName string
	TurnsTaken          int
	Score               int
	VisitedRooms        map[string]bool
}

// RenderMap takes a MapView and produces a box-drawing string representation of the map.
// Each room is a 5×3 box. Horizontal corridors are 7 chars wide, vertical corridors 1 char tall.
func RenderMap(view MapView) string {
	if len(view.AllRooms) == 0 {
		return ""
	}

	minX, minY, maxX, maxY := computeBounds(view.AllRooms)

	gridW := maxX - minX + 1
	gridH := maxY - minY + 1

	bufW := gridW*12 - 7
	bufH := gridH*4 - 1

	// Allocate buffer filled with spaces
	buf := make([][]rune, bufH)
	for i := range buf {
		buf[i] = make([]rune, bufW)
		for j := range buf[i] {
			buf[i][j] = ' '
		}
	}

	// Draw all room boxes
	for _, room := range view.AllRooms {
		ox, oy := roomOrigin(room.X-minX, room.Y-minY)
		interior := ' '
		if room == view.PlayerLocation {
			interior = '@'
		} else if view.VisitedRooms[room.Name] {
			interior = '.'
		}
		drawRoom(buf, ox, oy, interior)
	}

	// Draw corridors (only east and south to avoid double-drawing).
	// Check both sides for locked status since the generator may only lock one direction.
	for _, room := range view.AllRooms {
		for dir, exit := range room.Exits {
			switch dir {
			case "east":
				if exit.Room.X == room.X+1 && exit.Room.Y == room.Y {
					ox, oy := roomOrigin(room.X-minX, room.Y-minY)
					locked := exit.Locked || isExitLocked(exit.Room, "west")
					drawHCorridor(buf, ox, oy, locked)
				}
			case "south":
				if exit.Room.X == room.X && exit.Room.Y == room.Y+1 {
					ox, oy := roomOrigin(room.X-minX, room.Y-minY)
					locked := exit.Locked || isExitLocked(exit.Room, "north")
					drawVCorridor(buf, ox, oy, locked)
				}
			}
		}
	}

	return bufferToString(buf)
}

func computeBounds(rooms map[string]*world.Room) (minX, minY, maxX, maxY int) {
	first := true
	for _, room := range rooms {
		if first {
			minX, minY, maxX, maxY = room.X, room.Y, room.X, room.Y
			first = false
			continue
		}
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
	return
}

func isExitLocked(room *world.Room, dir string) bool {
	if e, ok := room.Exits[dir]; ok {
		return e.Locked
	}
	return false
}

func roomOrigin(gx, gy int) (int, int) {
	return gx * 12, gy * 4
}

func drawRoom(buf [][]rune, ox, oy int, interior rune) {
	// Top: ┌───┐
	buf[oy][ox] = '┌'
	buf[oy][ox+1] = '─'
	buf[oy][ox+2] = '─'
	buf[oy][ox+3] = '─'
	buf[oy][ox+4] = '┐'
	// Middle: │ X │
	buf[oy+1][ox] = '│'
	buf[oy+1][ox+1] = ' '
	buf[oy+1][ox+2] = interior
	buf[oy+1][ox+3] = ' '
	buf[oy+1][ox+4] = '│'
	// Bottom: └───┘
	buf[oy+2][ox] = '└'
	buf[oy+2][ox+1] = '─'
	buf[oy+2][ox+2] = '─'
	buf[oy+2][ox+3] = '─'
	buf[oy+2][ox+4] = '┘'
}

func drawHCorridor(buf [][]rune, ox, oy int, locked bool) {
	row := oy + 1
	rightWall := ox + 4
	leftWall := ox + 12

	if locked {
		buf[row][rightWall] = '╠'
		buf[row][leftWall] = '╣'
		for c := rightWall + 1; c < leftWall; c++ {
			buf[row][c] = '═'
		}
	} else {
		buf[row][rightWall] = '├'
		buf[row][leftWall] = '┤'
		for c := rightWall + 1; c < leftWall; c++ {
			buf[row][c] = '─'
		}
	}
}

func drawVCorridor(buf [][]rune, ox, oy int, locked bool) {
	col := ox + 2
	bottomWall := oy + 2
	topWall := oy + 4

	if locked {
		buf[bottomWall][col] = '╦'
		buf[topWall][col] = '╩'
		for r := bottomWall + 1; r < topWall; r++ {
			buf[r][col] = '║'
		}
	} else {
		buf[bottomWall][col] = '┬'
		buf[topWall][col] = '┴'
		for r := bottomWall + 1; r < topWall; r++ {
			buf[r][col] = '│'
		}
	}
}

func bufferToString(buf [][]rune) string {
	var b strings.Builder
	for i, row := range buf {
		line := strings.TrimRight(string(row), " ")
		b.WriteString(line)
		if i < len(buf)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// RenderHUD takes a MapView and produces a formatted string for the status bar.
func RenderHUD(view MapView) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Location: %s\n", view.CurrentLocationName))
	b.WriteString(fmt.Sprintf("Turns: %d\n", view.TurnsTaken))
	b.WriteString(fmt.Sprintf("Score: %d\n", view.Score))
	b.WriteString(strings.Repeat("-", 50)) // A separator line
	return b.String()
}
