package internal

import "fmt"

func FindAllPaths() {
	visited := make(map[string]bool)
	path := []string{}
	allPaths = [][]string{} // clear previous results just in case

	DFS(startRoom, visited, path)

	// print how many were found
	Log(fmt.Sprintf("Found %d valid paths from %s to %s", len(allPaths), startRoom, endRoom), "debug")
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
