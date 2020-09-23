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

func init() {
	rootCmd.AddCommand(group)
	group.AddCommand(createGroup)
	group.AddCommand(viewGroup)
	group.AddCommand(deleteGroup)
	group.AddCommand(listGroup)

	createGroup.Flags().StringVarP(&name, "name", "n", "", "group name")
	createGroup.Flags().StringVarP(&variables, "variables", "v", "", "json array of variables to set on the group level")
	createGroup.Flags().BoolVarP(&enable, "enable", "e", false, "should the group be enabled")
	createGroup.Flags().BoolVarP(&monitor, "monitor", "m", false, "should the group be monitored")

	viewGroup.Flags().IntVar(&id, "id", 0, "group id")
	viewGroup.Flags().StringVarP(&name, "name", "n", "", "group name")

	deleteGroup.Flags().IntVar(&id, "id", 0, "group id")
	deleteGroup.Flags().StringVarP(&name, "name", "n", "", "group name")
}

var group = &cobra.Command{
	Use:   "group",
	Short: "Managing inventory groups",
	Args:  cobra.MinimumNArgs(1),
}

var createGroup = &cobra.Command{
	Use:   "create",
	Short: "Create new group or update existing one",
	Run:   createGroupFunc,
}

func createGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var group datastructs.Group
	if jsonPath != "" {
		groupF, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatal(err)
		}
		if len(groupF) <= 0 {
			log.Fatal("File is empty or could not be found")
		}
		err = json.Unmarshal(groupF, &group)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		group = datastructs.Group{
			Name:      name,
			Variables: variables,
			Enabled:   enable,
			Monitored: monitor,
		}
	}

	if group.Variables == "" {
		group.Variables = "{}"
	}

	if err := client.CreateGroup(group); err != nil {
		log.Fatal(err)
	}

	createdGroup, err := client.ViewGroupByName(group.Name)
	if err != nil {
		log.Fatal(err)
	}
	printGroups([]datastructs.Group{createdGroup})
}

var viewGroup = &cobra.Command{
	Use:   "view",
	Short: "view group details",
	Long:  "View group information by name (-n,--name) or group id (--id)",
	Run:   viewGroupFunc,
}

func viewGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var group datastructs.Group
	var err error

	if len(name) > 0 {
		group, err = client.ViewGroupByName(name)
		if err != nil {
			log.Fatal(err)
		}
	} else if id != 0 {
		group, err = client.ViewGroupByID(id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Missing selector flag use --help to get available options")
	}
	printGroups([]datastructs.Group{group})
}

var deleteGroup = &cobra.Command{
	Use:   "delete",
	Short: "delete group from inventory (irevertable)",
	Long:  "delete group from inventory by name (-n,--name) or group id (--id)",
	Run:   deleteGroupFunc,
}

func deleteGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var group datastructs.Group
	var err error

	if len(name) > 0 {
		group, err = client.ViewGroupByName(name)
		if err != nil {
			log.Fatal(err)
		}
	} else if id != 0 {
		group, err = client.ViewGroupByID(id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Missing selector flag use --help to get available options")
	}
	affected, err := client.DeleteGroup(group)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Lines deleted %v\n", affected)
	}
}

var listGroup = &cobra.Command{
	Use:   "list",
	Short: "list all group in the inventory",
	Run:   listGroupFunc,
}

func listGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	groups, err := client.ListGroups()
	if err != nil {
		log.Fatal(err)
	}

	printGroups(groups)
}

func printGroups(groups []datastructs.Group) {
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "name", MinWidth: 12},
		{Header: "Enabled", MinWidth: 12},
		{Header: "Monitored", MinWidth: 12},
		{Header: "Variables", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}
	tbl.Separator = " | "
	for _, group := range groups {
		err = tbl.AddRow(group.ID, group.Name, group.Enabled, group.Monitored, group.Variables)
		if err != nil {
			log.Fatal(err)
		}
	}

	tbl.Print()
}
