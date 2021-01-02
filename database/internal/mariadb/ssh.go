package mariadb

import (
	"context"
	"fmt"
	"net"

	"github.com/via-justa/admiral/sshproxy"

	"github.com/jmoiron/sqlx"
	"github.com/via-justa/admiral/config"
	"golang.org/x/crypto/ssh"

	"github.com/go-sql-driver/mysql"
)

type viaSSHDialer struct {
	client *ssh.Client
}

func (sshd *viaSSHDialer) Dial(addr string) (net.Conn, error) {
	return sshd.client.Dial("tcp", addr)
}

// ProxyConnect starts new database connection via ssh proxy
func ProxyConnect(conf *config.Config) (*Database, error) {
	var db Database

	sshcon, err := sshproxy.NewSSHProxy(conf.SSHProxy)
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
