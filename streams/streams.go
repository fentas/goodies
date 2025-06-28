package streams

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fentas/goodies/output"
	"github.com/fentas/goodies/util"
)

type IO struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	OutFlags output.Opts
	Closer   func() error
}

type printable interface {
	String() string
}

type ScrubSecrets interface {
	ScrubSecrets() interface{}
}

func (p *IO) Print(o interface{}) error {
	if s, ok := o.(printable); ok {
		_, err := fmt.Fprint(p.Out, s.String())
		return err
	}
	if s, ok := o.(ScrubSecrets); ok {
		o = s.ScrubSecrets()
	}

	switch {
	case p.OutFlags.IsSet("json"):
		return util.DescribeJSON(p.Out, o)
	case p.OutFlags.IsSet("format"):
		return util.DescribeFormat(p.Out, o, p.OutFlags["format"][0])
	default: // IsSet("yaml")
		return util.DescribeYAML(p.Out, o)
	}
}

func (s *IO) OutToFile(path string) (io.Writer, error) {
	var err error
	s.Out, s.Closer, err = ToFile(path)
	return s.Out, err
}
func (s *IO) ErrOutToFile(path string) (io.Writer, error) {
	var err error
	s.ErrOut, s.Closer, err = ToFile(path)
	return s.ErrOut, err
}
func (s *IO) ToFile(path string) (io.Writer, error) {
	out, err := s.OutToFile(path)
	s.ErrOut = out
	return out, err
}

func ToFile(path string) (io.Writer, func() error, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, nil, err
	}
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)
	go func() {
		_, _ = io.Copy(f, r)
		exit <- true
	}()

	// function to be deferred
	return w, func() error {
		_ = w.Close()
		<-exit
		return f.Close()
	}, nil
}
