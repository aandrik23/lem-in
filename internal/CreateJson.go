package internal

import (
	"encoding/json"
	"os"
)

type SimulationDump struct {
	Start string `json:"start"`
	Rooms []Room `json:"rooms"`
	Moves []Move `json:"moves"`
}

func CreateJson() {
	if visualizer {
		dump := SimulationDump{
			Start: startRoom,
			Rooms: make([]Room, 0, len(rooms)),
			Moves: allMoves,
		}
		for _, r := range rooms {
			dump.Rooms = append(dump.Rooms, r)
		}

		f, err := os.Create("simulation.json")
		if err != nil {
			Log("could not create simulation.json: "+err.Error(), "error")
			return
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(dump); err != nil {
			Log("failed to write simulation.json: "+err.Error(), "error")
		}
	}
}
