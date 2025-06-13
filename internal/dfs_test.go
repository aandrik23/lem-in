package internal

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Minimal helpers to reset state for each test run
func resetGlobalsForTest() {
	tunnels = make(map[string][]string)
	rooms = make(map[string]Room)
	startRoom = ""
	endRoom = ""
	allPaths = nil
	allMoves = nil
	lastRoom = make(map[int]string)
}

// Helper to sort paths for comparison
func sortPaths(paths [][]string) {
	sort.Slice(paths, func(i, j int) bool {
		return joinPath(paths[i]) < joinPath(paths[j])
	})
}

func joinPath(path []string) string {
	return "/" + strings.Join(path, "/")
}

func TestFindAllPaths(t *testing.T) {
	resetGlobalsForTest()
	rooms = map[string]Room{
		"start": {Name: "start"},
		"A":     {Name: "A"},
		"B":     {Name: "B"},
		"C":     {Name: "C"},
		"end":   {Name: "end"},
	}
	tunnels = map[string][]string{
		"start": {"A", "B"},
		"A":     {"C"},
		"B":     {"C"},
		"C":     {"end"},
		"end":   {},
	}
	startRoom = "start"
	endRoom = "end"

	FindAllPaths()

	want := [][]string{
		{"start", "A", "C", "end"},
		{"start", "B", "C", "end"},
	}

	sortPaths(allPaths)
	sortPaths(want)
	if !reflect.DeepEqual(allPaths, want) {
		t.Errorf("FindAllPaths() = %v, want %v", allPaths, want)
	}
}
