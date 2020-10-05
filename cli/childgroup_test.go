package cli

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func TestCreateChildGroup(t *testing.T) {
	db = dbMock{}

	type args struct {
		parent datastructs.Group
		child  datastructs.Group
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert New",
			args: args{
				parent: datastructs.Group{ID: 3},
				child:  datastructs.Group{ID: 2},
			},
			wantErr: false,
		},
		{
			name: "Insert Duplicate",
			args: args{
				parent: datastructs.Group{ID: 2},
				child:  datastructs.Group{ID: 1},
			},
			wantErr: true,
		},
		{
			name: "Insert none-existing",
			args: args{
				parent: datastructs.Group{ID: 5},
				child:  datastructs.Group{ID: 1},
			},
			wantErr: true,
		},
		{
			name: "Child and Parent the same",
			args: args{
				parent: datastructs.Group{ID: 1},
				child:  datastructs.Group{ID: 1},
			},
			wantErr: true,
		},
		{
			name: "Relationship loop",
			args: args{
				parent: datastructs.Group{ID: 1, Name: "group1"},
				child:  datastructs.Group{ID: 3, Name: "group3"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateChildGroup(tt.args.parent, tt.args.child); (err != nil) != tt.wantErr {
				t.Errorf("CreateChildGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestViewChildGroupsByParent(t *testing.T) {
	db = dbMock{}

	type args struct {
		parent string
	}

	tests := []struct {
		name            string
		args            args
		wantChildGroups []datastructs.ChildGroupView
		wantErr         bool
	}{
		{
			name: "Get child-groups",
			args: args{
				parent: "group3",
			},
			wantChildGroups: []datastructs.ChildGroupView{
				datastructs.ChildGroupView{ID: 2, ChildID: 2, Child: "group2", ParentID: 3, Parent: "group3"},
			},
			wantErr: false,
		},
		{
			name: "Get none-existing child-groups",
			args: args{
				parent: "group1",
			},
			wantChildGroups: nil,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChildGroups, err := ViewChildGroupsByParent(tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewChildGroupsByParent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("ViewChildGroupsByParent() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func TestViewChildGroupsByChild(t *testing.T) {
	db = dbMock{}

	type args struct {
		child string
	}

	tests := []struct {
		name            string
		args            args
		wantChildGroups []datastructs.ChildGroupView
		wantErr         bool
	}{
		{
			name: "Get child-groups",
			args: args{
				child: "group1",
			},
			wantChildGroups: []datastructs.ChildGroupView{
				datastructs.ChildGroupView{ID: 1, ChildID: 1, Child: "group1", ParentID: 2, Parent: "group2"},
			},
			wantErr: false,
		},
		{
			name: "Get none-existing child-groups",
			args: args{
				child: "group3",
			},
			wantChildGroups: nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChildGroups, err := ViewChildGroupsByChild(tt.args.child)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewChildGroupsByChild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("ViewChildGroupsByChild() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func TestListChildGroups(t *testing.T) {
	db = dbMock{}

	tests := []struct {
		name            string
		wantChildGroups []datastructs.ChildGroupView
		wantErr         bool
	}{
		{
			name: "List child-groups",
			wantChildGroups: []datastructs.ChildGroupView{
				datastructs.ChildGroupView{ID: 1, ChildID: 1, Child: "group1", ParentID: 2, Parent: "group2"},
				datastructs.ChildGroupView{ID: 2, ChildID: 2, Child: "group2", ParentID: 3, Parent: "group3"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChildGroups, err := ListChildGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListChildGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChildGroups, tt.wantChildGroups) {
				t.Errorf("ListChildGroups() = %v, want %v", gotChildGroups, tt.wantChildGroups)
			}
		})
	}
}

func TestDeleteChildGroup(t *testing.T) {
	db = dbMock{}

	type args struct {
		childGroup datastructs.ChildGroup
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
				childGroup: datastructs.ChildGroup{Child: 1, Parent: 2},
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "delete none-existing child-group",
			args: args{
				childGroup: datastructs.ChildGroup{Child: 2, Parent: 3},
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := DeleteChildGroup(tt.args.childGroup)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteChildGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("DeleteChildGroup() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}
