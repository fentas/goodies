package logs

import (
	"log/slog"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

type progressTracker struct {
	*progress.Tracker
	done chan struct{}
	msg  chan string
}

func ProgressTracker(logger *slog.Logger, tracker *progress.Tracker) {
	h := logger.Handler().(*Handler)
	t := &progressTracker{
		Tracker: tracker,
		done:    make(chan struct{}),
		msg:     make(chan string),
	}
	h.tracker = t
	go func() {
		for {
			select {
			case <-h.tracker.done:
				return
			case msg := <-h.tracker.msg:
				h.m.Lock()
				if h.tracker == nil {
					h.m.Unlock()
					return
				}
				h.tracker.UpdateMessage(msg)
				h.m.Unlock()
			}
		}
	}()
}

func ProgressDone(logger *slog.Logger, message string, err error) {
	h := logger.Handler().(*Handler)
	if h.tracker != nil {
		if err != nil {
			h.tracker.UpdateMessage(err.Error())
			h.tracker.MarkAsErrored()
		} else {
			h.tracker.UpdateMessage(message)
			h.tracker.MarkAsDone()
		}
		h.tracker.done <- struct{}{}
		h.tracker = nil
		time.Sleep(time.Millisecond * 100)
	}
}
