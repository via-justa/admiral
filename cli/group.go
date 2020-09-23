package cli

import (
	"fmt"

	"github.com/via-justa/admiral/database"
	"github.com/via-justa/admiral/datastructs"
)

func (conf *Config) CreateGroup(group datastructs.Group) error {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return err
	}

	i, err := db.InsertGroup(group)
	if err != nil {
		return err
	} else if i == 0 {
		return fmt.Errorf("No lines affected")
	}

	return nil
}

func (conf *Config) ViewGroupByName(name string) (group datastructs.Group, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return group, err
	}

	group, err = db.SelectGroup(name, 0)
	if err != nil {
		return group, err
	}

	return group, nil
}

func (conf *Config) ViewGroupByID(id int) (group datastructs.Group, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return group, err
	}

	group, err = db.SelectGroup("", id)
	if err != nil {
		return group, err
	}

	return group, nil
}

func (conf *Config) ListGroups() (groups []datastructs.Group, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return groups, err
	}

	groups, err = db.GetGroups()
	if err != nil {
		return groups, err
	}

	return groups, nil
}

func (conf *Config) DeleteGroup(group datastructs.Group) (affected int64, err error) {
	db, err := database.Connect(conf.Database)
	if err != nil {
		return affected, err
	}

	affected, err = db.DeleteGroup(group)
	if err != nil {
		return affected, err
	}

	return affected, nil
}
