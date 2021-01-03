// nolint
package cmd

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

// admiral view host
func Example_viewHost() {
	testDB := prepEnv()

	defer testDB.Close()

	viewAsJSON = false

	viewHost([]string{})
	// Output:
	// IP           |Hostname     |domain       |Enabled      |Monitored    |Direct Groups |Inherited Groups
	// 1.1.1.1      |host1        |domain.local |true         |true         |group1        |
	// 2.2.2.2      |host2        |domain.local |true         |true         |group2        |
	// 3.3.3.3      |host3        |domain.local |true         |true         |group3        |group4,group5
}

// admiral view host host1
func Example_viewSingleHost() {
	testDB := prepEnv()

	defer testDB.Close()

	viewAsJSON = false

	viewHost([]string{"host1"})
	// Output:
	// IP           |Hostname     |domain       |Enabled      |Monitored    |Direct Groups |Inherited Groups
	// 1.1.1.1      |host1        |domain.local |true         |true         |group1        |
}

// admiral view host host1 -j
func Example_viewSingleHostAsJSON() {
	testDB := prepEnv()

	defer testDB.Close()

	viewAsJSON = true

	viewHost([]string{"host1"})
	// Output:
	// 	[
	//     {
	//         "ip": "1.1.1.1",
	//         "hostname": "host1",
	//         "domain": "domain.local",
	//         "variables": {
	//             "host_var1": {
	//                 "host_sub_var1": "host_sub_val1"
	//             }
	//         },
	//         "enable": true,
	//         "monitor": true,
	//         "direct_group": "group1"
	//     }
	// ]
}

func Test_listHosts(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name      string
		wantHosts []datastructs.Host
		wantErr   bool
	}{
		{
			name:      "List hosts",
			wantHosts: []datastructs.Host{testHost1, testHost2, testHost3},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHosts, err := listHosts()
			if (err != nil) != tt.wantErr {
				t.Errorf("listHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHosts, tt.wantHosts) {
				t.Errorf("listHosts() = %v, want %v", gotHosts, tt.wantHosts)
			}
		})
	}
}

func Test_scanHosts(t *testing.T) {
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
			name: "get exact match host1",
			args: args{
				val: "host1",
			},
			wantHosts: []datastructs.Host{testHost1},
			wantErr:   false,
		},
		{
			name: "get substring 1",
			args: args{
				val: "1",
			},
			wantHosts: []datastructs.Host{testHost1},
			wantErr:   false,
		},
		{
			name: "get substring host",
			args: args{
				val: "host",
			},
			wantHosts: []datastructs.Host{testHost1, testHost2, testHost3},
			wantErr:   false,
		},
		{
			name: "pass empty string",
			args: args{
				val: "",
			},
			wantHosts: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHosts, err := scanHosts(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHosts, tt.wantHosts) {
				t.Errorf("scanHosts() = %v, want %v", gotHosts, tt.wantHosts)
			}
		})
	}
}

// admiral view host-group
func Example_viewHostGroup() {
	testDB := prepEnv()

	defer testDB.Close()

	viewHostGroup([]string{})
	// Output:
	// Group        | Group ID     | Hostname     | Host ID
	// group1       | 1            | host1        | 1
	// group2       | 2            | host2        | 2
	// group3       | 3            | host3        | 3
}

// admiral view host-group group1
func Example_viewSingleHostGroup() {
	testDB := prepEnv()

	defer testDB.Close()

	viewHostGroup([]string{"group1"})
	// Output:
	// Group        | Group ID     | Hostname     | Host ID
	// group1       | 1            | host1        | 1
}
func Test_listHostGroups(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		wantHg  []datastructs.HostGroup
		wantErr bool
	}{
		{
			name:    "List host groups",
			wantHg:  []datastructs.HostGroup{testHostGroup1, testHostGroup2, testHostGroup3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHg, err := listHostGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("listHostGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHg, tt.wantHg) {
				t.Errorf("listHostGroups() = %v, want %v", gotHg, tt.wantHg)
			}
		})
	}
}

func Test_scanHostGroups(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		val string
	}
	tests := []struct {
		name           string
		args           args
		wantHostGroups []datastructs.HostGroup
		wantErr        bool
	}{
		{
			name: "get exact match group1",
			args: args{
				val: "group1",
			},
			wantHostGroups: []datastructs.HostGroup{testHostGroup1},
			wantErr:        false,
		},
		{
			name: "get substring 1",
			args: args{
				val: "1",
			},
			wantHostGroups: []datastructs.HostGroup{testHostGroup1},
			wantErr:        false,
		},
		{
			name: "get substring group",
			args: args{
				val: "group",
			},
			wantHostGroups: []datastructs.HostGroup{testHostGroup1, testHostGroup2, testHostGroup3},
			wantErr:        false,
		},
		{
			name: "pass empty string",
			args: args{
				val: "",
			},
			wantHostGroups: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostGroups, err := scanHostGroups(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanHostGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroups, tt.wantHostGroups) {
				t.Errorf("scanHostGroups() = %v, want %v", gotHostGroups, tt.wantHostGroups)
			}
		})
	}
}

// admiral view group
func Example_viewGroup() {
	testDB := prepEnv()

	defer testDB.Close()

	viewAsJSON = false

	viewGroup([]string{})
	// Output:
	// name         |Enabled      |Monitored    |Children count |Hosts count
	// group1       |true         |true         |0              |1
	// group2       |true         |true         |0              |1
	// group3       |true         |true         |0              |1
	// group4       |true         |true         |1              |0
	// group5       |true         |true         |2              |0
}

