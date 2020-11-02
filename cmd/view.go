package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

var viewAsJSON bool

func init() {
	rootCmd.AddCommand(view)

	view.AddCommand(viewHostVar)
	viewHostVar.Flags().BoolVarP(&viewAsJSON, "json", "j", false, "view in json format (present vars)")
	view.AddCommand(viewGroupVar)
	viewGroupVar.Flags().BoolVarP(&viewAsJSON, "json", "j", false, "view in json format (present vars)")
	view.AddCommand(viewChildVar)
}

var view = &cobra.Command{
	Use:     "view",
	Aliases: []string{"list", "ls", "get"},
	Short:   "view existing record",
}

var viewHostVar = &cobra.Command{
	Use:   "host",
	Short: "view existing host by substring of hostname or IP or view all records when no argument passed",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host

		var err error

		switch len(args) {
		case 0:
			hosts, err = listHosts()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			hosts, err = scanHosts(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		if viewAsJSON {
			for i := range hosts {
				_ = hosts[i].UnmarshalVars()
			}
			b, _ := json.MarshalIndent(hosts, "", "    ")
			fmt.Printf("%s\n", b)
		} else {
			printHosts(hosts)
		}
	},
}

func listHosts() (hosts []datastructs.Host, err error) {
	hosts, err = db.getHosts()
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

func scanHosts(val string) (hosts []datastructs.Host, err error) {
	hosts, err = db.scanHosts(val)
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

var viewGroupVar = &cobra.Command{
	Use:   "group",
	Short: "view existing group by substring of group name or view all records when no argument passed",
	Run: func(cmd *cobra.Command, args []string) {
		var groups []datastructs.Group

		var err error

		switch len(args) {
		case 0:
			groups, err = listGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			groups, err = scanGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		if viewAsJSON {
			for i := range groups {
				_ = groups[i].UnmarshalVars()
			}
			b, _ := json.MarshalIndent(groups, "", "    ")
			fmt.Printf("%s\n", b)
		} else {
			printGroups(groups)
		}
	},
}

func viewGroupByName(name string) (group datastructs.Group, err error) {
	group, err = db.selectGroup(name)
	if err != nil {
		return group, err
	} else if group.ID == 0 {
		return group, fmt.Errorf("requested group does not exists")
	}

	return group, nil
}

func listGroups() (groups []datastructs.Group, err error) {
	groups, err = db.getGroups()
	if err != nil {
		return groups, err
	}

	return groups, nil
}

func scanGroups(val string) (groups []datastructs.Group, err error) {
	groups, err = db.scanGroups(val)
	if err != nil {
		return groups, err
	}

	return groups, nil
}

var viewChildVar = &cobra.Command{
	Use:   "child",
	Short: "view existing child-group relationship by parent, child or view all records when no argument passed",
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup

		var err error

		switch len(args) {
		case 0:
			childGroups, err = listChildGroups()
			if err != nil {
				log.Fatal(err)
			}
		case 1:
			childGroups, err = scanChildGroups(args[0])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		printChildGroups(childGroups)
	},
}

func listChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = db.getChildGroups()
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func viewChildGroup(child, parent string) (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = db.selectChildGroup(child, parent)
	if err != nil {
		return childGroups, err
	} else if childGroups == nil {
		return childGroups, fmt.Errorf("no record matched request")
	}

	return childGroups, nil
}

func scanChildGroups(val string) (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = db.scanChildGroups(val)
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}
