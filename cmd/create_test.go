// nolint
package cmd

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func init() {
	User = testUser{}
}

func Test_createHostCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid create new",
			args:    []string{"host10"},
			wantErr: false,
		},
		{
			name:    "missing param",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many params",
			args:    []string{"host1", "extra-param"},
			wantErr: true,
		},
		{
			name:    "new with domain",
			args:    []string{"host10.domain.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createHostCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("createHostCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var emptyHost10 = `{
  "ip": "",
  "hostname": "host10",
  "domain": "domain.local",
  "variables": {},
  "enable": true,
  "monitor": true,
  "direct_group": ""
}`

var testHost1Edit = `{
  "ip": "1.1.1.1",
  "hostname": "host1",
  "domain": "domain.local",
  "variables": {
    "host_var1": {
      "host_sub_var1": "host_sub_val1"
    }
  },
  "enable": true,
  "monitor": true,
  "direct_group": "group1"
}`

func Test_returnHosts(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		val string
	}
	tests := []struct {
		name      string
		args      args
		wantHosts []datastructs.Host
		wantErr   bool
	}{
		{
			name: "Existing hostname",
			args: args{
				val: "host1",
			},
			wantHosts: []datastructs.Host{testHost1},
			wantErr:   false,
		},
		{
			name: "Existing FQDN",
			args: args{
				val: "host1.domain.local",
			},
			wantHosts: []datastructs.Host{testHost1},
			wantErr:   false,
		},
		{
			name: "None-existing FQDN",
			args: args{
				val: "host1.com",
			},
			wantHosts: []datastructs.Host{datastructs.Host{Domain: "domain.local", Monitored: true, Enabled: true}},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHosts, err := returnHosts(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("returnHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHosts, tt.wantHosts) {
				t.Errorf("returnHosts() = %v, want %v", gotHosts, tt.wantHosts)
			}
		})
	}
}

func Test_prepHostForEdit(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tmpHost := testHost1

	type args struct {
		hosts    *datastructs.Host
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		wantB   []byte
		wantErr bool
	}{
		{
			name: "Hostname does not exists",
			args: args{
				hosts:    &datastructs.Host{},
				hostname: "host10",
			},
			wantB:   []byte(emptyHost10),
			wantErr: false,
		},
		{
			name: "Hostname exists",
			args: args{
				hosts:    &testHost1,
				hostname: "host1",
			},
			wantB:   []byte(testHost1Edit),
			wantErr: false,
		},
		{
			name: "FQDN exists",
			args: args{
				hosts:    &testHost1,
				hostname: "host1.domain.local",
			},
			wantB:   []byte(testHost1Edit),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := prepHostForEdit(tt.args.hosts, tt.args.hostname)
			if (err != nil) != tt.wantErr {
				testHost1 = tmpHost
				t.Errorf("prepHostForEdit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				testHost1 = tmpHost
				t.Errorf("prepHostForEdit() = %s, want %s", gotB, tt.wantB)
			}
			testHost1 = tmpHost
		})
	}
}

func Test_confirmedHost(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		host *datastructs.Host
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := confirmedHost(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("confirmedHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_viewHostByHostname(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		hostname string
	}
	tests := []struct {
		name     string
		args     args
		wantHost datastructs.Host
		wantErr  bool
	}{
		{
			name: "View host by HostName",
			args: args{
				hostname: "host1",
			},
			wantHost: testHost1,
			wantErr:  false,
		},
		{
			name: "None existing host",
			args: args{
				hostname: "host10",
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
		{
			name: "Empty host HostName",
			args: args{
				hostname: "",
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, err := viewHostByHostname(tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("viewHostByHostname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("viewHostByHostname() = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

var testHost10 = datastructs.Host{
	ID:          10,
	Hostname:    "host10",
	Host:        "10.10.10.10",
	Domain:      "local",
	Variables:   "{\"var10\": \"val10\"}",
	Enabled:     true,
	Monitored:   true,
	DirectGroup: "group1",
}

func Test_createHost(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		host *datastructs.Host
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert host",
			args: args{
				host: &testHost10,
			},
			wantErr: false,
		},
		{
			name: "Insert Existing host without change",
			args: args{
				host: &testHost1,
			},
			wantErr: false,
		},
		{
			name: "Change Existing host",
			args: args{
				host: &datastructs.Host{
					ID:        1,
					Hostname:  "host1",
					Host:      "1.1.1.1",
					Variables: "{\"host_var1\": {\"host_sub_var1\": \"host_sub_val1\"}}",
					Enabled:   true,
					Monitored: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing host HostName",
			args: args{
				host: &datastructs.Host{Hostname: ""},
			},
			wantErr: true,
		},
		{
			name: "Missing host Host",
			args: args{
				host: &datastructs.Host{Host: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createHost(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("createHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_viewHostGroupByHost(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		host string
	}
	tests := []struct {
		name          string
		args          args
		wantHostGroup []datastructs.HostGroup
		wantErr       bool
	}{
		{
			name: "Get Host-group",
			args: args{
				host: "host1",
			},
			wantHostGroup: []datastructs.HostGroup{testHostGroup1},
			wantErr:       false,
		},
		{
			name: "Get none-existing host-groups",
			args: args{
				host: "host10",
			},
			wantHostGroup: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostGroup, err := viewHostGroupByHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("viewHostGroupByHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroup, tt.wantHostGroup) {
				t.Errorf("viewHostGroupByHost() = %v, want %v", gotHostGroup, tt.wantHostGroup)
			}
		})
	}
}

func Test_deleteHostGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		hostGroup *datastructs.HostGroup
	}
	tests := []struct {
		name         string
		args         args
		wantAffected int64
		wantErr      bool
	}{
		{
			name: "delete host-group",
			args: args{
				hostGroup: &testHostGroup1,
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "delete none-existing host-group",
			args: args{
				hostGroup: &datastructs.HostGroup{HostID: 10, GroupID: 10},
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := deleteHostGroup(tt.args.hostGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteHostGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("deleteHostGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func Test_createHostGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		host  *datastructs.Host
		group *datastructs.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert New",
			args: args{
				host:  &testHost3,
				group: &testGroup3,
			},
			wantErr: false,
		},
		{
			name: "Insert Update group",
			args: args{
				host:  &testHost1,
				group: &testGroup2,
			},
			wantErr: false,
		},
		{
			name: "Insert none-existing host",
			args: args{
				host:  &testHost10,
				group: &testGroup1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createHostGroup(tt.args.host, tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("createHostGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createGroupCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid new",
			args:    []string{"group10"},
			wantErr: false,
		},
		{
			name:    "missing param",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many params",
			args:    []string{"group1", "extra-param"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createGroupCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("createGroupCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var emptyGroup10 = `{
  "name": "group10",
  "variables": {},
  "enable": true,
  "monitor": true
}`

var testGroup1Edit = `{
  "name": "group1",
  "variables": {
    "group_var1": {
      "group_sub_var1": "group_sub_val1"
    }
  },
  "enable": true,
  "monitor": true
}`

func Test_prepGroupForEdit(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tmpGroup := testGroup1
	type args struct {
		groups *datastructs.Group
		name   string
	}
	tests := []struct {
		name    string
		args    args
		wantB   []byte
		wantErr bool
	}{
		{
			name: "case 0 - group does not exist",
			args: args{
				groups: &datastructs.Group{},
				name:   "group10",
			},
			wantB:   []byte(emptyGroup10),
			wantErr: false,
		},
		{
			name: "case 1 - group exist",
			args: args{
				groups: &testGroup1,
				name:   "group1",
			},
			wantB:   []byte(testGroup1Edit),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := prepGroupForEdit(tt.args.groups, tt.args.name)
			if (err != nil) != tt.wantErr {
				testGroup1 = tmpGroup
				t.Errorf("prepGroupForEdit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				testGroup1 = tmpGroup
				t.Errorf("prepGroupForEdit() = %s, want %s", gotB, tt.wantB)
			}
			testGroup1 = tmpGroup
		})
	}
}

var testGroup10 = datastructs.Group{
	ID:        10,
	Name:      "group10",
	Variables: "{\"var10\": \"val10\"}",
	Enabled:   true,
	Monitored: true,
}

func Test_createGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tmpGroup := testGroup1
	type args struct {
		group *datastructs.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert group",
			args: args{
				group: &testGroup10,
			},
			wantErr: false,
		},
		{
			name: "Insert Existing group without change",
			args: args{
				group: &testGroup1,
			},
			wantErr: false,
		},
		{
			name: "Change Existing group",
			args: args{
				group: &datastructs.Group{
					Name:      "group1",
					Variables: "{\"group_var1\": {\"group_sub_var1\": \"group_sub_val1\"}, \"var2\": \"val2\"}",
					Enabled:   true,
					Monitored: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing group name",
			args: args{
				group: &datastructs.Group{Name: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createGroup(tt.args.group); (err != nil) != tt.wantErr {
				testGroup1 = tmpGroup
				t.Errorf("createGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
			testGroup1 = tmpGroup
		})
	}
}

func Test_createChildCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid new",
			args:    []string{"group5", "group1"},
			wantErr: false,
		},
		{
			name:    "already exists",
			args:    []string{"group4", "group5"},
			wantErr: true,
		},
		{
			name:    "child does not exists",
			args:    []string{"group10", "group1"},
			wantErr: true,
		},
		{
			name:    "parent does not exists",
			args:    []string{"group1", "group10"},
			wantErr: true,
		},
		{
			name:    "relationship loop",
			args:    []string{"group5", "group4"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createChildCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("createChildCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createChildGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		parent *datastructs.Group
		child  *datastructs.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert New",
			args: args{
				parent: &testGroup2,
				child:  &testGroup1,
			},
			wantErr: false,
		},
		{
			name: "Insert Duplicate",
			args: args{
				parent: &testGroup4,
				child:  &testGroup3,
			},
			wantErr: true,
		},
		{
			name: "Insert none-existing",
			args: args{
				parent: &testGroup10,
				child:  &testGroup1,
			},
			wantErr: true,
		},
		{
			name: "Child and Parent the same",
			args: args{
				parent: &testGroup1,
				child:  &testGroup1,
			},
			wantErr: true,
		},
		{
			name: "Relationship loop",
			args: args{
				parent: &testGroup3,
				child:  &testGroup5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createChildGroup(tt.args.parent, tt.args.child); (err != nil) != tt.wantErr {
				t.Errorf("createChildGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
