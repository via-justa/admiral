package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

// CreateGroup accept group to create
func CreateGroup(group datastructs.Group) error {
	i, err := db.insertGroup(group)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("no lines affected")
	}

	return nil
}

// ViewGroupByName accept group name and return the group struct
func ViewGroupByName(name string) (group datastructs.Group, err error) {
	group, err = db.selectGroup(name, 0)
	if err != nil {
		return group, err
	} else if group.ID == 0 {
		return group, fmt.Errorf("requested group does not exists")
	}

	return group, nil
}

// ViewGroupByID accept group ID and return the group struct
func ViewGroupByID(id int) (group datastructs.Group, err error) {
	group, err = db.selectGroup("", id)
	if err != nil {
		return group, err
	} else if group.ID == 0 {
		return group, fmt.Errorf("requested group does not exists")
	}

	return group, nil
}

// ListGroups return all existing groups
func ListGroups() (groups []datastructs.Group, err error) {
	groups, err = db.getGroups()
	if err != nil {
		return groups, err
	}

	return groups, nil
}

// DeleteGroup accept group to remove
func DeleteGroup(group datastructs.Group) (affected int64, err error) {
	affected, err = db.deleteGroup(group)
	if err != nil {
		return affected, err
	} else if affected == 0 {
		return affected, fmt.Errorf("no record matched")
	}

	return affected, nil
}
