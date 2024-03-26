package progress

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	StyleDefaultChars = progress.StyleChars{
		BoxLeft:       "[",
		BoxRight:      "]",
		Finished:      "#",
		Finished25:    ".",
		Finished50:    ".",
		Finished75:    ".",
		Indeterminate: progress.IndeterminateIndicatorMovingBackAndForth("<#>", progress.DefaultUpdateFrequency/2),
		Unfinished:    ".",
	}
	StyleDefaultColors = progress.StyleColors{
		Message: text.Colors{text.FgWhite},
		Error:   text.Colors{text.FgRed},
		Percent: text.Colors{text.FgHiRed},
		Pinned:  text.Colors{text.BgHiBlack, text.FgWhite, text.Bold},
		Stats:   text.Colors{text.FgHiBlack},
		Time:    text.Colors{text.FgGreen},
		Tracker: text.Colors{text.FgYellow},
		Value:   text.Colors{text.FgCyan},
		Speed:   text.Colors{text.FgMagenta},
	}
	StyleDefaultOptions = progress.StyleOptions{
		DoneString:              "done!",
		ErrorString:             "fail!",
		ETAPrecision:            time.Second,
		ETAString:               "~ETA",
		PercentFormat:           "%4.2f%%",
		PercentIndeterminate:    " ??? ",
		Separator:               " ... ",
		SnipIndicator:           "~",
		SpeedPosition:           progress.PositionRight,
		SpeedPrecision:          time.Microsecond,
		SpeedOverallFormatter:   progress.FormatNumber,
		SpeedSuffix:             "/s",
		TimeDonePrecision:       time.Millisecond,
		TimeInProgressPrecision: time.Microsecond,
		TimeOverallPrecision:    time.Second,
	}
)

type StyleOverwrites struct {
	TimeInProgressPrecision time.Duration
}

func NewStyle(name string, visibility progress.StyleVisibility, overwrites *StyleOverwrites) progress.Style {
	style := progress.Style{
		Name:       name,
		Chars:      StyleDefaultChars,
		Colors:     StyleDefaultColors,
		Options:    StyleDefaultOptions,
		Visibility: visibility,
	}
	if overwrites != nil {
		style.Options.TimeInProgressPrecision = overwrites.TimeInProgressPrecision
	}
	return style
}
