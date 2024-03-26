package progress

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

func ProgressDone(tracker *progress.Tracker, message string, err error) {
	if err != nil {
		tracker.UpdateMessage(err.Error())
		tracker.MarkAsErrored()
	} else {
		tracker.UpdateMessage(message)
		tracker.MarkAsDone()
	}
	time.Sleep(time.Millisecond * 100)
}
