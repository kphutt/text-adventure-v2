package generator

import (
	"errors"
	"text-adventure-v2/world"
)

// validateWorld ensures that the generated world is solvable.
func validateWorld(startRoom *world.Room, allRooms map[string]*world.Room) error {
	var keyRoom, treasureRoom *world.Room
	for _, room := range allRooms {
		for _, item := range room.Items {
			if item.Name == "key" {
				keyRoom = room
			}
			if item.Name == "treasure" {
				treasureRoom = room
			}
		}
	}

	if keyRoom == nil {
		return errors.New("validator: no key found in the world")
	}
	if treasureRoom == nil {
		return errors.New("validator: no treasure room found in the world")
	}

	// Test 1: Can you get from the start room to the key's room?
	_, err := validatorBfs(startRoom, keyRoom, false)
	if err != nil {
		return errors.New("validator: could not find a path from start to key")
	}

	// Test 2: Can you get from the start room to the treasure room (ignoring locks)?
	_, err = validatorBfs(startRoom, treasureRoom, true)
	if err != nil {
		return errors.New("validator: could not find a path from start to treasure (ignoring locks)")
	}

	// Test 3: Is there really no path to the treasure room if locks are considered?
	// This ensures the door is actually blocking the path.
	_, err = validatorBfs(startRoom, treasureRoom, false)
	if err == nil {
		return errors.New("validator: a path to treasure exists without needing the key")
	}

	return nil
}

// validatorBfs finds if a path exists between two rooms.
func validatorBfs(start, end *world.Room, ignoreLocks bool) ([]*world.Room, error) {
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
				if !exit.Locked || ignoreLocks {
					visited[exit.Room] = true
					newPath := make([]*world.Room, len(path))
					copy(newPath, path)
					newPath = append(newPath, exit.Room)
					queue = append(queue, newPath)
				}
			}
		}
	}

	return nil, errors.New("no path found between start and end")
}
