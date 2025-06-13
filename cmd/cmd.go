package cmd

import (
	"fmt"
	"lemin/internal"
	"time"
)

func Cmd() {
	start := time.Now()

	file := internal.CheckArguments()
	internal.ValidateFileFormat(file)
	internal.ValidateConnectivity()
	internal.FindAllPaths()
	internal.Simulate()
	internal.CreateJson()
	internal.RunVisualizer()

	elapsed := time.Since(start)
	internal.Log(fmt.Sprintf("Execution time: %s\n", elapsed), "debug")
}
