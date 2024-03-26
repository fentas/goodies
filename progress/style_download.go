package progress

import (
	"github.com/jedib0t/go-pretty/v6/progress"
)

var (
	StyleDownload = NewStyle("download", progress.StyleVisibility{
		ETA:            true,
		ETAOverall:     true,
		Percentage:     true,
		Pinned:         false,
		Speed:          true,
		SpeedOverall:   true,
		Time:           true,
		Tracker:        true,
		TrackerOverall: false,
		Value:          true,
	}, nil)
)
