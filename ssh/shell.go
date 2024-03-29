package ssh

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// todo: if this should work we need
// - ignore start info of remote shell
// - chennel to check output if command is done

func (s *Client) Shell() error {
	if _, err := s.Connect(); err != nil {
		return err
	}

	var err error
	s.sshIn, err = s.Session.StdinPipe()
	if err != nil {
		return err
	}
	s.sshOut, err = s.Session.StdoutPipe()
	if err != nil {
		return err
	}
	s.sshErr, err = s.Session.StderrPipe()
	if err != nil {
		return err
	}

	return s.Session.Shell()
}

func (s *Client) Exec(commands ...string) error {
	if s.Session == nil || s.Session.Stdin == nil {
		s.Shell()
	}

	commands = append(commands, "exit")
	for _, command := range commands {
		command = fmt.Sprintf("%s\r\n", strings.TrimSpace(command))
		fmt.Println(command)
		if _, err := io.WriteString(s.sshIn, command); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}

func (s *Client) ReadStdout() ([]byte, error) {
	if s.sshOut == nil {
		return nil, fmt.Errorf("no stdout")
	}

	return io.ReadAll(s.sshOut)
}

func (s *Client) ReadStderr() ([]byte, error) {
	if s.sshErr == nil {
		return nil, fmt.Errorf("no stderr")
	}

	return io.ReadAll(s.sshErr)
}

func (s *Client) CombinedRead() ([]byte, error) {
	if s.sshOut == nil || s.sshErr == nil {
		return nil, fmt.Errorf("no stdout or stderr")
	}

	stdout, err := io.ReadAll(s.sshOut)
	if err != nil {
		return nil, err
	}

	stderr, err := io.ReadAll(s.sshErr)
	if err != nil {
		return nil, err
	}

	return append(stdout, stderr...), nil
}
