package internal

import (
	"fmt"
)

func Log(s string, errType string) {
	switch errType {
	case "error":
		fmt.Printf("[ERROR] %s\n", s)
	case "debug":
		// return // comment this line to show debug output
		fmt.Printf("[DEBUG] %s\n", s)
	}
}
