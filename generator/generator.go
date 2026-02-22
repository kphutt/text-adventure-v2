package generator

import (
	"fmt"
	"text-adventure-v2/world"
)

// Config holds the parameters for map generation.
type Config struct {
	NumberOfRooms     int
	MinPathToTreasure int
	ExtraItems        []string
	RoomNamePool      []string
	RoomDescPool      []string
}

// DefaultConfig provides sensible starting values for map generation.
func DefaultConfig() Config {
	return Config{
		NumberOfRooms:     10,
		MinPathToTreasure: 4,
		ExtraItems:        []string{"sword"},
		RoomNamePool: []string{
			"Dank Cellar",
			"Dusty Armory",
			"Forgotten Library",
			"Echoing Cavern",
			"Drafty Corridor",
			"Sunken Grotto",
			"Crystal Chamber",
			"Shadowy Antechamber",
			"Musty Crawlspace",
			"Alchemist's Laboratory",
		},
		RoomDescPool: []string{
			"You are in a small, damp room. A faint dripping sound echoes from a dark corner.",
			"The air is thick with the smell of old books and decaying paper. Shelves line the walls.",
			"A single torch flickers, casting long, dancing shadows across the cold stone floor.",
			"The ground is uneven and slick with moisture. Strange fungi glow with a soft, eerie light.",
			"You can feel a cold breeze, though you can't identify its source.",
			"This room is surprisingly ornate, with faded tapestries hanging on the walls.",
			"An old suit of armor stands in the corner, its helmet staring at you blankly.",
			"The ceiling is unusually high here, lost in the oppressive darkness above.",
		},
	}
}

// Generate orchestrates the creation of a new, random, and solvable game world.
func Generate(config Config) (*world.Room, error) {
	var err error
	const maxRetries = 10

	for i := 0; i < maxRetries; i++ {
		// Step 1: Build the raw world structure.
		var startRoom *world.Room
		var allRooms map[string]*world.Room
		startRoom, allRooms, err = buildWorld(config)
		if err != nil {
			continue // Should be rare, but retry if it happens
		}

		// Step 2: Place the puzzles and extra items.
		err = placePuzzles(config, startRoom, allRooms)
		if err != nil {
			continue // This can fail if the map is too simple, so we retry
		}

		// Step 3: Validate that the world is solvable.
		err = validateWorld(startRoom, allRooms)
		if err != nil {
			continue // Should be rare, but retry if validation fails
		}

		// If we get here, the world is valid.
		return startRoom, nil
	}

	// If we've exhausted all retries, return the last error encountered.
	return nil, fmt.Errorf("failed to generate a valid world after %d attempts: %w", maxRetries, err)
}
