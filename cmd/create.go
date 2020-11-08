package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/via-justa/admiral/datastructs"
)

func init() {
	rootCmd.AddCommand(create)

	create.AddCommand(createHostVar)
	create.AddCommand(createGroupVar)
	create.AddCommand(createChildVar)
}

var create = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "edit"},
	Short:   "create or modify existing record",
}

var createHostVar = &cobra.Command{
	Use:   "host",
	Short: "create or modify host",
	Long: "create new host or modify existing one, expecting argument host hostname/fqdn as the host to create or edit" +
		"the new or edited host would open in your favorite editor as editable json",
	Example: "admiral create host new-host\nadmiral create host new-host.domain.com\nadmiral edit host existing-host",
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []datastructs.Host
		var host datastructs.Host
		var err error

		switch len(args) {
		case 0:
			log.Fatal("no host hostname argument passed")
		case 1:
			hosts, err = returnHosts(args[0])
			if err != nil {
				log.Print(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		host, err = editHost(&hosts[0], args[0])
		if err != nil {
			log.Print(err)
		}

		printHosts([]datastructs.Host{host})

		if confirm() {
			err = confirmedHost(&host)
			if err != nil {
				log.Print(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}

func returnHosts(val string) (hosts []datastructs.Host, err error) {
	checkedVal := strings.SplitN(val, ".", 2)

	var tmp datastructs.Host

	tmp, err = viewHostByHostname(checkedVal[0])

	// return empty host if hostname or fqdn (if provided) does not exists
	if (err != nil && err.Error() != "requested host does not exists") ||
		(len(checkedVal) > 1 && tmp.Domain != checkedVal[1]) {
		return []datastructs.Host{datastructs.Host{}}, err
	}

	return []datastructs.Host{tmp}, err
}

func prepHostForEdit(host *datastructs.Host, val string) (b []byte, err error) {
	switch len(host.Hostname) {
	case 0:
		checkedVal := strings.SplitN(val, ".", 2)
		tmp := datastructs.Host{}

		tmp.Hostname = checkedVal[0]
		if len(checkedVal) > 1 {
			tmp.Domain = checkedVal[1]
		}

		tmp.Variables = "{}"

		err = tmp.UnmarshalVars()
		if err != nil {
			return b, err
		}

		b, err = json.MarshalIndent(tmp, "", "  ")
		if err != nil {
			return b, err
		}
	default:
		err = host.UnmarshalVars()
		if err != nil {
			return b, err
		}

		b, err = json.MarshalIndent(host, "", "  ")
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func editHost(host *datastructs.Host, val string) (returnHost datastructs.Host, err error) {
	var hostB []byte

	hostB, err = prepHostForEdit(host, val)
	if err != nil {
		return returnHost, err
	}

	modifiedHostB, err := Edit(hostB)
	if err != nil {
		return returnHost, err
	}

	err = json.Unmarshal(modifiedHostB, &returnHost)
	if err != nil {
		return returnHost, err
	}

	err = returnHost.MarshalVars()
	if err != nil {
		log.Print(err)
	}

	return returnHost, err
}

func confirmedHost(host *datastructs.Host) (err error) {
	var group datastructs.Group

	err = createHost(host)
	if err != nil && err.Error() != "no lines affected" {
		return err
	}

	if host.DirectGroup == "" {
		return fmt.Errorf("created host without group. please make sure to add the host to default group")
	}

	group, err = viewGroupByName(host.DirectGroup)
	if err != nil {
		return err
	}

	// if host already got host-group relationship first delete it
	existingHostGroup, err := viewHostGroupByHost(host.Hostname)
	if err != nil && err.Error() != "no record matched request" {
		return err
	} else if existingHostGroup != nil {
		_, err = deleteHostGroup(&existingHostGroup[0])
		if err != nil {
			return err
		}
	}

	// retrieving the created host to get its ID
	created, err := viewHostByHostname(host.Hostname)
	if err != nil {
		return err
	}

	err = createHostGroup(&created, &group)
	if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
		return err
	}

	return err
}

func viewHostByHostname(hostname string) (host datastructs.Host, err error) {
	host, err = db.selectHost(hostname)
	if err != nil {
		return host, err
	} else if host.Hostname == "" {
		return host, fmt.Errorf("requested host does not exists")
	}

	return host, nil
}

func createHost(host *datastructs.Host) error {
	i, err := db.insertHost(host)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

func viewHostGroupByHost(host string) (hostGroup []datastructs.HostGroup, err error) {
	hostGroup, err = db.selectHostGroup(host)
	if err != nil {
		return hostGroup, err
	} else if hostGroup == nil {
		return hostGroup, fmt.Errorf("no record matched request")
	}

	return hostGroup, nil
}

func deleteHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	affected, err = db.deleteHostGroup(hostGroup)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched request")
	}

	return affected, nil
}

func createHostGroup(host *datastructs.Host, group *datastructs.Group) error {
	hostGroup := &datastructs.HostGroup{
		HostID:  host.ID,
		GroupID: group.ID,
	}

	i, err := db.insertHostGroup(hostGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

var createGroupVar = &cobra.Command{
	Use:   "group",
	Short: "create or modify group",
	Long: "create new group or modify existing one by passing argument group name" +
		"the new or edited group would open in your favorite editor as editable json",
	Example: "admiral create group new-group\nadmiral edit group existing-group",
	Run: func(cmd *cobra.Command, args []string) {
		var group datastructs.Group
		var tmpGroup datastructs.Group
		var err error

		switch len(args) {
		case 0:
			log.Fatal("no group name argument passed")
		case 1:
			tmpGroup, err = viewGroupByName(args[0])
			if err != nil {
				log.Print(err)
			}
		default:
			log.Fatal("received too many arguments")
		}

		group, err = editGroup(&tmpGroup, args[0])
		if err != nil {
			log.Fatal(err)
		}

		printGroups([]datastructs.Group{group})
		if confirm() {
			err := createGroup(&group)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}

func prepGroupForEdit(group *datastructs.Group, name string) (b []byte, err error) {
	switch len(group.Name) {
	case 0:
		tmp := datastructs.Group{}
		tmp.Name = name
		tmp.Variables = "{}"

		err = tmp.UnmarshalVars()
		if err != nil {
			return b, err
		}

		b, err = json.MarshalIndent(tmp, "", "  ")
		if err != nil {
			return b, err
		}
	default:
		err = group.UnmarshalVars()
		if err != nil {
			return b, err
		}

		b, err = json.MarshalIndent(group, "", "  ")
		if err != nil {
			return b, err
		}
	}

	return b, err
}

func editGroup(group *datastructs.Group, val string) (returnGroup datastructs.Group, err error) {
	var groupB []byte

	groupB, err = prepGroupForEdit(group, val)
	if err != nil {
		return returnGroup, err
	}

	modifiedgroupB, err := Edit(groupB)
	if err != nil {
		return returnGroup, err
	}

	err = json.Unmarshal(modifiedgroupB, &group)
	if err != nil {
		return returnGroup, err
	}

	err = group.MarshalVars()
	if err != nil {
		return returnGroup, err
	}

	return returnGroup, err
}

func createGroup(group *datastructs.Group) error {
	i, err := db.insertGroup(group)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

var createChildVar = &cobra.Command{
	Use:   "child",
	Short: "create or modify existing child-group relationship",
	Long: "create or modify existing child-group relationship expecting ordered arguments child and parent group names." +
		" If the created relationship creates relationship loop an error will be returned",
	Example: "admiral create child child-group parent-group",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var childGroups []datastructs.ChildGroup
		var err error

		// check if relationship already exists
		childGroups, _ = viewChildGroup(args[0], args[1])
		if len(childGroups) != 0 {
			log.Fatal("Group relationship already exists")
		}

		child, err := viewGroupByName(args[0])
		if err != nil {
			log.Fatal(err)
		}

		parent, err := viewGroupByName(args[1])
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
			err = createChildGroup(&parent, &child)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("aborted")
		}
	},
}

func createChildGroup(parent *datastructs.Group, child *datastructs.Group) error {
	if child.ID == parent.ID {
		return fmt.Errorf("child and parent cannot be the same group")
	}

	isLoop := isRelationshipLoop(parent, child)
	if isLoop {
		return fmt.Errorf("relationship loop detected")
	}

	childGroup := &datastructs.ChildGroup{
		ParentID: parent.ID,
		ChildID:  child.ID,
	}

	i, err := db.insertChildGroup(childGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

func isRelationshipLoop(parent, child *datastructs.Group) bool {
	children := strings.Split(child.ChildGroups, ",")

	for _, c := range children {
		if parent.Name == c {
			return true
		}
	}

	return false
}
