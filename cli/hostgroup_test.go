package cli

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func TestCreateHostGroup(t *testing.T) {
	db = dbMock{}
	type args struct {
		host  datastructs.Host
		group datastructs.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert New",
			args: args{
				host:  datastructs.Host{ID: 2},
				group: datastructs.Group{ID: 2},
			},
			wantErr: false,
		},
		{
			name: "Insert Duplicate",
			args: args{
				host:  datastructs.Host{ID: 1},
				group: datastructs.Group{ID: 1},
			},
			wantErr: true,
		},
		{
			name: "Insert none-existing",
			args: args{
				host:  datastructs.Host{ID: 3},
				group: datastructs.Group{ID: 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateHostGroup(tt.args.host, tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("CreateHostGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestViewHostGroupByHost(t *testing.T) {
	db = dbMock{}
	type args struct {
		hostID int
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
				hostID: 1,
			},
			wantHostGroup: []datastructs.HostGroup{datastructs.HostGroup{ID: 1, Host: 1, Group: 1}},
			wantErr:       false,
		},
		{
			name: "Get none-existing host-groups",
			args: args{
				hostID: 2,
			},
			wantHostGroup: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostGroup, err := ViewHostGroupByHost(tt.args.hostID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewHostGroupByHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroup, tt.wantHostGroup) {
				t.Errorf("ViewHostGroupByHost() = %v, want %v", gotHostGroup, tt.wantHostGroup)
			}
		})
	}
}

func TestViewHostGroupByGroup(t *testing.T) {
	db = dbMock{}
	type args struct {
		groupID int
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
				groupID: 1,
			},
			wantHostGroup: []datastructs.HostGroup{datastructs.HostGroup{ID: 1, Host: 1, Group: 1}},
			wantErr:       false,
		},
		{
			name: "Get none-existing host-groups",
			args: args{
				groupID: 2,
			},
			wantHostGroup: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostGroup, err := ViewHostGroupByGroup(tt.args.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewHostGroupByGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroup, tt.wantHostGroup) {
				t.Errorf("ViewHostGroupByGroup() = %v, want %v", gotHostGroup, tt.wantHostGroup)
			}
		})
	}
}

func TestListHostGroup(t *testing.T) {
	db = dbMock{}
	tests := []struct {
		name           string
		wantHostGroups []datastructs.HostGroup
		wantErr        bool
	}{
		{
			name:           "List host-group",
			wantHostGroups: []datastructs.HostGroup{datastructs.HostGroup{ID: 1, Host: 1, Group: 1}},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostGroups, err := ListHostGroup()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListHostGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHostGroups, tt.wantHostGroups) {
				t.Errorf("ListHostGroup() = %v, want %v", gotHostGroups, tt.wantHostGroups)
			}
		})
	}
}

func TestDeleteHostGroup(t *testing.T) {
	db = dbMock{}
	type args struct {
		hostGroup datastructs.HostGroup
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
				hostGroup: datastructs.HostGroup{Host: 1, Group: 1},
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "delete none-existing host-group",
			args: args{
				hostGroup: datastructs.HostGroup{Host: 2, Group: 2},
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := DeleteHostGroup(tt.args.hostGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteHostGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("DeleteHostGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}
