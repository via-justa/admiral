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
	Args:              cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := copyHostCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func copyHostCase(args []string) error {
	var hosts datastructs.Hosts

	var err error

	hosts, err = scanHosts(args[0])
	if err != nil {
		return err
	}

	switch len(hosts) {
	case 0:
		return fmt.Errorf("source host does not exists")
	case 1:
		var domain string

		fqdn := strings.SplitN(args[1], ".", 2)
		if len(fqdn) > 1 {
			domain = fqdn[1]
		} else {
			domain = hosts[0].Domain
		}

		hosts = datastructs.Hosts{
			{
				Hostname:  fqdn[0],
				Domain:    domain,
				Variables: hosts[0].Variables,
				Enabled:   hosts[0].Enabled,
				Monitored: hosts[0].Monitored,
			},
		}
	default:
		return fmt.Errorf("source host matched to many records")
	}

	hosts, err = editHosts(&hosts)
	if err != nil {
		return err
	}

	printHosts(hosts)

	if User.confirm() {
		err = confirmedHosts(&hosts)
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
	Args:              cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := copyGroupCase(args); err != nil {
			log.Fatal(err)
		}
	},
}

func copyGroupCase(args []string) error {
	var group datastructs.Group

	var err error

	group, err = viewGroupByName(args[0])
	if err != nil {
		return err
	}

	group.Name = args[1]

	group, err = editGroup(&group)
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
