package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ------------------------------------------------------
func GetFile() string {
	if len(os.Args) < 2 {
		Log("insufficient arguments", "error")
		os.Exit(0)
	}

	var (
		file      string
		fileFound = false
	)

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-v", "--visualize":
			visualizer = true

		default:
			if strings.HasPrefix(arg, "-") {
				Log(fmt.Sprintf("unknown flag %q", arg), "error")
				os.Exit(0)
			}
			if fileFound {
				Log("too many positional arguments", "error")
				os.Exit(0)
			}
			file = arg
			fileFound = true
		}
	}
	if !fileFound {
		Log("no input file specified", "error")
		os.Exit(0)
	}
	// check if file exists
	_, err := os.Stat(file)
	if err != nil {
		Log("file does not exist", "error")
		os.Exit(0)
	}
	return file
}

// ------------------------------------------------------
func ValidateFileFormat(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		Log("failed to open file", "error")
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	numLine := 1
	for scanner.Scan() {
		line := scanner.Text()
		Continue := processLine(line, numLine)
		if Continue != true {
			Log("line "+strconv.Itoa(numLine)+": "+line, "error")
			os.Exit(0) // we already specify the errors so we exiting with a 0 state
		}
		numLine++
	}

	// make sure start && end rooms are not empty strings
	if strings.TrimSpace(startRoom) == "" {
		Log("no start room found.", "error")
		os.Exit(0)
	}
	if strings.TrimSpace(endRoom) == "" {
		Log("no end room found.", "error")
		os.Exit(0)
	}
	Log("Starting Room: "+startRoom+" Ending Room: "+endRoom, "debug")
}
