package cli

import (
	"fmt"

	"github.com/via-justa/admiral/datastructs"
)

type dbMock struct{}

// Simulate the following DB records

// host:
// id|host   |hostname |domain|variables       |enabled|monitored|
// --|-------|---------|------|----------------|-------|---------|
//  1|1.1.1.1|host1    |local |{"var1": "val1"}|      1|        1|
//  2|2.2.2.2|host2    |local |{"var2": "val2"}|      1|        0|

// group
// id|name  |variables       |enabled|monitored|
// --|------|----------------|-------|---------|
//  1|group1|{"var1": "val1"}|      1|        1|
//  2|group2|{"var2": "val2"}|      1|        1|
//  3|group3|{"var3": "val3"}|      1|        1|

// hostgroups
// id|host_id|group_id|
// --|-------|--------|
//  1|      1|       1|

// childgroups
// id|child_id|parent_id|
// --|--------|---------|
//  1|       1|        2|
//  2|       2|        3|

// Group

func (d dbMock) selectGroup(name string, id int) (returnedGroup datastructs.Group, err error) {
	switch {
	// Existing record group1
	case name == "group1" || id == 1:
		return datastructs.Group{ID: 1, Name: "group1", Variables: "{\"var1\": \"val1\"}", Enabled: true, Monitored: true}, nil
	// Existing record group2
	case name == "group2" || id == 2:
		return datastructs.Group{ID: 2, Name: "group2", Variables: "{\"var2\": \"val2\"}", Enabled: true, Monitored: true}, nil
	// Existing record group3
	case name == "group3" || id == 3:
		return datastructs.Group{ID: 3, Name: "group3", Variables: "{\"var3\": \"val3\"}", Enabled: true, Monitored: true}, nil
	// Error missing param
	case name == "" && id == 0:
		return returnedGroup, fmt.Errorf("missing param")
	// return empty if group does not exists
	default:
		return returnedGroup, nil
	}
}

func (d dbMock) getGroups() (groups []datastructs.Group, err error) {
	return []datastructs.Group{datastructs.Group{ID: 1, Name: "group1", Variables: "{\"var1\": \"val1\"}", Enabled: true, Monitored: true}, datastructs.Group{ID: 2, Name: "group2", Variables: "{\"var2\": \"val2\"}", Enabled: true, Monitored: true}, datastructs.Group{ID: 3, Name: "group3", Variables: "{\"var3\": \"val3\"}", Enabled: true, Monitored: true}}, nil
}

func (d dbMock) insertGroup(group datastructs.Group) (affected int64, err error) {
	switch {
	// insert existing group without changes
	case group == datastructs.Group{Name: "group1", Variables: "{\"var1\": \"val1\"}", Enabled: true, Monitored: true}:
		return 0, nil
	// Update existing group
	case group == datastructs.Group{Name: "group1", Variables: "{\"var1\": \"val1\", \"var2\": \"val2\"}", Enabled: true, Monitored: true}:
		return 1, nil
	// Insert new group
	case group.Name != "" && group.Name != "group1" && group.Name != "group2" && group.Name != "group3":
		return 1, nil
	// We can get here in case the group name is empty
	default:
		return 0, fmt.Errorf("query error")
	}
}

func (d dbMock) deleteGroup(group datastructs.Group) (affected int64, err error) {
	switch {
	case group.ID == 1 || group.ID == 2 || group.ID == 3:
		return 1, nil
	default:
		return 0, nil
	}
}

// Host

func (d dbMock) selectHost(hostname string, ip string, id int) (returnedHost datastructs.Host, err error) {
	switch {
	// Existing record
	case hostname == "host1" || ip == "1.1.1.1" || id == 1:
		return datastructs.Host{ID: 1, Hostname: "host1", Host: "1.1.1.1", Domain: "local", Variables: "{\"var1\": \"val1\"}", Enabled: true, Monitored: true, Groups: []string{}}, nil
	// Error missing param
	case hostname == "" && ip == "" && id == 0:
		return returnedHost, fmt.Errorf("missing param")
	// return empty if host does not exists
	default:
		return returnedHost, nil
	}
}

