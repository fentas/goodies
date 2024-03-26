package progress

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/progress"
)

type ReaderCloser struct {
	body    io.ReadCloser
	tracker *progress.Tracker
}

func NewReader(body io.ReadCloser, tracker *progress.Tracker) *ReaderCloser {
	return &ReaderCloser{
		body:    body,
		tracker: tracker,
	}
}

func (r *ReaderCloser) Read(p []byte) (n int, err error) {
	r.tracker.Increment(int64(len(p)))
	return r.body.Read(p)
}

func (r *ReaderCloser) Close() error {
	return r.body.Close()
}
