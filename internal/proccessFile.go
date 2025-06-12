package internal

import (
	"fmt"
	"strconv"
	"strings"
)

func processLine(line string, numLine int) bool {
	// validating ants
	if numLine == 1 { // first line always number of ants
		antsNumber, err := strconv.Atoi(line)
		if err != nil {
			Log("number of ants must be a digit", "error")
			return false
		}
		if antsNumber <= 0 {
			Log("number of ants must be > 0", "error")
			return false
		}
		ants = antsNumber
		Log(fmt.Sprintf("Number of ants: %d", antsNumber), "debug")
		return true
	}

	//validating rooms
	if strings.HasPrefix(line, "##start") {
		if startRoomFound == false { // used to enter only once
			expectingStartRoom = true
			startRoomFound = true
			return true
		} else {
			Log("Found more than one start rooms", "error")
			return false
		}
	}

	if strings.HasPrefix(line, "##end") {
		if endRoomFound == false {
			expectingEndRoom = true
			endRoomFound = true
			return true
		} else {
			Log("Found more than one end rooms", "error")
			return false
		}
	}

	if isRoomLine(line) {
		ok := getRoom(line) // create the room
		if ok {
			return true
		}
		return false
	}

	//validating tunels
	if isTunnelLine(line) {
		ok := getTunnel(line) // link the rooms
		if ok {
			return true
		}
		return false
	}

	return true // unkown comments will be ignored
}

// Tunnel Functions----------------------------------------------------

func getTunnel(line string) bool {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		Log("invalid tunnel format: "+line, "error")
		return false
	}

	a := strings.TrimSpace(parts[0])
	b := strings.TrimSpace(parts[1])

	// Validate both rooms exist
	if _, ok := rooms[a]; !ok {
		Log("unknown room in tunnel: "+a, "error")
		return false
	}

	if _, ok := rooms[b]; !ok {
		Log("unknown room in tunnel: "+b, "error")
		return false
	}

	if a == b {
		Log("invalid tunnel: self-link on "+a, "error")
		return false
	}

	// Check for duplicate
	if isConnected(a, b) {
		Log("duplicate tunnel between "+a+" and "+b, "error")
		return false
	}

	// Add bidirectional link
	tunnels[a] = append(tunnels[a], b)
	tunnels[b] = append(tunnels[b], a)
	return true
}

func isConnected(a, b string) bool {
	for _, neighbor := range tunnels[a] {
		if neighbor == b {
			return true
		}
	}
	for _, neighbor := range tunnels[b] {
		if neighbor == a {
			return true
		}
	}
	return false
}

func isTunnelLine(line string) bool {
	return strings.Count(line, "-") == 1 && !strings.HasPrefix(line, "#")
}

// Room Functions------------------------------------------------------

func getRoom(line string) bool { // create room
	parts := strings.Fields(line) // getRoom assumes the line has already passed isRoomLine validation
	name := parts[0]
	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil {
		Log("invalid coordinates for room: "+line, "error")
		return false
	}

	if _, exists := rooms[name]; exists {
		Log("duplicate room name: "+name, "error")
		return false
	}

	if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
		Log("invalid room name: "+name, "error")
		return false
	}

	room := Room{Name: name, X: x, Y: y}
	rooms[name] = room

	if expectingStartRoom {
		startRoom = name
		expectingStartRoom = false
	}

	if expectingEndRoom {
		endRoom = name
		expectingEndRoom = false
	}
	return true
}

func isRoomLine(line string) bool {
	parts := strings.Fields(line)
	return len(parts) == 3 && !strings.HasPrefix(line, "#") && !strings.Contains(line, "-")
}

// ------------------------------------------------------
