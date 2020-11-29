package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/via-justa/admiral/datastructs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.AddCommand(importHosts)
	importCmd.AddCommand(importGroups)
	importCmd.AddCommand(importChildren)
}

var importCmd = &cobra.Command{
	Use:        "import {hosts | groups | children} [file path]",
	ValidArgs:  []string{"host", "group", "children"},
	ArgAliases: []string{"hosts", "groups"},
	Short:      "bulk import hosts groups and child group relationships",
	Long:       "bulk import hosts, groups or child group relationships from json encoded file",
}

var importHosts = &cobra.Command{
	Use:   "hosts [file path]",
	Short: "bulk import hosts",
	Long:  "bulk import hosts from json encoded file [file path]",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := importHostsFromPath(args); err != nil {
			log.Fatal(err)
		}
	},
}

func importHostsFromPath(args []string) (err error) {
	var file []byte

	var hosts []datastructs.Host

	file, err = ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &hosts)
	if err != nil {
		return err
	}

	for i := range hosts {
		err = hosts[i].MarshalVars()
		if err != nil {
			return err
		}

		err = createHost(&hosts[i])
		if err != nil {
			return err
		}
	}

	return err
}

var importGroups = &cobra.Command{
	Use:   "groups [file path]",
	Short: "bulk import groups",
	Long:  "bulk import groups from json encoded file [file path]",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := importGroupsFromPath(args); err != nil {
			log.Fatal(err)
		}
	},
}

func importGroupsFromPath(args []string) (err error) {
	var file []byte

	var groups []datastructs.Group

	file, err = ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &groups)
	if err != nil {
		return err
	}

	for i := range groups {
		err = groups[i].MarshalVars()
		if err != nil {
			return err
		}

		err = createGroup(&groups[i])
		if err != nil {
			return err
		}
	}

	return err
}

var importChildren = &cobra.Command{
	Use:   "children [file path]",
	Short: "bulk import child group",
	Long:  "bulk import child group relationships from json encoded file [file path]",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := importChildrenFromPath(args); err != nil {
			log.Fatal(err)
		}
	},
}

func importChildrenFromPath(args []string) (err error) {
	var file []byte

	var children []datastructs.ChildGroup

	file, err = ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &children)
	if err != nil {
		return err
	}

	for i := range children {
		var child, parent datastructs.Group

		child, err = viewGroupByName(children[i].Child)
		if err != nil {
			return err
		}

		parent, err = viewGroupByName(children[i].Parent)
		if err != nil {
			return err
		}

		err = createChildGroup(&child, &parent)
		if err != nil {
			return err
		}
	}

	return err
}
