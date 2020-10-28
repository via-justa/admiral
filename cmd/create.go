package cmd

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/cli"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(create)

	create.AddCommand(createHost)
	create.AddCommand(createGroup)
	create.AddCommand(createChild)
}

var create = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "edit"},
	Short:   "create or modify existing record",
}

var createHost = &cobra.Command{
	Use:   "host",
	Short: "create or modify host. expecting one argument host hostname",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host
		var host datastructs.Host
		var group datastructs.Group
		var hostB []byte
		var err error

		switch len(args) {
		case 0:
			log.Fatal("no host hostname argument passed")
		case 1:
			var tmp datastructs.Host
			tmp, err = cli.ViewHostByHostname(args[0])
			if err != nil && err.Error() != "requested host does not exists" {
				log.Print(err)
			}
			hosts = []datastructs.Host{tmp}
		default:
			log.Fatal("received too many arguments")
		}

		switch len(hosts[0].Hostname) {
		case 0:
			tmp := datastructs.Host{}
			tmp.Hostname = args[0]
			tmp.Variables = "{}"
			err = tmp.UnmarshalVars()
			if err != nil {
				log.Print(err)
			}

			hostB, err = json.MarshalIndent(tmp, "", "  ")
			if err != nil {
				log.Print(err)
			}
		default:
			err = hosts[0].UnmarshalVars()
			if err != nil {
				log.Print(err)
			}

			hostB, err = json.MarshalIndent(hosts[0], "", "  ")
			if err != nil {
				log.Print(err)
			}
		}

		modifiedHostB, err := Edit(hostB)
		if err != nil {
			log.Print(err)
		}

		err = json.Unmarshal(modifiedHostB, &host)
		if err != nil {
			log.Print(err)
		}

		err = host.MarshalVars()
		if err != nil {
			log.Print(err)
		}

		printHosts([]datastructs.Host{host})
		if confirm() {
			err = cli.CreateHost(&host)
			if err != nil && err.Error() != "no lines affected" {
				log.Fatal(err)
			}

			group, err = cli.ViewGroupByName(host.DirectGroup)
			if err != nil {
				log.Fatal(err)
			}

			// if host already got host-group relationship first delete it
			existingHostGroup, err := cli.ViewHostGroupByHost(host.Hostname)
			if err != nil && err.Error() != "no record matched request" {
				log.Fatal(err)
			} else if existingHostGroup != nil {
				_, err = cli.DeleteHostGroup(&existingHostGroup[0])
				if err != nil {
					log.Fatal(err)
				}
			}

			// retrieving the created host to get its ID
			created, err := cli.ViewHostByHostname(args[0])
			if err != nil {
				log.Print(err)
			}

			err = cli.CreateHostGroup(&created, &group)
			if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
				log.Fatal(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}

var createGroup = &cobra.Command{
	Use:   "group",
	Short: "create or modify group. expecting one argument group name",
	Run: func(cmd *cobra.Command, args []string) {
		var groups []datastructs.Group
		var group datastructs.Group
		var groupB []byte
		var err error

		switch len(args) {
		case 0:
			log.Fatal("no group name argument passed")
		case 1:
			var tmp datastructs.Group
			tmp, err = cli.ViewGroupByName(args[0])
			if err != nil {
				log.Print(err)
			}
			groups = []datastructs.Group{tmp}
		default:
			log.Fatal("received too many arguments")
		}

		switch len(groups[0].Name) {
		case 0:
			tmp := datastructs.Group{}
			tmp.Name = args[0]
			tmp.Variables = "{}"
			err = tmp.UnmarshalVars()
			if err != nil {
				log.Print(err)
			}

			groupB, err = json.MarshalIndent(tmp, "", "  ")
			if err != nil {
				log.Print(err)
			}
		default:
			err = groups[0].UnmarshalVars()
			if err != nil {
				log.Print(err)
			}

			groupB, err = json.MarshalIndent(groups[0], "", "  ")
			if err != nil {
				log.Print(err)
			}
		}

		modifiedgroupB, err := Edit(groupB)
		if err != nil {
			log.Print(err)
		}

		err = json.Unmarshal(modifiedgroupB, &group)
		if err != nil {
			log.Print(err)
		}

		err = group.MarshalVars()
		if err != nil {
			log.Print(err)
		}

		printGroups([]datastructs.Group{group})
		if confirm() {
			err := cli.CreateGroup(&group)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}

var createChild = &cobra.Command{
	Use:   "child",
	Short: "create or modify existing child-group relationship",
	Long:  "create or modify existing child-group relationship expecting ordered arguments child and parent group names",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup
		var err error

		// check if relationship already exists
		childGroups, _ = cli.ViewChildGroup(args[0], args[1])
		if len(childGroups) != 0 {
			log.Fatal("Group relationship already exists")
		}

		child, err := cli.ViewGroupByName(args[0])
		if err != nil {
			log.Fatal(err)
		}

		parent, err := cli.ViewGroupByName(args[1])
		if err != nil {
			log.Fatal(err)
		}

		childGroups = []datastructs.ChildGroup{datastructs.ChildGroup{
			Parent:   parent.Name,
			ParentID: parent.ID,
			Child:    child.Name,
			ChildID:  child.ID,
		}}

		printChildGroups(childGroups)
		if confirm() {
			err = cli.CreateChildGroup(&parent, &child)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}
