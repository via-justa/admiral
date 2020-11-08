// nolint:
package cmd

import (
	"fmt"
	"strings"

	"github.com/via-justa/admiral/datastructs"
)

type dbMock struct{}

// Simulate the following DB records

// host:
// id|host   |hostname |domain|variables                                        |enabled|monitored|
// --|-------|---------|------|-------------------------------------------------|-------|---------|
//  1|1.1.1.1|host1    |local |{"host_var1": {"host_sub_var1": "host_sub_val1"}}|      1|        1|
//  2|2.2.2.2|host2    |local |{"var2": "val2"}                                 |      1|        1|
//  3|3.3.3.3|host3    |local |{"var3": "val3"}                                 |      1|        1|

// group
// id|name  |variables       |enabled|monitored|
// --|------|----------------|-------|---------|
//  1|group1|{"group_var1": {"group_sub_var1": "group_sub_val1"}}|      1|        1|
//  2|group2|{"var2": "val2"}|      1|        1|
//  3|group3|{"var3": "val3"}|      1|        1|
//  4|group4|{"var4": "val4"}|      1|        1|
//  5|group5|{"var5": "val5"}|      1|        1|

// hostgroups
// id|host_id|group_id|
// --|-------|--------|
//  1|      1|       1|
//  2|      2|       2|
//  3|      3|       3|

// childgroups
// id|child_id|parent_id|
// --|--------|---------|
//  1|       3|        4|
//  2|       4|        5|

var (
	testGroup1 = datastructs.Group{
		ID:        1,
		Name:      "group1",
		Variables: "{\"group_var1\": {\"group_sub_var1\": \"group_sub_val1\"}}",
		Enabled:   true,
		Monitored: true,
	}
	testGroup2 = datastructs.Group{
		ID:        2,
		Name:      "group2",
		Variables: "{\"var2\": \"val2\"}",
		Enabled:   true,
		Monitored: true,
	}
	testGroup3 = datastructs.Group{
		ID:        3,
		Name:      "group3",
		Variables: "{\"var3\": \"val3\"}",
		Enabled:   true,
		Monitored: true,
	}
	testGroup4 = datastructs.Group{
		ID:          4,
		Name:        "group4",
		Variables:   "{\"var4\": \"val4\"}",
		Enabled:     true,
		Monitored:   true,
		ChildGroups: "group3",
	}
	testGroup5 = datastructs.Group{
		ID:          5,
		Name:        "group5",
		Variables:   "{\"var5\": \"val5\"}",
		Enabled:     true,
		Monitored:   true,
		ChildGroups: "group3,group4",
	}
	testHost1 = datastructs.Host{
		ID:          1,
		Hostname:    "host1",
		Host:        "1.1.1.1",
		Domain:      "local",
		Variables:   "{\"host_var1\": {\"host_sub_var1\": \"host_sub_val1\"}}",
		Enabled:     true,
		Monitored:   true,
		DirectGroup: "group1",
	}
	testHost2 = datastructs.Host{
		ID:          2,
		Hostname:    "host2",
		Host:        "2.2.2.2",
		Domain:      "local",
		Variables:   "{\"var2\": \"val2\"}",
		Enabled:     true,
		Monitored:   true,
		DirectGroup: "group2",
	}
	testHost3 = datastructs.Host{
		ID:              3,
		Hostname:        "host3",
		Host:            "3.3.3.3",
		Domain:          "local",
		Variables:       "{\"var3\": \"val3\"}",
		Enabled:         true,
		Monitored:       true,
		DirectGroup:     "group3",
		InheritedGroups: "group4,group5",
	}
	testHostGroup1 = datastructs.HostGroup{
		ID:      1,
		HostID:  1,
		Host:    "host1",
		GroupID: 1,
		Group:   "group1",
	}
	testHostGroup2 = datastructs.HostGroup{
		ID:      2,
		HostID:  2,
		Host:    "host2",
		GroupID: 2,
		Group:   "group2",
	}
	testHostGroup3 = datastructs.HostGroup{
		ID:      3,
		HostID:  3,
		Host:    "host3",
		GroupID: 3,
		Group:   "group3",
	}
	testChild1 = datastructs.ChildGroup{
		ID:       1,
		ChildID:  3,
		Child:    "group3",
		ParentID: 4,
		Parent:   "group4",
	}
	testChild2 = datastructs.ChildGroup{
		ID:       2,
		ChildID:  4,
		Child:    "group4",
		ParentID: 5,
		Parent:   "group5",
	}
)

