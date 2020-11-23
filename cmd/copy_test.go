package cmd

import (
	"testing"
)

func Test_copyHostCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "missing param",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many params",
			args:    []string{"host1", "host10", "extra-param"},
			wantErr: true,
		},
		{
			name:    "source does not exists",
			args:    []string{"host9", "host10"},
			wantErr: true,
		},
		{
			name:    "valid source without domain",
			args:    []string{"host1", "host10"},
			wantErr: false,
		},
		{
			name:    "new with domain",
			args:    []string{"host1", "host10.domain.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyHostCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("copyHostCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_copyGroupCase(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid new",
			args:    []string{"group1", "group10"},
			wantErr: false,
		},
		{
			name:    "missing param",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many params",
			args:    []string{"group1", "group10", "extra-param"},
			wantErr: true,
		},
		{
			name:    "source does not exists",
			args:    []string{"group9", "group10"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyGroupCase(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("copyGroupCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
