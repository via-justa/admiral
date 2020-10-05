package database

import (
	"github.com/via-justa/admiral/database/internal/db"
	"github.com/via-justa/admiral/datastructs"
)

// Hosts

// SelectHost return host information. The function will search for the host in the following order:
// By hostname, if hostname is empty by host and if both hostname and host are empty by id
func SelectHost(hostname string, ip string, id int) (returnedHost datastructs.Host, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.SelectHost(hostname, ip, id)
}

// GetHosts return all hosts in the inventory
func GetHosts() (hosts []datastructs.Host, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.GetHosts()
}

// InsertHost accept Host to insert or update and return the number of affected rows and error if exists
func InsertHost(host *datastructs.Host) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.InsertHost(host)
}

// DeleteHost accept Host to delete and return the number of affected rows and error if exists
func DeleteHost(host *datastructs.Host) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.DeleteHost(host)
}

// Groups

// SelectGroup return group information. The function will search for the group in the following order:
// By name, if name is empty by id
func SelectGroup(name string, id int) (returnedGroup datastructs.Group, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.SelectGroup(name, id)
}

// GetGroups return all groups in the inventory
func GetGroups() (groups []datastructs.Group, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.GetGroups()
}

// InsertGroup accept Group to insert or update and return the number of affected rows and error if exists
func InsertGroup(group datastructs.Group) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.InsertGroup(group)
}

// DeleteGroup accept Group to delete and return the number of affected rows and error if exists
func DeleteGroup(group datastructs.Group) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.DeleteGroup(group)
}

// ChildGroups

// SelectChildGroup accept either child or parent id and return slice of ids for parent or child groups respectively.
// If child is provided will return slice of parent ids
// If parent is provided will return slice of child ids
// will error if none is provided
func SelectChildGroup(child, parent string) (childGroups []datastructs.ChildGroupView, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.SelectChildGroup(child, parent)
}

// GetChildGroups return all child groups relationships in the inventory
func GetChildGroups() (childGroups []datastructs.ChildGroupView, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.GetChildGroups()
}

// InsertChildGroup accept ChildGroup to insert and return the number of affected rows and error if exists
func InsertChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.InsertChildGroup(childGroup)
}

// DeleteChildGroup accept ChildGroup to delete and return the number of affected rows and error if exists
func DeleteChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.DeleteChildGroup(childGroup)
}

// HostGroups

// SelectHostGroup accept either host or group id and return slice of ids for groups or hosts respectively.
// If host is provided will return slice of groups ids
// If group is provided will return slice of hosts ids
// will error if none is provided
func SelectHostGroup(host, group string) (hostGroups []datastructs.HostGroupView, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.SelectHostGroup(host, group)
}

// GetHostGroups return all host groups relationships in the inventory
func GetHostGroups() (hostGroups []datastructs.HostGroupView, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.GetHostGroups()
}

// InsertHostGroup accept HostGroup to insert and return the number of affected rows and error if exists
func InsertHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.InsertHostGroup(hostGroup)
}

// DeleteHostGroup accept HostGroup to delete and return the number of affected rows and error if exists
func DeleteHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	conf := db.NewConfig()
	conn, _ := db.Connect(conf.Database)

	return conn.DeleteHostGroup(hostGroup)
}
