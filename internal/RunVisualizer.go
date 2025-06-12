package internal

import (
	"os"
	"os/exec"
)

var jsonFile = "simulation.json"

func RunVisualizer() {
	if visualizer {
		// 1) Run the Python visualizer
		cmd := exec.Command("python3", "python/visualizer.py", "--input", jsonFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		// 2) Check exit status
		if err != nil {
			Log("visualizer failed: "+err.Error(), "error")
		} else {
			// 3) On success, delete the JSON file
			if rmErr := os.Remove(jsonFile); rmErr != nil {
				// silently ignore or log at debug level
				Log("could not remove JSON file: "+rmErr.Error(), "debug")
			}
		}
	}
}