// Group
func (d dbMock) selectGroup(name string) (returnedGroup datastructs.Group, err error) {
	switch {
	// Existing record group1
	case name == "group1":
		return testGroup1, nil
	// Existing record group2
	case name == "group2":
		return testGroup2, nil
	// Existing record group3
	case name == "group3":
		return testGroup3, nil
	case name == "group4":
		return testGroup4, nil
	case name == "group5":
		return testGroup5, nil
	// Error missing param
	case name == "":
		return returnedGroup, fmt.Errorf("missing param")
	// return empty if group does not exists
	default:
		return returnedGroup, nil
	}
}

func (d dbMock) getGroups() (groups []datastructs.Group, err error) {
	return []datastructs.Group{testGroup1, testGroup2, testGroup3, testGroup4, testGroup5}, nil
}

func (d dbMock) insertGroup(group *datastructs.Group) (affected int64, err error) {
	switch {
	// insert existing group without changes
	case group.Name == "group1" &&
		group.Variables == "{\"group_var1\": {\"group_sub_var1\": \"group_sub_val1\"}}" &&
		group.Enabled == true &&
		group.Monitored == true:
		return 0, nil
	// Update existing group
	case group.Name == "group1" &&
		group.Variables == "{\"group_var1\": {\"group_sub_var1\": \"group_sub_val1\"}, \"var2\": \"val2\"}" &&
		group.Enabled == true &&
		group.Monitored == true:
		return 1, nil
	// Insert new group
	case group.Name != "group1" &&
		group.Name != "group2" &&
		group.Name != "group3" &&
		group.Name != "group4" &&
		group.Name != "group5" &&
		group.Name != "":
		return 1, nil
	// We can get here in case the group name is empty
	default:
		return 0, fmt.Errorf("query error")
	}
}

func (d dbMock) deleteGroup(group *datastructs.Group) (affected int64, err error) {
	switch {
	case group.ID == 1 || group.ID == 2 || group.ID == 3 || group.ID == 4 || group.ID == 5:
		return 1, nil
	default:
		return 0, nil
	}
}

func (d dbMock) scanGroups(val string) (groups []datastructs.Group, err error) {
	switch {
	// Error missing param
	case val == "":
		return groups, fmt.Errorf("missing param")
	case strings.Contains("group", val):
		return []datastructs.Group{testGroup1, testGroup2, testGroup3, testGroup4, testGroup5}, nil
	case strings.Contains("group1", val):
		return []datastructs.Group{testGroup1}, nil
	case strings.Contains("group2", val):
		return []datastructs.Group{testGroup2}, nil
	case strings.Contains("group3", val):
		return []datastructs.Group{testGroup3}, nil
	case strings.Contains("group4", val):
		return []datastructs.Group{testGroup4}, nil
	case strings.Contains("group5", val):
		return []datastructs.Group{testGroup5}, nil
	// return empty if group does not exists
	default:
		return groups, nil
	}
}

// Host

func (d dbMock) selectHost(hostname string) (returnedHost datastructs.Host, err error) {
	switch {
	// Existing record
	case hostname == "host1":
		return testHost1, nil
	case hostname == "host2":
		return testHost2, nil
	case hostname == "host3":
		return testHost3, nil
	// Error missing param
	case hostname == "":
		return returnedHost, fmt.Errorf("missing param")
	// return empty if host does not exists
	default:
		return returnedHost, nil
	}
}

func (d dbMock) getHosts() (hosts []datastructs.Host, err error) {
	return []datastructs.Host{testHost1, testHost2, testHost3}, nil
}

func (d dbMock) insertHost(host *datastructs.Host) (affected int64, err error) {
	switch {
	// insert existing host without changes
	case host.ID == 1 &&
		host.Hostname == "host1" &&
		host.Host == "1.1.1.1" &&
		host.Variables == "{\"host_var1\": {\"host_sub_var1\": \"host_sub_val1\"}}" &&
		host.Enabled == true &&
		host.Monitored == true:
		return affected, nil
	// Update existing host (Monitored)
	case host.ID == 1 &&
		host.Hostname == "host1" &&
		host.Host == "1.1.1.1" &&
		host.Variables == "{\"host_var1\": {\"host_sub_var1\": \"host_sub_val1\"}}" &&
		host.Enabled == true &&
		host.Monitored == false:
		return 1, nil
	// Insert new host
	case host.Hostname != "host1" &&
		host.Hostname != "host2" &&
		host.Hostname != "host3" &&
		host.Hostname != "" &&
		host.Host != "1.1.1.1" &&
		host.Host != "2.2.2.2" &&
		host.Host != "3.3.3.3" &&
		host.Host != "":
		return 1, nil
	// We can get here in case the host HostName or Host is empty
	default:
		return affected, fmt.Errorf("query error")
	}
}

func (d dbMock) deleteHost(host *datastructs.Host) (affected int64, err error) {
	switch {
	case host.ID == 1 || host.ID == 2 || host.ID == 3:
		return 1, nil
	default:
		return 0, nil
	}
}

