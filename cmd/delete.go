package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/cli"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(delete)

	delete.AddCommand(deleteHost)
	delete.AddCommand(deleteGroup)
	delete.AddCommand(deleteChild)
}

var delete = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm", "del"},
	Short:   "delete existing record",
}

var deleteHost = &cobra.Command{
	Use:   "host",
	Short: "delete existing host",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host

		var err error

		switch len(args) {
		case 0:
			log.Fatal("no host argument passed")
		case 1:
			hosts, err = cli.ScanHosts(args[0])
			if err != nil {
				log.Print(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		switch len(hosts) {
		case 0:
			log.Fatal("no host matched request")
		case 1:
			printHosts(hosts)
			if confirm() {
				affected, err := cli.DeleteHost(&hosts[0])
				if err != nil {
					log.Fatal(err)
				} else {
					fmt.Printf("Lines deleted %v\n", affected)
				}
			} else {
				log.Fatal("aborted")
			}
		default:
			printHosts(hosts)
			log.Fatal("too many results please adjust your request")
		}
	},
}

var deleteGroup = &cobra.Command{
	Use:   "group",
	Short: "delete existing group",
	Run: func(cmd *cobra.Command, args []string) {
		var groups []datastructs.Group

		var err error

		switch len(args) {
		case 0:
			log.Fatal("no group argument passed")
		case 1:
			groups, err = cli.ScanGroups(args[0])
			if err != nil {
				log.Print(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		switch len(groups) {
		case 0:
			log.Fatal("no group matched request")
		case 1:
			printGroups(groups)
			if confirm() {
				affected, err := cli.DeleteGroup(&groups[0])
				if err != nil {
					log.Fatal(err)
				} else {
					fmt.Printf("Lines deleted %v\n", affected)
				}
			} else {
				log.Fatal("aborted")
			}
		default:
			printGroups(groups)
			log.Fatal("too many results please adjust your request")
		}
	},
}

var deleteChild = &cobra.Command{
	Use:   "child",
	Short: "delete existing child-group relationship",
	Long:  "delete existing child-group relationship expecting ordered arguments child and parent group names",
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup

		var err error

		childGroups, err = cli.ViewChildGroup(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		printChildGroups(childGroups)
		if confirm() {
			affected, err := cli.DeleteChildGroup(&childGroups[0])
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("Lines deleted %v\n", affected)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}
