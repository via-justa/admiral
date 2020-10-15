package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

// CreateHostGroup accept host and group and create new relationship
func CreateHostGroup(host *datastructs.Host, group *datastructs.Group) error {
	hostGroup := &datastructs.HostGroup{
		Host:  host.ID,
		Group: group.ID,
	}

	i, err := db.insertHostGroup(hostGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

// ViewHostGroupByHost accept host ID and return all host-group relationships for the host
func ViewHostGroupByHost(host string) (hostGroup []datastructs.HostGroupView, err error) {
	hostGroup, err = db.selectHostGroup(host, "")
	if err != nil {
		return hostGroup, err
	} else if hostGroup == nil {
		return hostGroup, fmt.Errorf("no record matched request")
	}

	return hostGroup, nil
}

// ViewHostGroupByGroup accept group ID and return all host-group relationships for the group
func ViewHostGroupByGroup(group string) (hostGroup []datastructs.HostGroupView, err error) {
	hostGroup, err = db.selectHostGroup("", group)
	if err != nil {
		return hostGroup, err
	} else if hostGroup == nil {
		return hostGroup, fmt.Errorf("no record matched request")
	}

	return hostGroup, nil
}

// ListHostGroup list all existing host-group relationships
func ListHostGroup() (hostGroups []datastructs.HostGroupView, err error) {
	hostGroups, err = db.getHostGroups()
	if err != nil {
		return hostGroups, err
	}

	return hostGroups, nil
}

// DeleteHostGroup accept hostGroup to remove
func DeleteHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	affected, err = db.deleteHostGroup(hostGroup)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}