// admiral view group group1
func Example_viewSingleGroup() {
	testDB := prepEnv()

	defer testDB.Close()

	viewAsJSON = false

	viewGroup([]string{"group1"})
	// Output:
	// name         |Enabled      |Monitored    |Children count |Hosts count
	// group1       |true         |true         |0              |1
}

// admiral view group group1 -j
func Example_viewSingleGroupAsJSON() {
	testDB := prepEnv()

	defer testDB.Close()

	viewAsJSON = true

	viewGroup([]string{"group1"})
	// Output:
	// [
	//     {
	//         "name": "group1",
	//         "variables": {
	//             "group_var1": {
	//                 "group_sub_var1": "group_sub_val1"
	//             }
	//         },
	//         "enable": true,
	//         "monitor": true
	//     }
	// ]
}

func Test_viewGroupByName(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantGroup datastructs.Group
		wantErr   bool
	}{
		{
			name: "View group by name",
			args: args{
				name: "group1",
			},
			wantGroup: testGroup1,
			wantErr:   false,
		},
		{
			name: "None existing group",
			args: args{
				name: "group10",
			},
			wantGroup: datastructs.Group{},
			wantErr:   true,
		},
		{
			name: "Empty group name",
			args: args{
				name: "",
			},
			wantGroup: datastructs.Group{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroup, err := viewGroupByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("viewGroupByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroup, tt.wantGroup) {
				t.Errorf("viewGroupByName() = %v, want %v", gotGroup, tt.wantGroup)
			}
		})
	}
}

func Test_listGroups(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name       string
		wantGroups []datastructs.Group
		wantErr    bool
	}{
		{
			name:       "List groups",
			wantGroups: []datastructs.Group{testGroup1, testGroup2, testGroup3, testGroup4, testGroup5},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroups, err := listGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("listGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("listGroups() = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

func Test_scanGroups(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		val string
	}
	tests := []struct {
		name       string
		args       args
		wantGroups []datastructs.Group
		wantErr    bool
	}{
		{
			name: "get exact match group1",
			args: args{
				val: "group1",
			},
			wantGroups: []datastructs.Group{testGroup1},
			wantErr:    false,
		},
		{
			name: "get substring 1",
			args: args{
				val: "1",
			},
			wantGroups: []datastructs.Group{testGroup1},
			wantErr:    false,
		},
		{
			name: "get substring group",
			args: args{
				val: "group",
			},
			wantGroups: []datastructs.Group{testGroup1, testGroup2, testGroup3, testGroup4, testGroup5},
			wantErr:    false,
		},
		{
			name: "pass empty string",
			args: args{
				val: "",
			},
			wantGroups: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroups, err := scanGroups(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("scanGroups() = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

// admiral view child
func Example_viewChildGroup() {
	testDB := prepEnv()

	defer testDB.Close()

	viewChild([]string{})
	// Output:
	// Parent       |Parent ID    |Child        |Child ID
	// group4       |4            |group3       |3
	// group5       |5            |group4       |4
}

// admiral view child group5
func Example_viewSingleChildGroup() {
	testDB := prepEnv()

	defer testDB.Close()

	viewChild([]string{"group5"})
	// Output:
	// Parent       |Parent ID    |Child        |Child ID
	// group5       |5            |group4       |4
}

func Test_listChildGroups(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name            string
		wantChildGroups []datastructs.ChildGroup
		wantErr         bool
	}{
		{
			name:            "List child-groups",
			wantChildGroups: []datastructs.ChildGroup{testChild1, testChild2},
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChildGroups, err := listChildGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("listChildGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("listChildGroups() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func Test_viewChildGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		child  string
		parent string
	}
	tests := []struct {
		name            string
		args            args
		wantChildGroups []datastructs.ChildGroup
		wantErr         bool
	}{
		{
			name: "search group with both parent and child ",
			args: args{
				child:  "group3",
				parent: "group4",
			},
			wantChildGroups: []datastructs.ChildGroup{testChild1},
			wantErr:         false,
		},
		{
			name: "Get none-existing child-groups",
			args: args{
				child:  "group3",
				parent: "group10",
			},
			wantChildGroups: nil,
			wantErr:         true,
		},
		{
			name: "missing param",
			args: args{
				child:  "",
				parent: "group10",
			},
			wantChildGroups: nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChildGroups, err := viewChildGroup(tt.args.child, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("viewChildGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("viewChildGroup() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func Test_scanChildGroups(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		val string
	}
	tests := []struct {
		name            string
		args            args
		wantChildGroups []datastructs.ChildGroup
		wantErr         bool
	}{
		{
			name: "get exact match group3",
			args: args{
				val: "group3",
			},
			wantChildGroups: []datastructs.ChildGroup{testChild1},
			wantErr:         false,
		},
		{
			name: "get substring 3",
			args: args{
				val: "3",
			},
			wantChildGroups: []datastructs.ChildGroup{testChild1},
			wantErr:         false,
		},
		{
			name: "get substring group",
			args: args{
				val: "group",
			},
			wantChildGroups: []datastructs.ChildGroup{testChild1, testChild2},
			wantErr:         false,
		},
		{
			name: "pass empty string",
			args: args{
				val: "",
			},
			wantChildGroups: nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChildGroups, err := scanChildGroups(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanChildGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("scanChildGroups() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}
