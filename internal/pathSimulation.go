package internal

import (
	"fmt"
	"sort"
	"strings"
)

func Simulate() {
	paths, antsPerPath := SelectOptimalDisjointPathSet(ants, allPaths)
	if len(paths) == 0 {
		Log("No valid disjoint paths to simulate.", "error")
		return
	}
	simulateAnts(paths, antsPerPath)
}

// Ant represents the state of an ant in the simulation.
type Ant struct {
	ID    int // Ant identifier
	Path  int // Index into the paths slice (which path the ant is following)
	Index int // Current index (position) on that path
}

// moveAntsInTransit processes ants in transit so that their current room is freed
// immediately as they start to move. It returns the updated list of ants still in transit
// along with the movement commands for this turn.
func moveAntsInTransit(antsInTransit []Ant, paths [][]string, occupied map[string]bool) ([]Ant, []string) {
	output := []string{}
	newTransit := []Ant{}

	for i := 0; i < len(antsInTransit); i++ {
		ant := antsInTransit[i]
		path := paths[ant.Path]
		currentRoom := path[ant.Index]
		nextIndex := ant.Index + 1

		// 1. Free the current room IMMEDIATELY (unless it's start or end)
		if currentRoom != startRoom && currentRoom != endRoom {
			delete(occupied, currentRoom)
		}

		// 2. Then try to move as before...
		if nextIndex < len(path) {
			nextRoom := path[nextIndex]
			if !occupied[nextRoom] {
				move := fmt.Sprintf("L%d-%s", ant.ID, nextRoom)
				output = append(output, move)
				ant.Index = nextIndex

				if nextIndex < len(path)-1 {
					newTransit = append(newTransit, ant)
				}
				// Mark nextRoom as occupied only if not start/end
				if nextRoom != startRoom && nextRoom != endRoom {
					occupied[nextRoom] = true
				}
			} else {
				// If can't move, reoccupy current room
				if currentRoom != startRoom && currentRoom != endRoom {
					occupied[currentRoom] = true
				}
				newTransit = append(newTransit, ant)
			}
		}
	}

	return newTransit, output
}

// simulateAnts is the main simulation function. It repeatedly:
//  1. Moves ants already in transit,
//  2. Spawns new ants,
//  3. Prints all moves for that turn,
//  4. And finally updates the turn count.
//
// The simulation stops when no moves occur on a turn.
func simulateAnts(paths [][]string, quota []int) {
	if len(paths) == 0 {
		Log("no valid paths to simulate.", "error")
		return
	}

	antsInTransit := []Ant{}
	spawned := make([]int, len(paths))
	nextAnt := 1
	turns := 0

	for {
		// Reconstruct occupation map based on ants currently in rooms (excluding start/end)
		occupied := make(map[string]bool)
		for _, ant := range antsInTransit {
			path := paths[ant.Path]
			currentRoom := path[ant.Index]
			if currentRoom != startRoom && currentRoom != endRoom {
				occupied[currentRoom] = true
			}
		}

		turnOutput := []string{}

		// Move ants already in transit.
		var moves []string
		antsInTransit, moves = moveAntsInTransit(antsInTransit, paths, occupied)
		turnOutput = append(turnOutput, moves...)
		// spawn according to quota
		for i := range paths {
			room := paths[i][1] // the first room after start
			canSpawn := spawned[i] < quota[i]
			// if it's not the end, also require that it's free
			if room != endRoom {
				canSpawn = canSpawn && !occupied[room]
			}
			if !canSpawn {
				continue
			}

			// 1) mark that we've spawned one more on this path
			spawned[i]++

			// 2) enqueue the ant so it moves in the next phase
			antsInTransit = append(antsInTransit, Ant{
				ID:    nextAnt,
				Path:  i,
				Index: 1,
			})

			// 3) record the move (to be printed or turned into JSON)
			move := fmt.Sprintf("L%d-%s", nextAnt, room)
			turnOutput = append(turnOutput, move)

			// 4) reserve the room **only if** it's not the end
			if room != endRoom {
				occupied[room] = true
			}

			// 5) bump the ant ID counter
			nextAnt++
		}

		if len(turnOutput) == 0 {
			break
		}

		sort.Slice(turnOutput, func(i, j int) bool {
			var id1, id2 int
			fmt.Sscanf(turnOutput[i], "L%d-", &id1)
			fmt.Sscanf(turnOutput[j], "L%d-", &id2)
			return id1 < id2
		})
		if visualizer {
			for _, str := range turnOutput {
				antID, toRoom := parseMove(str)
				fromRoom := lastRoom[antID]
				allMoves = append(allMoves, Move{
					Turn: turns + 1,
					Ant:  antID,
					From: fromRoom,
					To:   toRoom,
				})
				lastRoom[antID] = toRoom
			}
		} else {
			fmt.Println(strings.Join(turnOutput, " "))
		}

		turns++
	}

	Log(fmt.Sprintf("Total number of turns: %d\n", turns), "debug")
}

// parseMove takes a string of the form "L<antID>-<roomName>"
// and returns the integer antID and the roomName.
func parseMove(s string) (antID int, room string) {
	// Split on the first dash into ["L3", "h"]
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		// malformed;
		return 0, ""
	}
	// Parse the ant number from the "Lx" piece
	fmt.Sscanf(parts[0], "L%d", &antID)
	// The second part is the room name
	room = parts[1]
	return
}
