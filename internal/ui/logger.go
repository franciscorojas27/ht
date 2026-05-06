package ui

import (
	"time"
	"github.com/fatih/color"
)

func TimeTaken(init time.Time) {
	if time.Since(init) > 2*time.Second {
		color.Red("Time taken: %s", time.Since(init))
	} else if time.Since(init) > 1*time.Second {
		color.Yellow("Time taken: %s", time.Since(init))
	} else {
		color.Green("Time taken: %s", time.Since(init))
	}
}