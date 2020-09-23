package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

func CreateChildGroup(parent datastructs.Group, child datastructs.Group) error {
	childGroup := datastructs.ChildGroup{
		Parent: parent.ID,
		Child:  child.ID,
	}

	// TODO: check child and parent not the same
	// TODO: check parent is not child of child (relationship loop)
	i, err := db.insertChildGroup(childGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func ViewChildGroupsByParent(parentID int) (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = db.selectChildGroup(0, parentID)
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func ViewChildGroupsByChild(childID int) (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = db.selectChildGroup(childID, 0)
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func ListChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	childGroups, err = db.getChildGroups()
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func DeleteChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	affected, err = db.deleteChildGroup(childGroup)
	if err != nil {
		return affected, err
	}

	return affected, nil
}
