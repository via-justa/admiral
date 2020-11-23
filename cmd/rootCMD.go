package cmd

import (
	"log"
	"os"

	"github.com/via-justa/admiral/config"
	"github.com/via-justa/admiral/database"

	"github.com/spf13/cobra"
)

var (
	// AppVersion is set in build time to the latest application version
	AppVersion string
	// Conf contain default configuration settings
	Conf *config.Config
	// DB connection to selected database backend
	DB database.DBInterface
	// User implement user action confirmation
	User userInt
)

var (
	rootCmd = &cobra.Command{
		Use:        "admiral command",
		ValidArgs:  []string{"copy", "create", "edit", "delete", "view", "list", "inventory", "prometheus"},
		ArgAliases: []string{"cp", "add", "remove", "rm", "del", "ls", "get", "inv", "prom"},
		Short:      "Admiral is a lightweight Ansible inventory database management tool",
		Long: `Admiral is a command line tool to manage ansible inventory. It can also 
expose the inventory to ansible as a full inventory structure. As monitoring is 
also important, the tool can also expose the inventory in Prometheus static file 
structure where all the host groups are set as host 'groups' label.

The tool is expecting to find a toml configuration file with the database details
in one of the following locations:
- /etc/admiral/config.toml
- ./config.toml
- $HOME/.admiral.toml

Example configuration file:
[database]
user = "root"
password = "local"
host = "localhost:3306"
db = "ansible"`,
	}
)

// nolint:errcheck
func init() {
	log.SetFlags(0)
}

//Execute starts the program
func Execute() {
	var err error

	Conf = config.NewConfig()
	User = newUser()

	if os.Args[1] != "docs" && os.Args[1] != "completion" {
		DB, err = database.Connect(Conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
