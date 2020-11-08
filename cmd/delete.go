package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(delete)

	delete.AddCommand(deleteHostVar)
	delete.AddCommand(deleteGroupVar)
	delete.AddCommand(deleteChildVar)
}

var delete = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm", "del"},
	Short:   "delete existing record",
}

var deleteHostVar = &cobra.Command{
	Use:     "host",
	Short:   "delete existing host",
	Example: "admiral delete host host1",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host

		var err error

		switch len(args) {
		case 0:
			log.Fatal("no host argument passed")
		case 1:
			hosts, err = scanHosts(args[0])
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
				affected, err := deleteHost(&hosts[0])
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

func deleteHost(host *datastructs.Host) (affected int64, err error) {
	affected, err = db.deleteHost(host)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}

var deleteGroupVar = &cobra.Command{
	Use:     "group",
	Short:   "delete existing group",
	Example: "admiral delete group group1",
	Run: func(cmd *cobra.Command, args []string) {
		var groups []datastructs.Group

		var err error

		switch len(args) {
		case 0:
			log.Fatal("no group argument passed")
		case 1:
			groups, err = scanGroups(args[0])
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
				affected, err := deleteGroup(&groups[0])
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

func deleteGroup(group *datastructs.Group) (affected int64, err error) {
	affected, err = db.deleteGroup(group)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}

var deleteChildVar = &cobra.Command{
	Use:     "child",
	Short:   "delete existing child-group relationship",
	Long:    "delete existing child-group relationship expecting ordered arguments child and parent group names",
	Example: "admiral delete child child-group parent-group",
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup

		var err error

		childGroups, err = viewChildGroup(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		printChildGroups(childGroups)
		if confirm() {
			affected, err := deleteChildGroup(&childGroups[0])
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

func deleteChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error) {
	affected, err = db.deleteChildGroup(childGroup)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}
