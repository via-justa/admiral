package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(copy)

	copy.AddCommand(copyHostVar)
	copy.AddCommand(copyGroupVar)
}

var copy = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"cp"},
	Short:   "create a new record from existing record",
}

var copyHostVar = &cobra.Command{
	Use:   "host",
	Short: "create a new host from existing one",
	Long:  "Use existing host record as template while creating a new host record",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host
		var existingHost datastructs.Host
		var host datastructs.Host
		var err error

		switch len(args) {
		case 0, 1:
			log.Print("please set source and destination host arguments")
			return
		case 2:
			hosts, err = returnHosts(args[0])
			if err != nil {
				log.Print(err)
				return
			}

			if len(hosts[0].Hostname) != 0 {
				existingHost = hosts[0]
				existingHost.Host = ""
				checkedVal := strings.SplitN(args[1], ".", 2)
				existingHost.Hostname = checkedVal[0]
				if len(checkedVal) > 1 {
					existingHost.Domain = checkedVal[1]
				}
			}
		default:
			log.Print("received too many arguments")
			return
		}

		host, err = editHost(&existingHost, args[1])
		if err != nil {
			log.Print(err)
		}

		printHosts([]datastructs.Host{host})

		if confirm() {
			err = confirmedHost(&host)
			if err != nil {
				log.Print(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}

var copyGroupVar = &cobra.Command{
	Use:   "group",
	Short: "create a new group from existing one",
	Long:  "Use existing group record as template while creating a new group record",
	Run: func(cmd *cobra.Command, args []string) {
		var templateGroup datastructs.Group
		var group datastructs.Group
		var err error

		switch len(args) {
		case 0, 1:
			log.Print("please set source and destination group arguments")
			return
		case 2:
			templateGroup, err = viewGroupByName(args[0])
			if err != nil {
				log.Print(err)
				return
			}

			templateGroup.Name = args[1]
		default:
			log.Fatal("received too many arguments")
		}

		group, err = editGroup(&templateGroup, args[1])
		if err != nil {
			log.Print(err)
			return
		}

		printGroups([]datastructs.Group{group})
		if confirm() {
			err := createGroup(&group)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}
