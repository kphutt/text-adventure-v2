package game

import (
	"fmt"
	"strings"
)

// NewGame creates a new game instance.
func NewGame() *Game {
	startRoom := CreateWorld()
	allRooms := make(map[string]*Room)
	GetAllRooms(startRoom, allRooms)

	player := &Player{
		Name:      "Player",
		Location:  startRoom,
		Inventory: make([]*Item, 0),
	}

	return &Game{
		Player:   player,
		AllRooms: allRooms,
		IsWon:    false,
	}
}

// GetAllRooms recursively finds all rooms starting from a given room.
func GetAllRooms(room *Room, rooms map[string]*Room) {
	rooms[room.Name] = room
	for _, exit := range room.Exits {
		if _, ok := rooms[exit.Room.Name]; !ok {
			GetAllRooms(exit.Room, rooms)
		}
	}
}

// HandleCommand processes a player command.
func (g *Game) HandleCommand(command string) (string, bool) {
	verb, noun := ParseInput(strings.ToLower(command))

	switch verb {
	case "quit", "q":
		return "Goodbye!", true
	case "help", "h":
		return "Instant Commands: w,a,s,d (move), e (take), i (inventory), u (unlock), l (look), q (quit)\nTyped Commands: go [dir], take [item], drop [item], unlock, help, quit", false
	case "look", "l":
		return g.Look(), false
	case "inventory", "i":
		return g.Inventory(), false
	case "go":
		return g.Move(noun)
	case "w", "a", "s", "d":
		var dir string
		dir = map[string]string{"w": "north", "a": "west", "s": "south", "d": "east"}[verb]
		return g.Move(dir)
	case "take":
		return g.Take(noun), false
	case "e":
		if len(g.Player.Location.Items) > 0 {
			return g.Take(g.Player.Location.Items[0].Name), false
		}
		return "There is nothing to take.", false
	case "drop":
		return g.Drop(noun), false
	case "unlock", "u":
		return g.Unlock(), false
	default:
		return "I don't understand that command.", false
	}
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
		return "", false // The main loop will call Look, so we return an empty message
	}
	return "You can't go that way.", false
}

// Take picks up an item from the current room.
func (g *Game) Take(itemName string) string {
	if itemName == "" {
		if len(g.Player.Location.Items) == 1 {
			itemName = g.Player.Location.Items[0].Name
		} else {
			return "What do you want to take?"
		}
	}

	for i, item := range g.Player.Location.Items {
		if item.Name == itemName {
			g.Player.Inventory = append(g.Player.Inventory, item)
			g.Player.Location.Items = append(g.Player.Location.Items[:i], g.Player.Location.Items[i+1:]...)
			return "You took the " + item.Name + "."
		}
	}
	return "You don't see that here."
}

// Drop drops an item into the current room.
func (g *Game) Drop(itemName string) string {
	if itemName == "" {
		return "What do you want to drop?"
	}

	for i, item := range g.Player.Inventory {
		if item.Name == itemName {
			g.Player.Location.Items = append(g.Player.Location.Items, item)
			g.Player.Inventory = append(g.Player.Inventory[:i], g.Player.Inventory[i+1:]...)
			return "You dropped the " + item.Name + "."
		}
	}
	return "You don't have that."
}

// Unlock unlocks a door.
func (g *Game) Unlock() string {
	if g.Player.Location.Name == "Dungeon" {
		if exit, ok := g.Player.Location.Exits["north"]; ok && exit.Locked {
			for _, item := range g.Player.Inventory {
				if item.Name == "key" {
					exit.Locked = false
					g.IsWon = true
					return "You unlocked the door!"
				}
			}
			return "You don't have the key."
		}
	}
	return "There is nothing to unlock here."
}