func (d dbMock) scanHosts(val string) (hosts []datastructs.Host, err error) {
	switch {
	// Error missing param
	case val == "":
		return hosts, fmt.Errorf("missing param")
	case strings.Contains("host", val):
		return []datastructs.Host{testHost1, testHost2, testHost3}, nil
	case strings.Contains("host1", val) || strings.Contains("1.1.1.1", val):
		return []datastructs.Host{testHost1}, nil
	case strings.Contains("host2", val) || strings.Contains("2.2.2.2", val):
		return []datastructs.Host{testHost2}, nil
	case strings.Contains("host3", val) || strings.Contains("3.3.3.3", val):
		return []datastructs.Host{testHost3}, nil
	// return empty if host does not exists
	default:
		return hosts, nil
	}
}

// Host-group

func (d dbMock) insertHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	switch {
	// Duplicate record
	case hostGroup.HostID == 1 && hostGroup.GroupID == 1:
		return 0, nil
	// Simulate forign key missing
	case hostGroup.HostID >= 4 || hostGroup.GroupID >= 4:
		return 0, fmt.Errorf("error")
	default:
		return 1, nil
	}
}

func (d dbMock) getHostGroups() (hosts []datastructs.HostGroup, err error) {
	return []datastructs.HostGroup{testHostGroup1, testHostGroup2, testHostGroup3}, nil
}

func (d dbMock) selectHostGroup(host string) (hostGroups []datastructs.HostGroup, err error) {
	switch {
	// Existing record
	case host == "host1":
		return []datastructs.HostGroup{testHostGroup1}, nil
	case host == "host2":
		return []datastructs.HostGroup{testHostGroup2}, nil
	// The rest of the request should return empty
	default:
		return hostGroups, nil
	}
}

func (d dbMock) deleteHostGroup(hostGroup *datastructs.HostGroup) (affected int64, err error) {
	switch {
	case hostGroup.HostID == 1 && hostGroup.GroupID == 1:
		return 1, nil
	default:
		return 0, nil
	}
}

func (d dbMock) scanHostGroups(val string) (hostGroups []datastructs.HostGroup, err error) {
	switch {
	case val == "":
		return hostGroups, fmt.Errorf("missing param")
	case strings.Contains("group", val):
		return []datastructs.HostGroup{testHostGroup1, testHostGroup2, testHostGroup3}, nil
	case strings.Contains("group1", val):
		return []datastructs.HostGroup{testHostGroup1}, nil
	case strings.Contains("group2", val):
		return []datastructs.HostGroup{testHostGroup2}, nil
	case strings.Contains("group3", val):
		return []datastructs.HostGroup{testHostGroup3}, nil
	// The rest of the request should return empty
	default:
		return hostGroups, nil
	}
}

// Child-group

func (d dbMock) insertChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error) {
	switch {
	// Duplicate record
	case childGroup.ChildID == 3 && childGroup.ParentID == 4:
		return 0, nil
	case childGroup.ChildID == 4 && childGroup.ParentID == 5:
		return 0, nil
	// Simulate forign key missing
	case childGroup.ChildID > 5 || childGroup.ParentID > 5:
		return 0, fmt.Errorf("Group does not exist")
	default:
		return 1, nil
	}
}

func (d dbMock) selectChildGroup(child, parent string) (childGroups []datastructs.ChildGroup, err error) {
	switch {
	case child == "group3" && parent == "group4":
		return []datastructs.ChildGroup{testChild1}, nil
	case child == "group4" && parent == "group5":
		return []datastructs.ChildGroup{testChild2}, nil
	default:
		return childGroups, nil
	}
}

func (d dbMock) getChildGroups() (childGroups []datastructs.ChildGroup, err error) {
	return []datastructs.ChildGroup{testChild1, testChild2}, nil
}

func (d dbMock) deleteChildGroup(childGroup *datastructs.ChildGroup) (affected int64, err error) {
	switch {
	case childGroup.ChildID == 3 || childGroup.ParentID == 4:
		return 1, nil
	case childGroup.ChildID == 4 || childGroup.ParentID == 5:
		return 1, nil
	default:
		return 0, nil
	}
}

func (d dbMock) scanChildGroups(val string) (childGroups []datastructs.ChildGroup, err error) {
	switch {
	case val == "":
		return childGroups, fmt.Errorf("missing param")
	case strings.Contains("group", val):
		return []datastructs.ChildGroup{testChild1, testChild2}, nil
	case strings.Contains("group3", val):
		return []datastructs.ChildGroup{testChild1}, nil
	case strings.Contains("group4", val):
		return []datastructs.ChildGroup{testChild1, testChild2}, nil
	case strings.Contains("group5", val):
		return []datastructs.ChildGroup{testChild2}, nil
	// The rest of the request should return empty
	default:
		return childGroups, nil
	}
}
