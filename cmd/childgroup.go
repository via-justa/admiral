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
	parentName string
	childName  string
)

func init() {
	rootCmd.AddCommand(childGroup)
	childGroup.AddCommand(createChildGroup)
	childGroup.AddCommand(viewChildGroup)
	childGroup.AddCommand(deleteChildGroup)
	childGroup.AddCommand(listChildGroup)

	childGroup.PersistentFlags().StringVarP(&parentName, "parent", "p", "", "parent group name")
	childGroup.PersistentFlags().StringVarP(&childName, "child", "c", "", "child group name")

}

type ChildGroupByName struct {
	Parent string
	Child  string
}

var childGroup = &cobra.Command{
	Use:   "child-group",
	Short: "Managing groups relationship",
	Args:  cobra.MinimumNArgs(1),
}

var createChildGroup = &cobra.Command{
	Use:   "create",
	Short: "Create new child-group relationship or update existing one",
	Run:   createChildGroupFunc,
}

func createChildGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var childGroupByName ChildGroupByName
	if jsonPath != "" {
		childGroupByNameF, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatal(err)
		}
		if len(childGroupByNameF) <= 0 {
			log.Fatal("File is empty or could not be found")
		}
		err = json.Unmarshal(childGroupByNameF, &childGroupByName)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		childGroupByName = ChildGroupByName{
			Parent: parentName,
			Child:  childName,
		}
	}

	parentGroup, err := client.ViewGroupByName(childGroupByName.Parent)
	if err != nil {
		log.Fatal(err)
	}

	childGroup, err := client.ViewGroupByName(childGroupByName.Child)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.CreateChildGroup(parentGroup, childGroup); err != nil {
		log.Fatal(err)
	}

	createdChildGroup, err := client.ViewChildGroupsByChild(childGroup.ID)
	if err != nil {
		log.Fatal(err)
	}
	printChildGroups(createdChildGroup)
}

var viewChildGroup = &cobra.Command{
	Use:   "view",
	Short: "view child-group relationship",
	Long:  "View child-group relationship by parent (-p, --parent) or child (-c, --child)",
	Run:   viewChildGroupFunc,
}

func viewChildGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var childGroups []datastructs.ChildGroup

	if len(parentName) > 0 {
		group, err := client.ViewGroupByName(parentName)
		if err != nil {
			log.Fatal(err)
		}

		childGroups, err = client.ViewChildGroupsByParent(group.ID)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(childName) > 0 {
		group, err := client.ViewGroupByName(childName)
		if err != nil {
			log.Fatal(err)
		}

		childGroups, err = client.ViewChildGroupsByChild(group.ID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Missing selector flag use --help to get available options")
	}
	printChildGroups(childGroups)
}

var deleteChildGroup = &cobra.Command{
	Use:   "delete",
	Short: "delete child-group relationship from inventory (irevertable)",
	Long:  "delete child-group relationship from inventory require parent (-p, --parent) and child (-c, --child)",
	Run:   deleteChildGroupFunc,
}

func deleteChildGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()
	var childGroup datastructs.ChildGroup
	var err error

	if len(parentName) == 0 || len(childName) == 0 {
		log.Fatal("Missing selector flag use --help to get available options")
	}

	pGroup, err := client.ViewGroupByName(parentName)
	if err != nil {
		log.Fatal(err)
	}

	cGroup, err := client.ViewGroupByName(childName)
	if err != nil {
		log.Fatal(err)
	}

	childGroup = datastructs.ChildGroup{
		Parent: pGroup.ID,
		Child:  cGroup.ID,
	}

	affected, err := client.DeleteChildGroup(childGroup)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Lines deleted %v\n", affected)
	}
}

var listChildGroup = &cobra.Command{
	Use:   "list",
	Short: "print all child-group relationship on the inventory",
	Run:   listChildGroupFunc,
}

func listChildGroupFunc(cmd *cobra.Command, args []string) {
	client := cli.NewConfig()

	childGroups, err := client.ListChildGroups()
	if err != nil {
		log.Fatal(err)
	}

	printChildGroups(childGroups)
}

func printChildGroups(childGroups []datastructs.ChildGroup) {
	client := cli.NewConfig()
	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID"},
		{Header: "Parent", MinWidth: 12},
		{Header: "Child", MinWidth: 12},
	}...)
	if err != nil {
		log.Fatal(err)
	}
	tbl.Separator = " | "
	for _, childGroup := range childGroups {
		pGroup, _ := client.ViewGroupByID(childGroup.Parent)
		cGroup, _ := client.ViewGroupByID(childGroup.Child)

		err = tbl.AddRow(childGroup.ID, pGroup.Name, cGroup.Name)
		if err != nil {
			log.Fatal(err)
		}
	}

	tbl.Print()
}
