package internal

import "fmt"

func FindAllPaths() {
	visited := make(map[string]bool)
	path := []string{}
	allPaths = [][]string{} // clear previous results just in case
	tunnelDebug()
	DFS(startRoom, visited, path)
	pathDebug()
}

func DFS(current string, visited map[string]bool, path []string) {
	// Mark current room as visited and add to path
	visited[current] = true
	path = append(path, current)

	// If we reached the endRoom, save a copy of the path
	if current == endRoom {
		// Make a deep copy of path to avoid mutation
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		allPaths = append(allPaths, pathCopy)
	}

	// Recurse into neighbors
	for _, neighbor := range tunnels[current] {
		if !visited[neighbor] {
			DFS(neighbor, visited, path)
		}
	}

	// Backtrack: unmark current room
	visited[current] = false
}

func pathDebug() {
	// print how many were found
	for i, path := range allPaths {
		Log(fmt.Sprintf("Path %d: %v", i, path), "debug")
	}
	Log(fmt.Sprintf("Found %d valid paths  [START]:%s -> [END]:%s", len(allPaths), startRoom, endRoom), "debug")
}

func tunnelDebug() {
	Log("Tunnels", "debug")
	for k, v := range tunnels {
		Log(fmt.Sprintf("Tunnel : %q: %q", k, v), "debug")
	}
	for name := range rooms {
		Log(fmt.Sprintf("Room: %q", name), "debug")
	}
}
