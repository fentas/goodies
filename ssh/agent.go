package ssh

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func sshAgentAuth() ssh.AuthMethod {
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	a := agent.NewClient(conn)
	return ssh.PublicKeysCallback(a.Signers)
}