func (d dbMock) getHosts() (hosts []datastructs.Host, err error) {
	return []datastructs.Host{datastructs.Host{ID: 1, Hostname: "host1", Host: "1.1.1.1", Domain: "local", Variables: "{\"var1\": \"val1\"}", Enabled: true, Monitored: true, Groups: []string{}}, datastructs.Host{ID: 2, Hostname: "host2", Host: "2.2.2.2", Domain: "local", Variables: "{\"var2\": \"val2\"}", Enabled: true, Monitored: false, Groups: []string{}}}, nil
}

func (d dbMock) insertHost(host datastructs.Host) (affected int64, err error) {
	switch {
	// insert existing host without changes
	case host.ID == 1 && host.Hostname == "host1" && host.Host == "1.1.1.1" && host.Variables == "{\"var1\": \"val1\"}" && host.Enabled == true && host.Monitored == true:
		return affected, nil
	// Update existing host (Monitored)
	case host.ID == 1 && host.Hostname == "host1" && host.Host == "1.1.1.1" && host.Variables == "{\"var1\": \"val1\"}" && host.Enabled == true && host.Monitored == false:
		return 1, nil
	// Insert new host
	case host.Hostname != "" && host.Hostname != "host1" && host.Hostname != "host2" && host.Host != "" && host.Host != "1.1.1.1" && host.Host != "2.2.2.2":
		return 1, nil
	// We can get here in case the host HostName or Host is empty
	default:
		return affected, fmt.Errorf("query error")
	}
}

func (d dbMock) deleteHost(host datastructs.Host) (affected int64, err error) {
	switch {
	case host.ID == 1 || host.ID == 2:
		return 1, nil
	default:
		return 0, nil
	}
}

// Host-group

func (d dbMock) insertHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	switch {
	// Duplicate record
	case hostGroup.Host == 1 && hostGroup.Group == 1:
		return 0, nil
	// Simulate forign key missing
	case hostGroup.Host >= 3 || hostGroup.Group >= 3:
		return 0, fmt.Errorf("Group does not exist")
	default:
		return 1, nil
	}
}

func (d dbMock) selectHostGroup(host, group int) (hostGroups []datastructs.HostGroup, err error) {
	switch {
	// Existing record
	case host == 1 || group == 1:
		return []datastructs.HostGroup{datastructs.HostGroup{ID: 1, Host: 1, Group: 1}}, nil
	// The rest of the request should return empty
	default:
		return hostGroups, nil
	}
}

func (d dbMock) getHostGroups() (hostGroups []datastructs.HostGroup, err error) {
	return []datastructs.HostGroup{datastructs.HostGroup{ID: 1, Host: 1, Group: 1}}, nil
}

func (d dbMock) deleteHostGroup(hostGroup datastructs.HostGroup) (affected int64, err error) {
	switch {
	case hostGroup.Host == 1 || hostGroup.Group == 1:
		return 1, nil
	default:
		return 0, nil
	}
}

// Child-group

func (d dbMock) insertChildGroup(childGroup datastructs.ChildGroup) (affected int64, err error) {
	switch {
	// Duplicate record
	case childGroup.Child == 1 && childGroup.Parent == 2:
		return 0, nil
	// Simulate forign key missing
	case childGroup.Child > 3 || childGroup.Parent > 3:
		return 0, fmt.Errorf("Group does not exist")
	default:
		return 1, nil
	}
}

func (d dbMock) selectChildGroup(child, parent int) (childGroups []datastructs.ChildGroup, err error) {
	switch {
	// Existing record 1 - group 1 is child of group 2
	case child == 1 || parent == 2:
		return []datastructs.ChildGroup{datastructs.ChildGroup{ID: 1, Child: 1, Parent: 2}}, nil
	// Existing record 2 - group 2 is child of group 3
	case child == 2 || parent == 3:
		return []datastructs.ChildGroup{datastructs.ChildGroup{ID: 2, Child: 2, Parent: 3}}, nil
	// The rest of the request should return empty
	default:
		return childGroups, nil
	}
}

func (d dbMock) getChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	return []datastructs.ChildGroup{datastructs.ChildGroup{ID: 1, Child: 1, Parent: 2}, datastructs.ChildGroup{ID: 2, Child: 2, Parent: 3}}, nil
}

func (d dbMock) deleteChildGroup(ChildGroup datastructs.ChildGroup) (affected int64, err error) {
	switch {
	case ChildGroup.Child == 1 || ChildGroup.Parent == 2:
		return 1, nil
	default:
		return 0, nil
	}
}
