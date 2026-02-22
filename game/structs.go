package game

import "text-adventure-v2/world"

// Game holds the entire state of the game.
type Game struct {
	Player       *world.Player
	AllRooms     map[string]*world.Room
	IsWon        bool
	Turns        int
	VisitedRooms map[string]bool
}
