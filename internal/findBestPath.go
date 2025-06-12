package internal

import (
	"fmt"
	"sort"
)

func FindBestPaths() {
	maxPaths := len(allPaths)

	Log("Evaluating additional disjoint paths (based on steps)...", "debug")
	sortedPaths := getSortedPathsBySteps()

	// Save the best step path (shortest path)
	bestStepPath = sortedPaths[0]

	// Remove the first path from the sorted list before selecting disjoint paths
	remainingPaths := sortedPaths[1:]

	// Select disjoint paths excluding the best path
	bestStepDisjointPaths = append([][]string{bestStepPath}, selectDisjointPaths(remainingPaths, maxPaths-1)...)

	// Ensure that, besides the best path, additional paths have a unique first intermediate room.
	// In other words, the second, third, etc. paths should not start with the same room as bestStepPath.
	var uniqueFirstPaths [][]string
	uniqueFirstPaths = append(uniqueFirstPaths, bestStepPath)
	for _, path := range bestStepDisjointPaths[1:] {
		// path[1] is the first intermediate room after start.
		if path[1] != bestStepPath[1] {
			duplicate := false
			// Also check among already accepted paths to ensure uniqueness.
			for _, up := range uniqueFirstPaths {
				if up[1] == path[1] {
					duplicate = true
					break
				}
			}
			if !duplicate {
				uniqueFirstPaths = append(uniqueFirstPaths, path)
			}
		}
	}
	bestStepDisjointPaths = uniqueFirstPaths

	// Log the final ordered disjoint paths.
	for i, path := range bestStepDisjointPaths {
		Log(fmt.Sprintf("Step Path %d : %v", i+1, path), "debug")
	}
}

// Sorts allPaths by the number of steps (path length) in ascending order.
func getSortedPathsBySteps() [][]string {
	sortedPaths := make([][]string, len(allPaths))
	copy(sortedPaths, allPaths)

	// Sort by length (fewest steps first).
	sort.Slice(sortedPaths, func(i, j int) bool {
		return len(sortedPaths[i]) < len(sortedPaths[j])
	})

	return sortedPaths
}

func selectDisjointPaths(paths [][]string, max int) [][]string {
	selected := [][]string{}
	used := map[string]int{}
	threshold := 1 // Set threshold to 0 for strict, or increase it to allow some overlap.

	for _, path := range paths {
		overlap := 0
		for _, room := range path[1 : len(path)-1] { // Exclude start and end.
			if used[room] > 0 {
				overlap++
			}
		}
		// Accept the path if the number of overlapping rooms is within the allowed threshold.
		if overlap <= threshold {
			selected = append(selected, path)
			for _, room := range path[1 : len(path)-1] {
				used[room]++
			}
			if len(selected) >= max {
				break
			}
		}
	}
	return selected
}
