// nolint: rowserrcheck,lll,golint
package sqlite

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/via-justa/admiral/datastructs"
)

var testDB *Database

func prepEnv() {
	var err error

	os.Remove(dbConfig.Path)

	testDB, err = Connect(&dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = testDB.PopulateTestData("../../../fixtures")
	if err != nil {
		log.Fatal(err)
	}
}

func TestDatabase_SelectHost(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name             string
		Conn             *sqlx.DB
		hostname         string
		wantReturnedHost datastructs.Host
		wantErr          bool
	}{
		{
			name:             "get exact match host1",
			Conn:             testDB.Conn,
			hostname:         "host1",
			wantReturnedHost: testHost1,
			wantErr:          false,
		},
		{
			name:             "get substring host",
			Conn:             testDB.Conn,
			hostname:         "host",
			wantReturnedHost: datastructs.Host{},
			wantErr:          false,
		},
		{
			name:             "pass empty string",
			Conn:             testDB.Conn,
			hostname:         "",
			wantReturnedHost: datastructs.Host{},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}

			gotReturnedHost, err := db.SelectHost(tt.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.SelectHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReturnedHost, tt.wantReturnedHost) {
				t.Errorf("Database.SelectHost() = %v, want %v", gotReturnedHost, tt.wantReturnedHost)
			}
		})
	}
}

func TestDatabase_GetHosts(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name      string
		Conn      *sqlx.DB
		wantHosts []datastructs.Host
		wantErr   bool
	}{
		{
			name:      "get all hosts",
			Conn:      testDB.Conn,
			wantHosts: []datastructs.Host{testHost1, testHost2, testHost3},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotHosts, err := db.GetHosts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHosts, tt.wantHosts) {
				t.Errorf("Database.GetHosts() = %v, want %v", gotHosts, tt.wantHosts)
			}
		})
	}
}

func TestDatabase_InsertHost(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	editTestHost1 := testHost1
	editTestHost1.Enabled = false

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		host         *datastructs.Host
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "insert host10",
			Conn:         testDB.Conn,
			host:         &createTestHost10,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "update host1",
			Conn:         testDB.Conn,
			host:         &editTestHost1,
			wantAffected: 1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.InsertHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.InsertHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.InsertHost() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_DeleteHost(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		host         *datastructs.Host
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "delete host1",
			Conn:         testDB.Conn,
			host:         &testHost1,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "delete none-existing host10",
			Conn:         testDB.Conn,
			host:         &createTestHost10,
			wantAffected: 0,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.DeleteHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.DeleteHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.DeleteHost() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_ScanHosts(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name      string
		Conn      *sqlx.DB
		val       string
		wantHosts []datastructs.Host
		wantErr   bool
	}{
		{
			name:      "get exact match host1",
			Conn:      testDB.Conn,
			val:       "host1",
			wantHosts: []datastructs.Host{testHost1},
			wantErr:   false,
		},
		{
			name:      "get substring 1",
			Conn:      testDB.Conn,
			val:       "1",
			wantHosts: []datastructs.Host{testHost1},
			wantErr:   false,
		},
		{
			name:      "get substring host",
			Conn:      testDB.Conn,
			val:       "host",
			wantHosts: []datastructs.Host{testHost1, testHost2, testHost3},
			wantErr:   false,
		},
		{
			name:      "pass empty string",
			Conn:      testDB.Conn,
			val:       "",
			wantHosts: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotHosts, err := db.ScanHosts(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.ScanHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHosts, tt.wantHosts) {
				t.Errorf("Database.ScanHosts() = %v, want %v", gotHosts, tt.wantHosts)
			}
		})
	}
}

func TestDatabase_SelectGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name              string
		Conn              *sqlx.DB
		nameV             string
		wantReturnedGroup datastructs.Group
		wantErr           bool
	}{
		{
			name:              "View group group1",
			Conn:              testDB.Conn,
			nameV:             "group1",
			wantReturnedGroup: testGroup1,
			wantErr:           false,
		},
		{
			name:              "None existing group",
			Conn:              testDB.Conn,
			nameV:             "group10",
			wantReturnedGroup: datastructs.Group{},
			wantErr:           false,
		},
		{
			name:              "Empty group name",
			Conn:              testDB.Conn,
			nameV:             "",
			wantReturnedGroup: datastructs.Group{},
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotReturnedGroup, err := db.SelectGroup(tt.nameV)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.SelectGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReturnedGroup, tt.wantReturnedGroup) {
				t.Errorf("Database.SelectGroup() = %v, want %v", gotReturnedGroup, tt.wantReturnedGroup)
			}
		})
	}
}

