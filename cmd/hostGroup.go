package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/tatsushid/go-prettytable"
	"github.com/via-justa/admiral/cli"
	"github.com/via-justa/admiral/datastructs"
)

var (
	groupName string
)

func init() {
	rootCmd.AddCommand(hostGroup)
	hostGroup.AddCommand(createHostGroup)
	hostGroup.AddCommand(viewHostGroup)
	hostGroup.AddCommand(deleteHostGroup)
	hostGroup.AddCommand(listHostGroup)

	hostGroup.PersistentFlags().StringVarP(&name, "hostname", "n", "", "base hostname")
	hostGroup.PersistentFlags().StringVarP(&groupName, "group", "g", "", "main group the host will be in")

}

type HostGroupByName struct {
	Host  string
	Group string
}

var hostGroup = &cobra.Command{
	Use:   "host-group",
	Short: "Managing host to groups relationship",
	Args:  cobra.MinimumNArgs(1),
}

var createHostGroup = &cobra.Command{
	Use:   "create",
	Short: "Create new host-group relationship or update existing one",
	Run:   createHostGroupFunc,
}

func createHostGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var hostGroupByName HostGroupByName
	if jsonPath != "" {
		hostGroupByNameF, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatal(err)
		}
		if len(hostGroupByNameF) <= 0 {
			log.Fatal("File is empty or could not be found")
		}
		err = json.Unmarshal(hostGroupByNameF, &hostGroupByName)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		hostGroupByName = HostGroupByName{
			Host:  name,
			Group: groupName,
		}
	}

	host, err := client.ViewHostByHostname(hostGroupByName.Host)
	if err != nil {
		log.Fatal(err)
	}

	group, err := client.ViewGroupByName(hostGroupByName.Group)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.CreateHostGroup(host, group); err != nil {
		log.Fatal(err)
	}

	createdHostGroup, err := client.ViewHostGroupByHost(host.ID)
	if err != nil {
		log.Fatal(err)
	}
	printHostGroups(createdHostGroup)
}

var viewHostGroup = &cobra.Command{
	Use:   "view",
	Short: "view host-group relationship",
	Long:  "View host-group relationship by hostname (-n, --hostname) or group (-g, --group)",
	Run:   viewHostGroupFunc,
}

func viewHostGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var hostGroup []datastructs.HostGroup

	if len(name) > 0 {
		host, err := client.ViewHostByHostname(name)
		if err != nil {
			log.Fatal(err)
		}

		hostGroup, err = client.ViewHostGroupByHost(host.ID)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(groupName) > 0 {
		group, err := client.ViewGroupByName(groupName)
		if err != nil {
			log.Fatal(err)
		}

		hostGroup, err = client.ViewHostGroupByGroup(group.ID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Missing selector flag use --help to get available options")
	}
	printHostGroups(hostGroup)
}

var deleteHostGroup = &cobra.Command{
	Use:   "delete",
	Short: "delete host-group relationship from inventory (irevertable)",
	Long:  "delete host-group relationship from inventory require hostname (-n,--hostname) and group (-g, --group)",
	Run:   deleteHostGroupFunc,
}

func deleteHostGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var hostGroup datastructs.HostGroup
	var err error

	if len(name) == 0 || len(groupName) == 0 {
		log.Fatal("Missing selector flag use --help to get available options")
	}

	host, err := client.ViewHostByHostname(name)
	if err != nil {
		log.Fatal(err)
	}

	group, err := client.ViewGroupByName(groupName)
	if err != nil {
		log.Fatal(err)
	}

	hostGroup = datastructs.HostGroup{
		Host:  host.ID,
		Group: group.ID,
	}

	affected, err := client.DeleteHostGroup(hostGroup)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Lines deleted %v\n", affected)
	}
}

var listHostGroup = &cobra.Command{
	Use:   "list",
	Short: "list all host-group relationships in the inventory",
	Run:   listHostGroupFunc,
}

func listHostGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()

	hostGroups, err := client.ListHostGroup()
	if err != nil {
		log.Fatal(err)
	}

	printHostGroups(hostGroups)
}

func printHostGroups(hostGroups []datastructs.HostGroup) {
	client := cli.NewConfig()
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "Group", MinWidth: 12},
		{Header: "Hostname", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}
	tbl.Separator = " | "
	for _, hostGroup := range hostGroups {
		group, _ := client.ViewGroupByID(hostGroup.Group)
		host, _ := client.ViewHostByID(hostGroup.Host)

		err = tbl.AddRow(hostGroup.ID, group.Name, host.Hostname)
		if err != nil {
			log.Fatal(err)
		}
	}

	tbl.Print()
}