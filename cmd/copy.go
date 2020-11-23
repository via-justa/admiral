package cmd

import (
	"fmt"
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
	Use:        "copy",
	ValidArgs:  []string{"host", "group"},
	ArgAliases: []string{"hosts", "groups"},
	Aliases:    []string{"cp"},
	Short:      "create a new record from existing record",
}

var copyHostVar = &cobra.Command{
	Use:     "host {'source hostname' | 'host fqdn'} 'new hostname'",
	Aliases: []string{"hosts"},
	Short:   "create a new host from existing one",
	Long: "Use existing host record as template while creating a new host record," +
		"the new host would open in your favorite editor as editable json",
	Example: "admiral copy host existing-host " +
		"new-host\nadmiral copy host existing-host.domain.local new-host.domain.com",
	ValidArgsFunction: hostsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := copyHostCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func copyHostCase(args []string) error {
	var hosts []datastructs.Host

	var existingHost datastructs.Host

	var host datastructs.Host

	var err error

	switch len(args) {
	case 0, 1:
		return fmt.Errorf("please set source and destination host arguments")
	case 2:
		hosts, err = returnHosts(args[0])
		if err != nil {
			return err
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
		return fmt.Errorf("received too many arguments")
	}

	host, err = editHost(&existingHost, args[1])
	if err != nil {
		return err
	}

	printHosts([]datastructs.Host{host})

	if User.confirm() {
		err = confirmedHost(&host)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("aborted")
	}

	return nil
}

var copyGroupVar = &cobra.Command{
	Use:     "group  'source group' 'new group'",
	Aliases: []string{"groups"},
	Short:   "create a new group from existing one",
	Long: "Use existing group record as template while creating a new group record," +
		"the new group would open in your favorite editor as editable json",
	Example:           "admiral copy existing-group new-group",
	ValidArgsFunction: groupsArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := copyGroupCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func copyGroupCase(args []string) error {
	var templateGroup datastructs.Group

	var group datastructs.Group

	var err error

	switch len(args) {
	case 0, 1:
		return fmt.Errorf("please set source and destination group arguments")
	case 2:
		templateGroup, err = viewGroupByName(args[0])
		if err != nil {
			return err
		}

		templateGroup.Name = args[1]
	default:
		return fmt.Errorf("received too many arguments")
	}

	group, err = editGroup(&templateGroup, args[1])
	if err != nil {
		return err
	}

	printGroups([]datastructs.Group{group})

	if User.confirm() {
		err := createGroup(&group)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("aborted")
	}

	return nil
}
