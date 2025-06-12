package internal

import (
	"os"
)

// We use BFS for quick lookup of at least one valid connection start -> end
func ValidateConnectivity() {
	visited := make(map[string]bool) // keep track of rooms visited
	queue := []string{startRoom}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == endRoom {
			return // success: path exists
		}

		visited[current] = true

		for _, neighbor := range tunnels[current] { // gives you all rooms connected to the current room
			if !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}

	// If we finished BFS without finding startRoom -> endRoom
	Log("no path from start to end", "error")
	os.Exit(0)
}
