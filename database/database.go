package database

import (
	"github.com/via-justa/admiral/config"
	"github.com/via-justa/admiral/database/internal/mariadb"
	"github.com/via-justa/admiral/database/internal/sqlite"
	"github.com/via-justa/admiral/datastructs"
)

// DBInterface expose backend database functions
type DBInterface interface {
	// hosts
	SelectHost(hostname string) (returnedHost datastructs.Host, err error)
	GetHosts() (hosts []datastructs.Host, err error)
	InsertHost(host *datastructs.Host) (affected int64, err error)
	DeleteHost(host *datastructs.Host) (affected int64, err error)
	ScanHosts(val string) (hosts []datastructs.Host, err error)
	// groups
	SelectGroup(name string) (returnedGroup datastructs.Group, err error)
	GetGroups() (groups []datastructs.Group, err error)
	InsertGroup(group *datastructs.Group) (affected int64, err error)
	DeleteGroup(group *datastructs.Group) (affected int64, err error)
	ScanGroups(val string) (groups []datastructs.Group, err error)
	// childGroups
	SelectChildGroup(child, parent string) (childGroups []datastructs.ChildGroup, err error)
	GetChildGroups() (childGroups []datastructs.ChildGroup, err error)
	InsertChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error)
	DeleteChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error)
	ScanChildGroups(val string) (childGroups []datastructs.ChildGroup, err error)
	// HOstGroups
	SelectHostGroup(host string) (hostGroups []datastructs.HostGroup, err error)
	GetHostGroups() (hostGroups []datastructs.HostGroup, err error)
	InsertHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error)
	DeleteHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error)
	ScanHostGroups(val string) (hostGroups []datastructs.HostGroup, err error)
	// Demo Data
	PopulateTestData(fixturesPath string) (err error)
	Close() (err error)
}

// Connect return database connection from config
func Connect(conf *config.Config) (db DBInterface, err error) {
	switch {
	case conf.MariaDB.Host != "":
		return mariadb.Connect(conf.MariaDB)
	case conf.SQLite.Path != "":
		return sqlite.Connect(&conf.SQLite)
	}

	return db, err
}
