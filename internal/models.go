package internal

var (
	visualizer = false // default visualization is off data is printed on terminal
	ants       int     // number of ants

	expectingStartRoom bool
	startRoomFound     = false

	expectingEndRoom bool
	endRoomFound     = false

	startRoom string // saves start room name
	endRoom   string // saves quick room name
)

// tunnels created after reading the file
var tunnels = make(map[string][]string)

// used for quick room lookup
var rooms = make(map[string]Room)

// Rooms created after reading the file
type Room struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// contains all the paths from DFS
var allPaths [][]string

var (
	bestStepPath          []string   // best path from step calculator
	bestStepDisjointPaths [][]string // includes bestStepPath and others
)

//-------------------------------------------------------------------------
/// Create json for python visualizer

// Move is a JSON-serializable record of a single ant move.
type Move struct {
	Turn int    `json:"turn"`
	Ant  int    `json:"ant"`
	From string `json:"from"`
	To   string `json:"to"`
}

// allMoves accumulates every Move in visualizer mode.
var allMoves []Move

// lastRoom[antID] = the room that antID occupied at the end of the last turn.
var lastRoom = make(map[int]string)
