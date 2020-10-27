package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/cli"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(view)

	view.AddCommand(viewHost)
	view.AddCommand(viewGroup)
	view.AddCommand(viewChild)
}

var view = &cobra.Command{
	Use:     "view",
	Aliases: []string{"list", "ls", "get"},
	Short:   "view existing record",
}

var viewHost = &cobra.Command{
	Use:   "host",
	Short: "view existing host by substring of hostname or IP or view all records when no argument passed",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host

		var err error

		switch len(args) {
		case 0:
			hosts, err = cli.ListHosts()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			hosts, err = cli.ScanHosts(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		printHosts(hosts)
	},
}

var viewGroup = &cobra.Command{
	Use:   "group",
	Short: "view existing group by substring of group name or view all records when no argument passed",
	Run: func(cmd *cobra.Command, args []string) {
		var groups []datastructs.Group

		var err error

		switch len(args) {
		case 0:
			groups, err = cli.ListGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			groups, err = cli.ScanGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		printGroups(groups)
	},
}

var viewChild = &cobra.Command{
	Use:   "child",
	Short: "view existing child-group relationship by parent or/and child or view all records when no argument passed",
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup

		var err error

		switch len(args) {
		case 0:
			childGroups, err = cli.ListChildGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			childGroups, err = cli.ScanChildGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		printChildGroups(childGroups)
	},
}
