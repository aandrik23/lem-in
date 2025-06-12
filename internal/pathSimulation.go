package internal

import (
	"fmt"
	"sort"
	"strings"
)

func Simulate() {
	// Determine the number of paths to use based on the number of ants
	numPaths := 1
	if ants > 3 && ants <= 6 {
		numPaths = 3
	} else {
		numPaths = len(bestStepDisjointPaths)
	}

	// Clamp to available paths
	if numPaths > len(bestStepDisjointPaths) {
		numPaths = len(bestStepDisjointPaths)
	}

	// 1) slice out the paths we will actually use
	paths := bestStepDisjointPaths[:numPaths]

	// 2) compute each path’s “cost” (number of edges)
	costs := make([]int, len(paths))
	for i, p := range paths {
		costs[i] = len(p) - 1
	}

	// 3) compute exactly how many ants each path should carry
	antsPerPath := ComputeAntsPerPath(costs, ants)

	// 4) hand off to simulateAnts (now with quotas)
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

	// Process ants in reverse order so that ants closer to the end are processed first.
	for i := len(antsInTransit) - 1; i >= 0; i-- {
		ant := antsInTransit[i]
		path := paths[ant.Path]
		currentRoom := path[ant.Index]
		nextIndex := ant.Index + 1

		// Free the current room immediately, since the ant is going to try to leave.
		delete(occupied, currentRoom)

		// Check if there is a valid next room and if it is not occupied.
		// Check if there is a next room and if it is not occupied.
		if nextIndex < len(path) && !occupied[path[nextIndex]] {
			move := fmt.Sprintf("L%d-%s", ant.ID, path[nextIndex])
			output = append(output, move)
			ant.Index = nextIndex

			// If this ant has not yet reached the final room, add it back to transit.
			if nextIndex < len(path)-1 {
				newTransit = append(newTransit, ant)
				occupied[path[nextIndex]] = true // Reserve the room.
			}
			// If the ant reaches the final room, do not add it back.
		} else {
			// If the ant couldn’t move, re-reserve its current room and keep it in transit.
			if currentRoom != endRoom {
				occupied[currentRoom] = true
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

	// These values manage the ant simulation state.
	antsInTransit := []Ant{}
	spawned := make([]int, len(paths))
	nextAnt := 1
	turns := 0

	// occupied tracks the rooms in use during the current turn.
	occupied := make(map[string]bool)

	// The simulation loop runs until no moves are produced.
	for {

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

		// If no moves were made this turn, the simulation is complete.
		if len(turnOutput) == 0 {
			break
		}

		// Sort moves by ant ID for consistent ordering in output.
		sort.Slice(turnOutput, func(i, j int) bool {
			var id1, id2 int
			fmt.Sscanf(turnOutput[i], "L%d-", &id1)
			fmt.Sscanf(turnOutput[j], "L%d-", &id2)
			return id1 < id2
		})
		if visualizer {
			// For every move string, parse it and append to allMoves.
			for _, str := range turnOutput {
				antID, toRoom := parseMove(str)
				fromRoom := lastRoom[antID] // lookup where it was
				allMoves = append(allMoves, Move{
					Turn: turns + 1,
					Ant:  antID,
					From: fromRoom,
					To:   toRoom,
				})
				lastRoom[antID] = toRoom // update for next time
			}
		} else {
			// Legacy behavior
			fmt.Println(strings.Join(turnOutput, " "))
		}

		turns++

		// Reset the occupied map for the next turn.
		occupied = make(map[string]bool)

	}

	// Log the total number of turns (only count turns in which moves were executed).
	Log(fmt.Sprintf("Total number of turns: %d\n", turns), "debug")
}

// parseMove takes a string of the form "L<antID>-<roomName>"
// and returns the integer antID and the roomName.
func parseMove(s string) (antID int, room string) {
	// Split on the first dash into ["L3", "h"]
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		// malformed; you could choose to log or panic here instead
		return 0, ""
	}
	// Parse the ant number from the "L3" piece
	fmt.Sscanf(parts[0], "L%d", &antID)
	// The second part is the room name
	room = parts[1]
	return
}
