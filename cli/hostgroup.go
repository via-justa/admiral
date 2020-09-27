package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

func CreateHostGroup(host datastructs.Host, group datastructs.Group) error {
	hostGroup := datastructs.HostGroup{
		Host:  host.ID,
		Group: group.ID,
	}

	i, err := db.insertHostGroup(hostGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func ViewHostGroupByHost(hostID int) (hostGroup []datastructs.HostGroup, err error) {
	hostGroup, err = db.selectHostGroup(hostID, 0)
	if err != nil {
		return hostGroup, err
	} else if hostGroup == nil {
		return hostGroup, fmt.Errorf("No record matched request")
	}

	return hostGroup, nil
}

func ViewHostGroupByGroup(groupID int) (hostGroup []datastructs.HostGroup, err error) {
	hostGroup, err = db.selectHostGroup(0, groupID)
	if err != nil {
		return hostGroup, err
	} else if hostGroup == nil {
		return hostGroup, fmt.Errorf("No record matched request")
	}

	return hostGroup, nil
}

func ListHostGroup() (hostGroups []datastructs.HostGroup, err error) {
	hostGroups, err = db.getHostGroups()
	if err != nil {
		return hostGroups, err
	}

	return hostGroups, nil
}

func DeleteHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	affected, err = db.deleteHostGroup(hostGroup)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("No record matched")
	}

	return affected, nil
}
