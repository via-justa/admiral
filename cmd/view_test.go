// nolint:
package cmd

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func Test_listHosts(t *testing.T) {
	db = dbMock{}

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
	db = dbMock{}

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

func Test_viewGroupByName(t *testing.T) {
	db = dbMock{}

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
	db = dbMock{}

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
	db = dbMock{}

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

func Test_listChildGroups(t *testing.T) {
	db = dbMock{}

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
	db = dbMock{}

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
	db = dbMock{}

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
