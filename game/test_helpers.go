package game

import "text-adventure-v2/world"

// This file contains helper functions for creating specific test worlds.

// createSimpleLayout creates a simple, predictable 3-room world for general testing.
// Layout: Room A <--> Room B (starts here) <--> Room C
func createSimpleLayout() *Game {
	roomA := &world.Room{Name: "Room A", Description: "This is Room A.", Exits: make(map[string]*world.Exit)}
	roomB := &world.Room{Name: "Room B", Description: "This is Room B.", Exits: make(map[string]*world.Exit)}
	roomC := &world.Room{Name: "Room C", Description: "This is Room C.", Exits: make(map[string]*world.Exit)}

	roomA.Exits["east"] = &world.Exit{Room: roomB}
	roomB.Exits["west"] = &world.Exit{Room: roomA}
	roomB.Exits["east"] = &world.Exit{Room: roomC}
	roomC.Exits["west"] = &world.Exit{Room: roomB}

	player := &world.Player{
		Name:     "Test Player",
		Location: roomB,
	}

	allRooms := map[string]*world.Room{
		roomA.Name: roomA,
		roomB.Name: roomB,
		roomC.Name: roomC,
	}

	return &Game{
		Player:   player,
		AllRooms: allRooms,
	}
}

// createLayoutWithItems creates a world specifically for testing take and drop.
// Layout: Room A (has "test_item") <--> Room B (starts here)
func createLayoutWithItems() *Game {
	roomA := &world.Room{
		Name:        "Room A",
		Description: "This is Room A.",
		Exits:       make(map[string]*world.Exit),
		Items:       []*world.Item{{Name: "test_item", Description: "A test item."}},
	}
	roomB := &world.Room{Name: "Room B", Description: "This is Room B.", Exits: make(map[string]*world.Exit)}

	roomA.Exits["east"] = &world.Exit{Room: roomB}
	roomB.Exits["west"] = &world.Exit{Room: roomA}

	player := &world.Player{
		Name:     "Test Player",
		Location: roomB,
	}

	allRooms := map[string]*world.Room{
		roomA.Name: roomA,
		roomB.Name: roomB,
	}

	return &Game{
		Player:   player,
		AllRooms: allRooms,
	}
}

// createLayoutWithLock creates a world specifically for testing the unlock command.
// Layout: Room A (has "key") <--> Room B (starts here) <--> [LOCKED] <--> Room C
func createLayoutWithLock() *Game {
	roomA := &world.Room{
		Name:        "Room A",
		Description: "This is Room A.",
		Exits:       make(map[string]*world.Exit),
		Items:       []*world.Item{{Name: "key", Description: "A test key."}},
	}
	roomB := &world.Room{Name: "Room B", Description: "This is Room B.", Exits: make(map[string]*world.Exit)}
	roomC := &world.Room{Name: "Room C", Description: "This is Room C.", Exits: make(map[string]*world.Exit)}

	roomA.Exits["east"] = &world.Exit{Room: roomB}
	roomB.Exits["west"] = &world.Exit{Room: roomA}
	roomB.Exits["east"] = &world.Exit{Room: roomC, Locked: true}
	roomC.Exits["west"] = &world.Exit{Room: roomB}

	player := &world.Player{
		Name:     "Test Player",
		Location: roomB,
	}

	allRooms := map[string]*world.Room{
		roomA.Name: roomA,
		roomB.Name: roomB,
		roomC.Name: roomC,
	}

	return &Game{
		Player:   player,
		AllRooms: allRooms,
	}
}

// createLayoutWithWinCondition creates a world for testing the win condition.
// Layout: Room A (has "key") <--> Room B (starts here) <--> [LOCKED] <--> Treasure Room
func createLayoutWithWinCondition() *Game {
	roomA := &world.Room{
		Name:        "Room A",
		Description: "This is Room A.",
		Exits:       make(map[string]*world.Exit),
		Items:       []*world.Item{{Name: "key", Description: "A test key."}},
	}
	roomB := &world.Room{Name: "Room B", Description: "This is Room B.", Exits: make(map[string]*world.Exit)}
	treasureRoom := &world.Room{Name: "Treasure Room", Description: "The treasure is here!", Exits: make(map[string]*world.Exit)}

	roomA.Exits["east"] = &world.Exit{Room: roomB}
	roomB.Exits["west"] = &world.Exit{Room: roomA}
	roomB.Exits["east"] = &world.Exit{Room: treasureRoom, Locked: true}
	treasureRoom.Exits["west"] = &world.Exit{Room: roomB}

	player := &world.Player{
		Name:     "Test Player",
		Location: roomB,
	}

	allRooms := map[string]*world.Room{
		roomA.Name:        roomA,
		roomB.Name:        roomB,
		treasureRoom.Name: treasureRoom,
	}

	return &Game{
		Player:   player,
		AllRooms: allRooms,
	}
}
