package progress

import (
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

type Writer struct {
	progress.Writer
}

func NewWriter(style progress.Style, out io.Writer) *Writer {
	pw := progress.NewWriter()
	pw.SetOutputWriter(out)
	pw.SetAutoStop(false)
	pw.SetTrackerLength(20)
	pw.SetSortBy(progress.SortByPercent)
	pw.SetTrackerPosition(progress.PositionRight)

	pw.SetStyle(style)
	switch style.Name {
	case StyleUnkown.Name:
		pw.SetUpdateFrequency(time.Millisecond * 500)
		pw.SetTrackerPosition(progress.PositionLeft)
	case StyleDownload.Name:
		pw.SetUpdateFrequency(time.Millisecond * 200)
		pw.SetMessageLength(30)
	}

	return &Writer{
		Writer: pw,
	}
}

func (w *Writer) AddTracker(message string, total int64) *progress.Tracker {
	units := progress.UnitsDefault
	switch w.Style().Name {
	case StyleDownload.Name:
		units = progress.UnitsBytes
	}

	tracker := progress.Tracker{
		Message: message,
		Total:   total,
		Units:   units,
	}

	w.Writer.AppendTracker(&tracker)
	return &tracker
}
