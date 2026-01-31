package game

// CreateWorld initializes and connects all the rooms in the game.
func CreateWorld() *Room {
	// Create rooms
	hall := &Room{
		Name:        "Hall",
		Description: "You are in a long, dark hall. A single torch flickers on the wall.",
		Exits:       make(map[string]*Exit),
		Items:       make([]*Item, 0),
		X:           0,
		Y:           1,
	}
	closet := &Room{
		Name:        "Closet",
		Description: "You are in a small, dusty closet. It smells of old wood.",
		Exits:       make(map[string]*Exit),
		Items:       make([]*Item, 0),
		X:           1,
		Y:           1,
	}
	dungeon := &Room{
		Name:        "Dungeon",
		Description: "You are in a cold, damp dungeon. You hear a faint dripping sound.",
		Exits:       make(map[string]*Exit),
		Items:       make([]*Item, 0),
		X:           0,
		Y:           0,
	}
	treasureRoom := &Room{
		Name:        "Treasure Room",
		Description: "You have found the treasure room! A large chest sits in the center.",
		Exits:       make(map[string]*Exit),
		Items:       make([]*Item, 0),
		X:           0,
		Y:           -1,
	}

	// Create items
	sword := &Item{Name: "sword", Description: "A sharp, pointy sword."}
	key := &Item{Name: "key", Description: "A small, rusty key."}
	treasure := &Item{Name: "treasure", Description: "A chest full of gold!"}


	// Add items to rooms
	hall.Items = append(hall.Items, sword)
	closet.Items = append(closet.Items, key)
	treasureRoom.Items = append(treasureRoom.Items, treasure)

	// Connect rooms
	hall.Exits["north"] = &Exit{Room: dungeon}
	hall.Exits["east"] = &Exit{Room: closet}
	closet.Exits["west"] = &Exit{Room: hall}
	dungeon.Exits["south"] = &Exit{Room: hall}
	dungeon.Exits["north"] = &Exit{Room: treasureRoom, Locked: true}
	treasureRoom.Exits["south"] = &Exit{Room: dungeon}

	return hall
}
