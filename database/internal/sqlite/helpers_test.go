// nolint:
package sqlite

import (
	"github.com/via-justa/admiral/config"
	"github.com/via-justa/admiral/datastructs"
)

// Simulate the following DB records

// host:
// id|host   |hostname |domain       |variables                                        |enabled|monitored|
// --|-------|---------|-------------|-------------------------------------------------|-------|---------|
//  1|1.1.1.1|host1    |domain.local |{"host_var1": {"host_sub_var1": "host_sub_val1"}}|      1|        1|
//  2|2.2.2.2|host2    |domain.local |{"host_var2": "host_val2"}                       |      1|        1|
//  3|3.3.3.3|host3    |domain.local |{"host_var3": "host_val3"}                       |      1|        1|

// group
// id|name  |variables                                           |enabled|monitored|
// --|------|----------------------------------------------------|-------|---------|
//  1|group1|{"group_var1": {"group_sub_var1": "group_sub_val1"}}|      1|        1|
//  2|group2|{"group_var2": "group_val2"}                        |      1|        1|
//  3|group3|{"group_var3": "group_val3"}                        |      1|        1|
//  4|group4|{"group_var4": "group_val4"}                        |      1|        1|
//  5|group5|{"group_var5": "group_val5"}                        |      1|        1|

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
		NumHosts:  1,
	}
	testGroup2 = datastructs.Group{
		ID:        2,
		Name:      "group2",
		Variables: "{\"group_var2\": \"group_val2\"}",
		Enabled:   true,
		Monitored: true,
		NumHosts:  1,
	}
	testGroup3 = datastructs.Group{
		ID:        3,
		Name:      "group3",
		Variables: "{\"group_var3\": \"group_val3\"}",
		Enabled:   true,
		Monitored: true,
		NumHosts:  1,
	}
	testGroup4 = datastructs.Group{
		ID:          4,
		Name:        "group4",
		Variables:   "{\"group_var4\": \"group_val4\"}",
		Enabled:     true,
		Monitored:   true,
		ChildGroups: "group3",
		NumChildren: 1,
	}
	testGroup5 = datastructs.Group{
		ID:          5,
		Name:        "group5",
		Variables:   "{\"group_var5\": \"group_val5\"}",
		Enabled:     true,
		Monitored:   true,
		ChildGroups: "group3,group4",
		NumChildren: 2,
	}
	testHost1 = datastructs.Host{
		ID:          1,
		Hostname:    "host1",
		Host:        "1.1.1.1",
		Domain:      "domain.local",
		Variables:   "{\"host_var1\": {\"host_sub_var1\": \"host_sub_val1\"}}",
		Enabled:     true,
		Monitored:   true,
		DirectGroup: "group1",
	}
	testHost2 = datastructs.Host{
		ID:          2,
		Hostname:    "host2",
		Host:        "2.2.2.2",
		Domain:      "domain.local",
		Variables:   "{\"host_var2\": \"host_val2\"}",
		Enabled:     true,
		Monitored:   true,
		DirectGroup: "group2",
	}
	testHost3 = datastructs.Host{
		ID:              3,
		Hostname:        "host3",
		Host:            "3.3.3.3",
		Domain:          "domain.local",
		Variables:       "{\"host_var3\": \"host_val3\"}",
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

var createTestHost10 = datastructs.Host{
	ID:          10,
	Hostname:    "host10",
	Host:        "10.10.10.10",
	Domain:      "domain.local",
	Variables:   "{\"host_var10\": \"host_val10\"}",
	Enabled:     true,
	Monitored:   true,
	DirectGroup: "group1",
}

var createTestGroup10 = datastructs.Group{
	ID:        1,
	Name:      "group10",
	Variables: "{\"group_var10\": {\"group_sub_var10\": \"group_sub_val10\"}}",
	Enabled:   true,
	Monitored: true,
}

var createTestChild10 = datastructs.ChildGroup{
	ID:       10,
	ChildID:  2,
	Child:    "group2",
	ParentID: 3,
	Parent:   "group3",
}

var createTestChild11 = datastructs.ChildGroup{
	ID:       11,
	ChildID:  10,
	Child:    "group10",
	ParentID: 11,
	Parent:   "group11",
}

var createTestHostGroup10 = datastructs.HostGroup{
	ID:      10,
	HostID:  3,
	Host:    "host3",
	GroupID: 5,
	Group:   "group5",
}

var createTestHostGroup11Err = datastructs.HostGroup{
	ID:      11,
	HostID:  10,
	Host:    "host10",
	GroupID: 10,
	Group:   "group10",
}

var testDefaultConfig = config.DefaultsConfig{
	Domain:    "domain.local",
	Monitored: true,
	Enabled:   true,
}

var testConf = config.Config{
	Defaults: testDefaultConfig,
}

var dbConfig = config.SQLiteConfig{
	Path:   "admiral.sqlite",
	Memory: true,
}
