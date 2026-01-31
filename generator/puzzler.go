package generator

import (
	"errors"
	"math/rand"
	"text-adventure-v2/world"
)

// placePuzzles finds a path for the main puzzle and places the key and locked door.
func placePuzzles(config Config, startRoom *world.Room, allRooms map[string]*world.Room) error {
	// Find the longest path to a dead end to be the treasure room.
	path, err := findLongestPath(startRoom, allRooms)
	if err != nil {
		return err
	}

	if len(path) < config.MinPathToTreasure {
		return errors.New("could not find a path long enough to satisfy MinPathToTreasure")
	}

	// The end of the path is the treasure room.
	treasureRoom := path[len(path)-1]
	treasureRoom.Name = "Treasure Room"
	treasureRoom.Description = "You have found the treasure room! A large chest sits in the center."
	treasureRoom.Items = append(treasureRoom.Items, &world.Item{Name: "treasure", Description: "A chest full of gold!"})

	// Place the locked door in the middle of the path.
	doorIndex := len(path) / 2
	doorRoom := path[doorIndex]
	nextRoomInPath := path[doorIndex+1]
	for _, exit := range doorRoom.Exits {
		if exit.Room == nextRoomInPath {
			exit.Locked = true
			break
		}
	}

	// Place the key somewhere on the path before the locked door.
	keyIndex := rand.Intn(doorIndex)
	keyRoom := path[keyIndex]
	keyRoom.Items = append(keyRoom.Items, &world.Item{Name: "key", Description: "A small, rusty key."})

	// Place extra items
	for _, itemName := range config.ExtraItems {
		for {
			var randomRoom *world.Room
			// Find a random room
			for _, r := range allRooms {
				randomRoom = r
				if rand.Float32() < 0.5 { // Add some randomness to room selection
					break
				}
			}

			// Don't place items in the start room or rooms that already have items.
			if randomRoom != startRoom && len(randomRoom.Items) == 0 {
				randomRoom.Items = append(randomRoom.Items, &world.Item{Name: itemName, Description: "An extra item."})
				break
			}
		}
	}

	return nil
}

// findLongestPath uses BFS to find the longest path from the start room to any other room.
func findLongestPath(start *world.Room, allRooms map[string]*world.Room) ([]*world.Room, error) {
	var longestPath []*world.Room
	
	for _, room := range allRooms {
        if room == start {
            continue
        }

		path, err := bfs(start, room)
		if err != nil {
			continue // Should not happen in a connected graph
		}
		if len(path) > len(longestPath) {
			longestPath = path
		}
	}

	if len(longestPath) == 0 {
		return nil, errors.New("could not find any path from the start room")
	}
	
	return longestPath, nil
}


// bfs finds the shortest path between two rooms.
func bfs(start, end *world.Room) ([]*world.Room, error) {
	queue := [][]*world.Room{{start}}
	visited := map[*world.Room]bool{start: true}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node == end {
			return path, nil
		}

		for _, exit := range node.Exits {
			if !visited[exit.Room] {
				visited[exit.Room] = true
				newPath := make([]*world.Room, len(path))
				copy(newPath, path)
				newPath = append(newPath, exit.Room)
				queue = append(queue, newPath)
			}
		}
	}

	return nil, errors.New("no path found between start and end")
}
