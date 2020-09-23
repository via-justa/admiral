package cli

import (
	"fmt"

	"github.com/via-justa/admiral/database"
	"github.com/via-justa/admiral/datastructs"
)

func (conf *Config) CreateChildGroup(parent datastructs.Group, child datastructs.Group) error {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return err
	}

	childGroup := datastructs.ChildGroup{
		Parent: parent.ID,
		Child:  child.ID,
	}

	i, err := db.InsertChildGroup(childGroup)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func (conf *Config) ViewChildGroupsByParent(parentID int) (childGroups []datastructs.ChildGroup, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return childGroups, err
	}

	childGroups, err = db.SelectChildGroup(0, parentID)
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func (conf *Config) ViewChildGroupsByChild(childID int) (childGroups []datastructs.ChildGroup, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return childGroups, err
	}

	childGroups, err = db.SelectChildGroup(childID, 0)
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func (conf *Config) ListChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return childGroups, err
	}

	childGroups, err = db.GetChildGroups()
	if err != nil {
		return childGroups, err
	}

	return childGroups, nil
}

func (conf *Config) DeleteChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return affected, err
	}

	affected, err = db.DeleteChildGroup(childGroup)
	if err != nil {
		return affected, err
	}

	return affected, nil
}
