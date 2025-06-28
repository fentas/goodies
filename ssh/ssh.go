package ssh

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/fentas/goodies/logs"
	"github.com/fentas/goodies/streams"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	defaultTimeout = 10 * time.Second
)

type ClientOptions struct {
	Host      string
	User      string
	Password  string
	IOStreams *streams.IO
	Timeout   time.Duration
}

type Client struct {
	ctx     context.Context
	logger  *slog.Logger
	Options *ClientOptions

	info    *RemoteInfo
	Session *ssh.Session

	sshIn  io.WriteCloser
	sshOut io.Reader
	sshErr io.Reader
}

func NewClient(ctx context.Context, logger *slog.Logger, options *ClientOptions) *Client {
	return &Client{
		Options: options,
		ctx:     ctx,
		logger: logger.
			With(slog.String("pkg", "ssh")).
			With(slog.String("host", options.Host)),
	}
}

func (s *Client) RawClient() (*ssh.Client, error) {
	authMethods := []ssh.AuthMethod{}
	s.authPassword(&authMethods)
	s.authInteractive(&authMethods)
	s.authGSSAPI(&authMethods)
	if conn := s.authAgent(&authMethods); conn != nil {
		s.logger.Log(s.ctx, logs.LevelTrace, "Connected to SSH_AUTH_SOCK")
		defer conn.Close()
	}

	if len(authMethods) == 0 {
		s.logger.ErrorContext(s.ctx, "No authentication methods available")
		return nil, fmt.Errorf("authentication methods exhausted")
	}

	timeout := s.Options.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(s.ctx, "tcp", s.Options.Host)
	if err != nil {
		s.logger.ErrorContext(s.ctx, "Failed to dial", "error", err)
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, s.Options.Host, &ssh.ClientConfig{
		User:            s.Options.User,
		Timeout:         timeout,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		s.logger.ErrorContext(s.ctx, "Failed to create client connection", "error", err)
		return nil, err
	}
	return ssh.NewClient(c, chans, reqs), nil
}

// Connect to the server, open a session and run a command.
func (s *Client) Connect() (*ssh.Session, error) {
	if s.Session != nil {
		return s.Session, nil
	}

	client, err := s.RawClient()
	if err != nil {
		return nil, err
	}

	s.Session, err = client.NewSession()
	if err != nil {
		s.logger.ErrorContext(s.ctx, "Failed to create session", "error", err)
		return nil, err
	}

	return s.Session, err
}

func (s *Client) authPassword(auth *[]ssh.AuthMethod) {
	if s.Options.Password == "" {
		s.logger.Log(s.ctx, logs.LevelTrace, "No password provided")
		return
	}

	*auth = append(*auth, ssh.Password(s.Options.Password))
}
func (s *Client) authAgent(auth *[]ssh.AuthMethod) net.Conn {
	if os.Getenv("SSH_AUTH_SOCK") == "" {
		s.logger.Log(s.ctx, logs.LevelTrace, "No SSH_AUTH_SOCK")
		return nil
	}

	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		s.logger.InfoContext(s.ctx, "Failed to connect to SSH_AUTH_SOCK", "error", err)
		return nil
	}
	*auth = append(*auth, ssh.PublicKeysCallback(
		agent.NewClient(conn).Signers,
	))
	return conn
}
func (s *Client) authInteractive(auth *[]ssh.AuthMethod) {
	// todo
}
func (s *Client) authGSSAPI(auth *[]ssh.AuthMethod) {
	// todo
}

func (s *Client) Close() {
	if s.Session != nil {
		_ = s.Session.Close()
		s.Session = nil
		return
	}
	if s.sshIn != nil {
		_ = s.sshIn.Close()
	}
}
