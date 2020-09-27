package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

func CreateGroup(group datastructs.Group) error {
	i, err := db.insertGroup(group)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

func ViewGroupByName(name string) (group datastructs.Group, err error) {
	group, err = db.selectGroup(name, 0)
	if err != nil {
		return group, err
	} else if group == (datastructs.Group{}) {
		return group, fmt.Errorf("requested group does not exists")
	}

	return group, nil
}

func ViewGroupByID(id int) (group datastructs.Group, err error) {
	group, err = db.selectGroup("", id)
	if err != nil {
		return group, err
	} else if group == (datastructs.Group{}) {
		return group, fmt.Errorf("requested group does not exists")
	}

	return group, nil
}

func ListGroups() (groups []datastructs.Group, err error) {
	groups, err = db.getGroups()
	if err != nil {
		return groups, err
	}

	return groups, nil
}

func DeleteGroup(group datastructs.Group) (affected int64, err error) {
	affected, err = db.deleteGroup(group)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}
