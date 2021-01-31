// nolint
package cmd

import (
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func init() {
	User = testUser{}
}

func Test_deleteHostCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid",
			args:    []string{"host1"},
			wantErr: false,
		},
		{
			name:    "host does not exists",
			args:    []string{"host10"},
			wantErr: true,
		},
		{
			name:    "too many hosts",
			args:    []string{"host"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteHostCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("deleteHostCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteHost(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		host *datastructs.Host
	}
	tests := []struct {
		name         string
		args         args
		wantAffected int64
		wantErr      bool
	}{
		{
			name: "Delete host",
			args: args{
				host: &testHost1,
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "Delete non-existing host",
			args: args{
				host: &testHost10,
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := deleteHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("deleteHost() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func Test_deleteGroupCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid",
			args:    []string{"group1"},
			wantErr: false,
		},
		{
			name:    "group does not exists",
			args:    []string{"group10"},
			wantErr: true,
		},
		{
			name:    "too many groups",
			args:    []string{"group"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteGroupCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("deleteGroupCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		group *datastructs.Group
	}
	tests := []struct {
		name         string
		args         args
		wantAffected int64
		wantErr      bool
	}{
		{
			name: "Delete group",
			args: args{
				group: &testGroup1,
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "Delete non-existing group",
			args: args{
				group: &testGroup10,
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := deleteGroup(tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("deleteGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

var testChild10 = datastructs.ChildGroup{
	ID:       10,
	ChildID:  1,
	Child:    "group1",
	ParentID: 2,
	Parent:   "group2",
}

func Test_deleteChildCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid",
			args:    []string{"group4", "group5"},
			wantErr: false,
		},
		{
			name:    "does not exists",
			args:    []string{"group1", "group2"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteChildCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("deleteChildCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteChildGroup(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	type args struct {
		childGroup *datastructs.ChildGroup
	}
	tests := []struct {
		name         string
		args         args
		wantAffected int64
		wantErr      bool
	}{
		{
			name: "delete child-group",
			args: args{
				childGroup: &testChild1,
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "delete none-existing child-group",
			args: args{
				childGroup: &testChild10,
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := deleteChildGroup(tt.args.childGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteChildGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("deleteChildGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}
