package game

// Item represents an object that can be picked up and dropped.
type Item struct {
	Name        string
	Description string
}

// Exit represents a connection from one room to another.
type Exit struct {
	Room   *Room
	Locked bool
}

// Room represents a location in the game world.
type Room struct {
	Name        string
	Description string
	Exits       map[string]*Exit
	Items       []*Item
	X, Y        int
}

// Player represents the user in the game.
type Player struct {
	Name      string
	Location  *Room
	Inventory []*Item
}

// Game holds the entire state of the game.
type Game struct {
	Player   *Player
	AllRooms map[string]*Room
	IsWon    bool
}