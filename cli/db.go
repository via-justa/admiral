package cli

import (
	"github.com/via-justa/admiral/database"
	"github.com/via-justa/admiral/datastructs"
)

type dbInterface interface {
	selectHost(hostname string, ip string, id int) (returnedHost datastructs.Host, err error)
	getHosts() (hosts []datastructs.Host, err error)
	insertHost(host datastructs.Host) (affected int64, err error)
	deleteHost(host datastructs.Host) (affected int64, err error)
	selectGroup(name string, id int) (returnedGroup datastructs.Group, err error)
	getGroups() (groups []datastructs.Group, err error)
	insertGroup(group datastructs.Group) (affected int64, err error)
	deleteGroup(group datastructs.Group) (affected int64, err error)
	selectChildGroup(child, parent int) (childGroups []datastructs.ChildGroup, err error)
	getChildGroups() (childGroups []datastructs.ChildGroup, err error)
	insertChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error)
	deleteChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error)
	selectHostGroup(host, group int) (hostGroups []datastructs.HostGroup, err error)
	getHostGroups() (hostGroups []datastructs.HostGroup, err error)
	insertHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error)
	deleteHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error)
}

type dbReal struct{}

var db dbInterface

func init() {
	db = dbReal{}
}

func (d dbReal) selectHost(hostname string, ip string, id int) (returnedHost datastructs.Host, err error) {
	return database.SelectHost(hostname, ip, id)
}

func (d dbReal) getHosts() (hosts []datastructs.Host, err error) {
	return database.GetHosts()
}

func (d dbReal) insertHost(host datastructs.Host) (affected int64, err error) {
	return database.InsertHost(host)
}

func (d dbReal) deleteHost(host datastructs.Host) (affected int64, err error) {
	return database.DeleteHost(host)
}

func (d dbReal) selectGroup(name string, id int) (returnedGroup datastructs.Group, err error) {
	return database.SelectGroup(name, id)
}

func (d dbReal) getGroups() (groups []datastructs.Group, err error) {
	return database.GetGroups()
}

func (d dbReal) insertGroup(group datastructs.Group) (affected int64, err error) {
	return database.InsertGroup(group)
}

func (d dbReal) deleteGroup(group datastructs.Group) (affected int64, err error) {
	return database.DeleteGroup(group)
}

func (d dbReal) selectChildGroup(child, parent int) (childGroups []datastructs.ChildGroup, err error) {
	return database.SelectChildGroup(child, parent)
}

func (d dbReal) getChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	return database.GetChildGroups()
}

func (d dbReal) insertChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	return database.InsertChildGroup(childGroup)
}

func (d dbReal) deleteChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	return database.DeleteChildGroup(childGroup)
}

func (d dbReal) selectHostGroup(host, group int) (hostGroups []datastructs.HostGroup, err error) {
	return database.SelectHostGroup(host, group)
}

func (d dbReal) getHostGroups() (hostGroups []datastructs.HostGroup, err error) {
	return database.GetHostGroups()
}

func (d dbReal) insertHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	return database.InsertHostGroup(hostGroup)
}

func (d dbReal) deleteHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	return database.DeleteHostGroup(hostGroup)
}
