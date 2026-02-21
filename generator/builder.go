package generator

import (
	"errors"
	"math/rand"
	"text-adventure-v2/world"
)

// buildWorld creates the raw structure of the world (rooms and their connections).
func buildWorld(config Config) (*world.Room, map[string]*world.Room, error) {

	if len(config.RoomNamePool) < config.NumberOfRooms {
		return nil, nil, errors.New("not enough unique room names in the pool for the number of rooms requested")
	}

	startRoom := &world.Room{
		Name:        "Starting Room",
		Description: "You find yourself in a plain room with a single, sturdy door.",
		Exits:       make(map[string]*world.Exit),
		Items:       make([]*world.Item, 0),
		X:           0,
		Y:           0,
	}

	allRooms := make(map[string]*world.Room)
	allRooms[startRoom.Name] = startRoom

	grid := make(map[int]map[int]*world.Room)
	grid[0] = make(map[int]*world.Room)
	grid[0][0] = startRoom

	roomNamePool := make([]string, len(config.RoomNamePool))
	copy(roomNamePool, config.RoomNamePool)
	rand.Shuffle(len(roomNamePool), func(i, j int) { roomNamePool[i], roomNamePool[j] = roomNamePool[j], roomNamePool[i] })

	for i := 1; i < config.NumberOfRooms; i++ {
		var created bool
		for !created {
			// Pick a random existing room to branch off from
			var randomRoom *world.Room
			for _, r := range allRooms {
				randomRoom = r
				break
			}

			// Pick a random direction
			dirs := []string{"north", "south", "east", "west"}
			dir := dirs[rand.Intn(len(dirs))]

			dx, dy := 0, 0
			var oppositeDir string
			switch dir {
			case "north":
				dy = -1
				oppositeDir = "south"
			case "south":
				dy = 1
				oppositeDir = "north"
			case "east":
				dx = 1
				oppositeDir = "west"
			case "west":
				dx = -1
				oppositeDir = "east"
			}

			newX, newY := randomRoom.X+dx, randomRoom.Y+dy

			// Check if the space is already occupied
			if _, exists := grid[newX][newY]; !exists {
				// Create new room
				newName := roomNamePool[i-1]
				newDesc := config.RoomDescPool[rand.Intn(len(config.RoomDescPool))]
				newRoom := &world.Room{
					Name:        newName,
					Description: newDesc,
					Exits:       make(map[string]*world.Exit),
					Items:       make([]*world.Item, 0),
					X:           newX,
					Y:           newY,
				}

				// Connect rooms
				randomRoom.Exits[dir] = &world.Exit{Room: newRoom}
				newRoom.Exits[oppositeDir] = &world.Exit{Room: randomRoom}

				// Add to collections
				allRooms[newRoom.Name] = newRoom
				if grid[newX] == nil {
					grid[newX] = make(map[int]*world.Room)
				}
				grid[newX][newY] = newRoom
				created = true
			}
		}
	}

	return startRoom, allRooms, nil
}
