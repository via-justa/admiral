package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	AppVersion string
	DebugLevel bool

	jsonPath  string
	name      string
	enable    bool
	monitor   bool
	variables string
	id        int
)

var (
	rootCmd = &cobra.Command{
		Use:   "admiral",
		Short: "Admiral is a lightweight Ansible inventory database management tool",
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
	rootCmd.PersistentFlags().BoolVar(&DebugLevel, "debug", false, "Sets log level to debug")
	rootCmd.PersistentFlags().StringVar(&jsonPath, "file", "", "Path to JSON encoded file")
}

//Execute starts the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}