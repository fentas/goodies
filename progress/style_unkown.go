package progress

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

var (
	StyleUnkown = NewStyle("unkown", progress.StyleVisibility{
		ETA:            false,
		ETAOverall:     false,
		Percentage:     false,
		Pinned:         false,
		Speed:          false,
		SpeedOverall:   false,
		Time:           true,
		Tracker:        true,
		TrackerOverall: false,
		Value:          false,
	}, &StyleOverwrites{
		TimeInProgressPrecision: time.Second,
	})
)