func TestDatabase_GetGroups(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name       string
		Conn       *sqlx.DB
		wantGroups []datastructs.Group
		wantErr    bool
	}{
		{
			name:       "get all groups",
			Conn:       testDB.Conn,
			wantGroups: []datastructs.Group{testGroup1, testGroup2, testGroup3, testGroup4, testGroup5},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotGroups, err := db.GetGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("Database.GetGroups() = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

func TestDatabase_InsertGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	editGroup1 := testGroup1
	editGroup1.Enabled = false

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		group        *datastructs.Group
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "insert group10",
			Conn:         testDB.Conn,
			group:        &createTestGroup10,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "update group1",
			Conn:         testDB.Conn,
			group:        &editGroup1,
			wantAffected: 1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.InsertGroup(tt.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.InsertGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.InsertGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_DeleteGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		group        *datastructs.Group
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "delete group1",
			Conn:         testDB.Conn,
			group:        &testGroup1,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "delete none-existing group10",
			Conn:         testDB.Conn,
			group:        &createTestGroup10,
			wantAffected: 0,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.DeleteGroup(tt.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.DeleteGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_ScanGroups(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name       string
		Conn       *sqlx.DB
		val        string
		wantGroups []datastructs.Group
		wantErr    bool
	}{
		{
			name:       "get exact match group1",
			Conn:       testDB.Conn,
			val:        "group1",
			wantGroups: []datastructs.Group{testGroup1},
			wantErr:    false,
		},
		{
			name:       "get substring 1",
			Conn:       testDB.Conn,
			val:        "1",
			wantGroups: []datastructs.Group{testGroup1},
			wantErr:    false,
		},
		{
			name:       "get substring group",
			Conn:       testDB.Conn,
			val:        "group",
			wantGroups: []datastructs.Group{testGroup1, testGroup2, testGroup3, testGroup4, testGroup5},
			wantErr:    false,
		},
		{
			name:       "pass empty string",
			Conn:       testDB.Conn,
			val:        "",
			wantGroups: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotGroups, err := db.ScanGroups(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.ScanGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("Database.ScanGroups() = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

func TestDatabase_SelectChildGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name            string
		Conn            *sqlx.DB
		child           string
		parent          string
		wantChildGroups []datastructs.ChildGroup
		wantErr         bool
	}{
		{
			name:            "view group with both parent and child ",
			Conn:            testDB.Conn,
			child:           "group3",
			parent:          "group4",
			wantChildGroups: []datastructs.ChildGroup{testChild1},
			wantErr:         false,
		},
		{
			name:            "None-existing child-groups",
			Conn:            testDB.Conn,
			child:           "group3",
			parent:          "group10",
			wantChildGroups: nil,
			wantErr:         false,
		},
		{
			name:            "missing param",
			Conn:            testDB.Conn,
			child:           "",
			parent:          "group10",
			wantChildGroups: nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotChildGroups, err := db.SelectChildGroup(tt.child, tt.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.SelectChildGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("Database.SelectChildGroup() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func TestDatabase_GetChildGroups(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name            string
		Conn            *sqlx.DB
		wantChildGroups []datastructs.ChildGroup
		wantErr         bool
	}{
		{
			name:            "get all child-groups",
			Conn:            testDB.Conn,
			wantChildGroups: []datastructs.ChildGroup{testChild1, testChild2},
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotChildGroups, err := db.GetChildGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetChildGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("Database.GetChildGroups() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func TestDatabase_InsertChildGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		childGroup   *datastructs.ChildGroup
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "insert child-group10",
			Conn:         testDB.Conn,
			childGroup:   &createTestChild10,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "none-existing child-group11",
			Conn:         testDB.Conn,
			childGroup:   &createTestChild11,
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.InsertChildGroup(tt.childGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.InsertChildGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.InsertChildGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_DeleteChildGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		childGroup   *datastructs.ChildGroup
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "delete child-group1",
			Conn:         testDB.Conn,
			childGroup:   &testChild1,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "delete none-existing group10",
			Conn:         testDB.Conn,
			childGroup:   &createTestChild10,
			wantAffected: 0,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.DeleteChildGroup(tt.childGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.DeleteChildGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.DeleteChildGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_ScanChildGroups(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name            string
		Conn            *sqlx.DB
		val             string
		wantChildGroups []datastructs.ChildGroup
		wantErr         bool
	}{
		{
			name:            "get exact match group3",
			Conn:            testDB.Conn,
			val:             "group3",
			wantChildGroups: []datastructs.ChildGroup{testChild1},
			wantErr:         false,
		},
		{
			name:            "get substring 3",
			Conn:            testDB.Conn,
			val:             "3",
			wantChildGroups: []datastructs.ChildGroup{testChild1},
			wantErr:         false,
		},
		{
			name:            "get substring group",
			Conn:            testDB.Conn,
			val:             "group",
			wantChildGroups: []datastructs.ChildGroup{testChild1, testChild2},
			wantErr:         false,
		},
		{
			name:            "pass empty string",
			Conn:            testDB.Conn,
			val:             "",
			wantChildGroups: nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotChildGroups, err := db.ScanChildGroups(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.ScanChildGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("Database.ScanChildGroups() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func TestDatabase_SelectHostGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name           string
		Conn           *sqlx.DB
		host           string
		wantHostGroups []datastructs.HostGroup
		wantErr        bool
	}{
		{
			name:           "view hosts for host1",
			Conn:           testDB.Conn,
			host:           "host1",
			wantHostGroups: []datastructs.HostGroup{testHostGroup1},
			wantErr:        false,
		},
		{
			name:           "None-existing host-group",
			Conn:           testDB.Conn,
			host:           "group5",
			wantHostGroups: nil,
			wantErr:        false,
		},
		{
			name:           "missing param",
			Conn:           testDB.Conn,
			host:           "",
			wantHostGroups: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotHostGroups, err := db.SelectHostGroup(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.SelectHostGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroups, tt.wantHostGroups) {
				t.Errorf("Database.SelectHostGroup() = %v, want %v", gotHostGroups, tt.wantHostGroups)
			}
		})
	}
}

func TestDatabase_GetHostGroups(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name           string
		Conn           *sqlx.DB
		wantHostGroups []datastructs.HostGroup
		wantErr        bool
	}{
		{
			name:           "view hosts for host1",
			Conn:           testDB.Conn,
			wantHostGroups: []datastructs.HostGroup{testHostGroup1, testHostGroup2, testHostGroup3},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotHostGroups, err := db.GetHostGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetHostGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroups, tt.wantHostGroups) {
				t.Errorf("Database.GetHostGroups() = %v, want %v", gotHostGroups, tt.wantHostGroups)
			}
		})
	}
}

func TestDatabase_InsertHostGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	editTestHostGroup1 := testHostGroup1
	editTestHostGroup1.GroupID = 2

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		hostGroup    *datastructs.HostGroup
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "insert hostGroup10 (host3, group5)",
			Conn:         testDB.Conn,
			hostGroup:    &createTestHostGroup10,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "update testHostGroup1 (host1, group2)",
			Conn:         testDB.Conn,
			hostGroup:    &editTestHostGroup1,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "none-existing FK",
			Conn:         testDB.Conn,
			hostGroup:    &createTestHostGroup11Err,
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.InsertHostGroup(tt.hostGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.InsertHostGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.InsertHostGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_DeleteHostGroup(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name         string
		Conn         *sqlx.DB
		hostGroup    *datastructs.HostGroup
		wantAffected int64
		wantErr      bool
	}{
		{
			name:         "delete host-group1",
			Conn:         testDB.Conn,
			hostGroup:    &testHostGroup1,
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name:         "delete none-existing host-group10",
			Conn:         testDB.Conn,
			hostGroup:    &createTestHostGroup10,
			wantAffected: 0,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotAffected, err := db.DeleteHostGroup(tt.hostGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.DeleteHostGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("Database.DeleteHostGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestDatabase_ScanHostGroups(t *testing.T) {
	prepEnv()

	defer testDB.Conn.Close()

	tests := []struct {
		name           string
		Conn           *sqlx.DB
		val            string
		wantHostGroups []datastructs.HostGroup
		wantErr        bool
	}{
		{
			name:           "get exact match group1",
			Conn:           testDB.Conn,
			val:            "group1",
			wantHostGroups: []datastructs.HostGroup{testHostGroup1},
			wantErr:        false,
		},
		{
			name:           "get substring 1",
			Conn:           testDB.Conn,
			val:            "1",
			wantHostGroups: []datastructs.HostGroup{testHostGroup1},
			wantErr:        false,
		},
		{
			name:           "get substring group",
			Conn:           testDB.Conn,
			val:            "group",
			wantHostGroups: []datastructs.HostGroup{testHostGroup1, testHostGroup2, testHostGroup3},
			wantErr:        false,
		},
		{
			name:           "pass empty string",
			Conn:           testDB.Conn,
			val:            "",
			wantHostGroups: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Database{
				Conn: tt.Conn,
			}
			gotHostGroups, err := db.ScanHostGroups(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.ScanHostGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroups, tt.wantHostGroups) {
				t.Errorf("Database.ScanHostGroups() = %v, want %v", gotHostGroups, tt.wantHostGroups)
			}
		})
	}
}
