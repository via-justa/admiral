package cmd

import (
	"testing"
)

func Test_importHostsFromPath(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "import valid file",
			args:    []string{"../fixtures/files/hosts.json"},
			wantErr: false,
		},
		{
			name:    "import file does not exists",
			args:    []string{"../fixtures/files/none-existing.json"},
			wantErr: true,
		},
		{
			name:    "import corrupted file",
			args:    []string{"../fixtures/files/hosts-corrupted.json"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := importHostsFromPath(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("importHostsFromPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_importGroupsFromPath(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "import valid file",
			args:    []string{"../fixtures/files/groups.json"},
			wantErr: false,
		},
		{
			name:    "import file does not exists",
			args:    []string{"../fixtures/files/none-existing.json"},
			wantErr: true,
		},
		{
			name:    "import corrupted file",
			args:    []string{"../fixtures/files/groups-corrupted.json"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := importGroupsFromPath(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("importGroupsFromPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_importChildrenFromPath(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "import valid file",
			args:    []string{"../fixtures/files/children.json"},
			wantErr: false,
		},
		{
			name:    "import file does not exists",
			args:    []string{"../fixtures/files/none-existing.json"},
			wantErr: true,
		},
		{
			name:    "import corrupted file",
			args:    []string{"../fixtures/files/children-corrupted.json"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := importChildrenFromPath(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("importChildrenFromPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
