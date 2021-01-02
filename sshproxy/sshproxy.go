package sshproxy

import (
	"fmt"
	"net"
	"os"

	"github.com/via-justa/admiral/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// NewSSHProxy starts new ssh connection to proxy commends thru
func NewSSHProxy(conf config.SSHProxy) (*ssh.Client, error) {
	var agentClient agent.Agent

	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close() // nolint: errcheck

		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            conf.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint:gosec
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}
	// When there's a non empty password add the password AuthMethod
	if conf.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return conf.Password, nil
		}))
	}

	// Connect to the SSH Server
	return ssh.Dial("tcp", net.JoinHostPort(conf.Host, fmt.Sprint(conf.Port)), sshConfig)
}
