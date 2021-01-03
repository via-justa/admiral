package mariadb

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/via-justa/admiral/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/go-sql-driver/mysql"
)

type viaSSHDialer struct {
	client *ssh.Client
}

// Dial implement ssh dialer
func (sshd *viaSSHDialer) Dial(addr string) (net.Conn, error) {
	return sshd.client.Dial("tcp", addr)
}

func newSSHProxy(conf config.SSHProxy) (*ssh.Client, error) {
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

// ProxyConnect starts new database connection via ssh proxy
func ProxyConnect(conf *config.Config) (*Database, error) {
	var db Database

	sshcon, err := newSSHProxy(conf.SSHProxy)
	if err != nil {
		return &db, err
	}

	mysql.RegisterDialContext("mysql+tcp", func(ctx context.Context, addr string) (net.Conn, error) {
		sshd := &viaSSHDialer{sshcon}
		return sshd.Dial(addr)
	})

	dbConfig := mysql.Config{
		User:                 conf.MariaDB.User,
		Passwd:               conf.MariaDB.Password,
		Net:                  "mysql+tcp",
		Addr:                 net.JoinHostPort(conf.MariaDB.Host, fmt.Sprint(conf.MariaDB.Port)),
		DBName:               conf.MariaDB.DB,
		AllowNativePasswords: true,
	}

	db.Conn, err = sqlx.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return &db, err
	}

	err = db.Conn.Ping()
	if err != nil {
		return &db, err
	}

	return &db, err
}
