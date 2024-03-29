package ssh

import (
	"bytes"
	"fmt"
	"strings"
)

func (s *Client) Output(cmds ...string) ([]byte, error) {
	if _, err := s.Connect(); err != nil {
		return nil, err
	}

	cmd := strings.Join(cmds, " && ")
	out, err := s.Session.Output(cmd)
	s.Close()
	return bytes.TrimSpace(out), err
}

func (s *Client) OutputPtrs(ptrs []*string) error {
	stdout, err := s.Output(ptrsToStrings(ptrs)...)
	if err != nil {
		return err
	}
	if err := ptrsCopyRawOutput(ptrs, stdout); err != nil {
		return err
	}
	return nil
}

func (s *Client) CombinedOutput(cmds ...string) ([]byte, error) {
	if _, err := s.Connect(); err != nil {
		return nil, err
	}

	cmd := strings.Join(cmds, " && ")
	out, err := s.Session.CombinedOutput(cmd)
	s.Close()
	return bytes.TrimSpace(out), err
}

func (s *Client) Run(cmds ...string) error {
	if _, err := s.Connect(); err != nil {
		return err
	}
	if s.Options.IOStreams == nil {
		s.logger.ErrorContext(s.ctx, "IOStreams is not set")
		return fmt.Errorf("IOStreams is not set")
	}

	s.Session.Stdout = s.Options.IOStreams.Out
	s.Session.Stderr = s.Options.IOStreams.ErrOut

	cmd := strings.Join(cmds, " && ")
	err := s.Session.Run(cmd)
	s.Close()
	return err
}

// func (s *Client) Run(command string) ([]byte, []byte, error) {
// 	if _, err := s.Connect(); err != nil {
// 		return nil, nil, err
// 	}

// 	var stdout, stderr bytes.Buffer
// 	s.Session.Stdout = &stdout
// 	s.Session.Stderr = &stderr
// 	if err := s.Session.Run(command); err != nil {
// 		return nil, nil, err
// 	}
// 	s.Close()
// 	return stdout.Bytes(), stderr.Bytes(), nil
// }

func ptrsCopyRawOutput(ptrs []*string, raw []byte) error {
	lines := strings.Split(string(raw), "\n")
	if len(lines) > len(ptrs) {
		return fmt.Errorf("unexpected number of lines in output: %d; expected: %d", len(lines), len(ptrs))
	}
	for i := 0; i < len(ptrs); i++ {
		value := ""
		if i < len(lines) {
			value = lines[i]
		}
		*ptrs[i] = value
	}
	return nil
}

func ptrsToStrings(ptrs []*string) []string {
	strs := make([]string, len(ptrs))
	for i, ptr := range ptrs {
		strs[i] = *ptr
	}
	return strs
}
