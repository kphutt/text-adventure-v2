package game

import (
	"fmt"
	"strings"
	"text-adventure-v2/generator"
	"text-adventure-v2/world"
)

// NewGame creates a new game instance.
func NewGame() *Game {
	// Generate a new world
	startRoom, err := generator.Generate(generator.DefaultConfig())
	if err != nil {
		// For now, we'll panic. In a real application, you might want to handle this more gracefully.
		panic(fmt.Sprintf("failed to generate world: %v", err))
	}

	allRooms := make(map[string]*world.Room)
	GetAllRooms(startRoom, allRooms)

	player := &world.Player{
		Name:      "Player",
		Location:  startRoom,
		Inventory: make([]*world.Item, 0),
	}

	return &Game{
		Player:       player,
		AllRooms:     allRooms,
		IsWon:        false,
		Turns:        0,
		VisitedRooms: map[string]bool{startRoom.Name: true},
	}
}

// GetAllRooms recursively finds all rooms starting from a given room.
func GetAllRooms(room *world.Room, rooms map[string]*world.Room) {
	if _, ok := rooms[room.Name]; ok {
		return
	}
	rooms[room.Name] = room
	for _, exit := range room.Exits {
		GetAllRooms(exit.Room, rooms)
	}
}

// HandleCommand processes a player command and updates the game state.
func (g *Game) HandleCommand(command string) (string, bool) {
	verb, noun := ParseInput(strings.ToLower(command))

	var msg string
	var success, shouldExit bool

	switch verb {
	case "quit", "q":
		return "Goodbye!", true
	case "help", "h":
		return "Instant Commands: w,a,s,d (move), e (take), i (inventory), u (unlock), l (look), q (quit)\nTyped Commands: go [dir], take [item], drop [item], unlock, score, help, quit", false
	case "look", "l":
		return g.Look(), false
	case "inventory", "i":
		return g.Inventory(), false
	case "score":
		return fmt.Sprintf("Score: %d", g.Score()), false
	case "go":
		msg, success = g.Move(noun)
	case "w", "a", "s", "d":
		var dir string
		dir = map[string]string{"w": "north", "a": "west", "s": "south", "d": "east"}[verb]
		msg, success = g.Move(dir)
	case "take":
		msg, success = g.Take(noun)
	case "e":
		if len(g.Player.Location.Items) > 0 {
			msg, success = g.Take(g.Player.Location.Items[0].Name)
		} else {
			msg, success = "There is nothing to take.", false
		}
	case "drop":
		msg, success = g.Drop(noun)
	case "unlock", "u":
		msg, success, shouldExit = g.Unlock()
	default:
		msg, success = "I don't understand that command.", false
	}

	if success {
		g.Turns++
	}

	return msg, shouldExit
}

// Score returns the player's current score.
// 10 points per inventory item, 5 points per room visited.
func (g *Game) Score() int {
	return len(g.Player.Inventory)*10 + len(g.VisitedRooms)*5
}

// Look returns the description of the player's current location.
func (g *Game) Look() string {
	var b strings.Builder
	b.WriteString(g.Player.Location.Description + "\n")
	if len(g.Player.Location.Items) > 0 {
		b.WriteString("You see the following items:\n")
		for _, item := range g.Player.Location.Items {
			fmt.Fprintf(&b, "- %s\n", item.Name)
		}
	}
	b.WriteString("Exits:\n")
	for dir := range g.Player.Location.Exits {
		fmt.Fprintf(&b, "- %s\n", dir)
	}
	return b.String()
}

// Inventory returns a string listing the player's inventory.
func (g *Game) Inventory() string {
	if len(g.Player.Inventory) == 0 {
		return "You are not carrying anything."
	}
	var b strings.Builder
	b.WriteString("You have the following items:\n")
	for _, item := range g.Player.Inventory {
		fmt.Fprintf(&b, "- %s\n", item.Name)
	}
	return b.String()
}

// Move moves the player in the given direction.
func (g *Game) Move(direction string) (string, bool) {
	if exit, ok := g.Player.Location.Exits[direction]; ok {
		if exit.Locked {
			return "The door is locked.", false
		}
		g.Player.Location = exit.Room
		g.VisitedRooms[exit.Room.Name] = true
		return "", true
	}
	return "You can't go that way.", false
}

// Take picks up an item from the current room.
func (g *Game) Take(itemName string) (string, bool) {
	if itemName == "" {
		if len(g.Player.Location.Items) == 1 {
			itemName = g.Player.Location.Items[0].Name
		} else {
			return "What do you want to take?", false
		}
	}

	for i, item := range g.Player.Location.Items {
		if strings.ToLower(item.Name) == strings.ToLower(itemName) {
			g.Player.Inventory = append(g.Player.Inventory, item)
			g.Player.Location.Items = append(g.Player.Location.Items[:i], g.Player.Location.Items[i+1:]...)
			return "You took the " + item.Name + ".", true
		}
	}
	return "You don't see that here.", false
}

// Drop drops an item into the current room.
func (g *Game) Drop(itemName string) (string, bool) {
	if itemName == "" {
		return "What do you want to drop?", false
	}

	for i, item := range g.Player.Inventory {
		if strings.ToLower(item.Name) == strings.ToLower(itemName) {
			g.Player.Location.Items = append(g.Player.Location.Items, item)
			g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
			return "You dropped the " + item.Name + ".", true
		}
	}
	return "You don't have that.", false
}

// Unlock unlocks a door.
func (g *Game) Unlock() (string, bool, bool) {
	var lockedExit *world.Exit
	for _, exit := range g.Player.Location.Exits {
		if exit.Locked {
			lockedExit = exit
			break
		}
	}

	if lockedExit == nil {
		return "There is nothing to unlock here.", false, false
	}

	hasKey := false
	for _, item := range g.Player.Inventory {
		if item.Name == "key" {
			hasKey = true
			break
		}
	}

	if !hasKey {
		return "You don't have the key.", false, false
	}

	lockedExit.Locked = false
	if lockedExit.Room.Name == "Treasure Room" {
		g.IsWon = true
		return "You unlocked the door! You win!", false, true
	}
	return "You unlocked the door.", true, false
}
