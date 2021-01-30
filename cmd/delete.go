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
	Use:        "delete",
	Aliases:    []string{"remove", "rm", "del"},
	ValidArgs:  []string{"host", "group", "child"},
	ArgAliases: []string{"hosts", "groups"},
	Short:      "delete existing record",
}

var deleteHostVar = &cobra.Command{
	Use:               "host hostname",
	Short:             "delete existing host",
	Example:           "admiral delete host host1",
	ValidArgsFunction: hostsArgsFunc,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteHostCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func deleteHostCase(args []string) error {
	var hosts []datastructs.Host

	var err error

	hosts, err = scanHosts(args[0])
	if err != nil {
		return err
	}

	switch len(hosts) {
	case 0:
		return fmt.Errorf("no host matched request")
	case 1:
		printHosts(hosts)

		if User.confirm() {
			affected, err := deleteHost(&hosts[0])
			if err != nil {
				return err
			}

			fmt.Printf("lines deleted %v\n", affected)
		} else {
			return fmt.Errorf("aborted")
		}
	default:
		printHosts(hosts)
		return fmt.Errorf("too many results please adjust your request")
	}

	return nil
}

func deleteHost(host *datastructs.Host) (affected int64, err error) {
	affected, err = DB.DeleteHost(host)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}

var deleteGroupVar = &cobra.Command{
	Use:               "group 'group name'",
	Short:             "delete existing group",
	Example:           "admiral delete group group1",
	ValidArgsFunction: groupsArgsFunc,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteGroupCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func deleteGroupCase(args []string) error {
	var groups []datastructs.Group

	var err error

	groups, err = scanGroups(args[0])
	if err != nil {
		return err
	}

	switch len(groups) {
	case 0:
		return fmt.Errorf("no group matched request")
	case 1:
		printGroups(groups)

		if User.confirm() {
			affected, err := deleteGroup(&groups[0])
			if err != nil {
				return err
			}

			fmt.Printf("lines deleted %v\n", affected)
		} else {
			return fmt.Errorf("aborted")
		}
	default:
		printGroups(groups)
		return fmt.Errorf("too many results please adjust your request")
	}

	return nil
}

func deleteGroup(group *datastructs.Group) (affected int64, err error) {
	affected, err = DB.DeleteGroup(group)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}

var deleteChildVar = &cobra.Command{
	Use:               "child 'child group' 'parent group'",
	Short:             "delete existing child-group relationship",
	Long:              "delete existing child-group relationship expecting ordered arguments child and parent group names",
	Example:           "admiral delete child child-group parent-group",
	ValidArgsFunction: groupsArgsFunc,
	Args:              cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteChildCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func deleteChildCase(args []string) error {
	var childGroups []datastructs.ChildGroup

	var err error

	childGroups, err = viewChildGroup(args[0], args[1])
	if err != nil {
		return (err)
	}

	printChildGroups(childGroups)

	if User.confirm() {
		affected, err := deleteChildGroup(&childGroups[0])
		if err != nil {
			return err
		}

		fmt.Printf("lines deleted %v\n", affected)
	} else {
		return fmt.Errorf("aborted")
	}

	return nil
}

func deleteChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error) {
	affected, err = DB.DeleteChildGroup(childGroup)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}
