package cli

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func TestCreateGroup(t *testing.T) {
	db = dbMock{}

	type args struct {
		group datastructs.Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert group",
			args: args{
				group: datastructs.Group{
					Name:      "group4",
					Variables: "{\"var4\": \"val4\"}",
					Enabled:   true,
					Monitored: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Insert Existing group without change",
			args: args{
				group: datastructs.Group{
					Name:      "group1",
					Variables: "{\"var1\": \"val1\"}",
					Enabled:   true,
					Monitored: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Change Existing group",
			args: args{
				group: datastructs.Group{
					Name:      "group1",
					Variables: "{\"var1\": \"val1\", \"var2\": \"val2\"}",
					Enabled:   true,
					Monitored: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing group name",
			args: args{
				group: datastructs.Group{Name: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateGroup(tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestViewGroupByName(t *testing.T) {
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
			wantGroup: datastructs.Group{
				ID:        1,
				Name:      "group1",
				Variables: "{\"var1\": \"val1\"}",
				Enabled:   true,
				Monitored: true,
			},
			wantErr: false,
		},
		{
			name: "None existing group",
			args: args{
				name: "group4",
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
			gotGroup, err := ViewGroupByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewGroupByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroup, tt.wantGroup) {
				t.Errorf("ViewGroupByName() = %v, want %v", gotGroup, tt.wantGroup)
			}
		})
	}
}

func TestViewGroupByID(t *testing.T) {
	db = dbMock{}

	type args struct {
		id int
	}
	tests := []struct {
		name      string
		args      args
		wantGroup datastructs.Group
		wantErr   bool
	}{
		{
			name: "View group by id",
			args: args{
				id: 1,
			},
			wantGroup: datastructs.Group{
				ID:        1,
				Name:      "group1",
				Variables: "{\"var1\": \"val1\"}",
				Enabled:   true,
				Monitored: true,
			},
			wantErr: false,
		},
		{
			name: "None existing group",
			args: args{
				id: 4,
			},
			wantGroup: datastructs.Group{},
			wantErr:   true,
		},
		{
			name: "Empty group name",
			args: args{
				id: 0,
			},
			wantGroup: datastructs.Group{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroup, err := ViewGroupByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroup, tt.wantGroup) {
				t.Errorf("ViewGroupByID() = %v, want %v", gotGroup, tt.wantGroup)
			}
		})
	}
}

func TestListGroups(t *testing.T) {
	db = dbMock{}

	tests := []struct {
		name       string
		wantGroups []datastructs.Group
		wantErr    bool
	}{
		{
			name: "List groups",
			wantGroups: []datastructs.Group{
				datastructs.Group{
					ID:        1,
					Name:      "group1",
					Variables: "{\"var1\": \"val1\"}",
					Enabled:   true,
					Monitored: true,
				},
				datastructs.Group{
					ID:        2,
					Name:      "group2",
					Variables: "{\"var2\": \"val2\"}",
					Enabled:   true,
					Monitored: true,
				},
				datastructs.Group{
					ID:        3,
					Name:      "group3",
					Variables: "{\"var3\": \"val3\"}",
					Enabled:   true,
					Monitored: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroups, err := ListGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("ListGroups() = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	db = dbMock{}

	type args struct {
		group datastructs.Group
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
				group: datastructs.Group{ID: 1},
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "Delete non-existing group",
			args: args{
				group: datastructs.Group{ID: 4},
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := DeleteGroup(tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("DeleteGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}
