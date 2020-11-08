package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tatsushid/go-prettytable"
	"github.com/via-justa/admiral/datastructs"
)

const separator = " | "

func printHosts(hosts []datastructs.Host) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "IP", MinWidth: 12},
		{Header: "Hostname", MinWidth: 12},
		{Header: "domain", MinWidth: 12},
		{Header: "Enabled", MinWidth: 12},
		{Header: "Monitored", MinWidth: 12},
		{Header: "Direct Groups", MinWidth: 12},
		{Header: "Inherited Groups", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = separator

	for i := range hosts {
		err = tbl.AddRow(hosts[i].ID, hosts[i].Host, hosts[i].Hostname, hosts[i].Domain, hosts[i].Enabled,
			hosts[i].Monitored, hosts[i].DirectGroup, hosts[i].InheritedGroups)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func printGroups(groups []datastructs.Group) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "name", MinWidth: 12},
		{Header: "Enabled", MinWidth: 12},
		{Header: "Monitored", MinWidth: 12},
		{Header: "Children count", MinWidth: 12},
		{Header: "Hosts count", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = separator

	for _, group := range groups {
		err = tbl.AddRow(group.ID, group.Name, group.Enabled, group.Monitored, group.NumChildren, group.NumHosts)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func printChildGroups(childGroups []datastructs.ChildGroup) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "Parent", MinWidth: 12},
		{Header: "Parent ID", MinWidth: 12},
		{Header: "Child", MinWidth: 12},
		{Header: "Child ID", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	// nolint: goconst
	tbl.Separator = separator

	for _, childGroup := range childGroups {
		err = tbl.AddRow(childGroup.ID, childGroup.Parent, childGroup.ParentID, childGroup.Child, childGroup.ChildID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func printHostGroups(hostGroups []datastructs.HostGroup) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "Group", MinWidth: 12},
		{Header: "Group ID", MinWidth: 12},
		{Header: "Hostname", MinWidth: 12},
		{Header: "Host ID", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}

	tbl.Separator = " | "

	for _, hostGroup := range hostGroups {
		err = tbl.AddRow(hostGroup.Group, hostGroup.GroupID, hostGroup.Host, hostGroup.HostID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// nolint: errcheck,gosec
	tbl.Print()
}

func confirm() bool {
	r := bufio.NewReader(os.Stdin)

	for tries := 2; tries > 0; tries-- {
		fmt.Print("Please confirm [y/n]: ")

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}
