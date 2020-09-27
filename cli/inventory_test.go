package cli

import (
	"reflect"
	"testing"
)

var inv = `{"_meta":{"hostvars":{"host1.local":{"var1":"val1"},"host2.local":{"var2":"val2"}}},"group1":{"hosts":["host1.local"],"vars":{"var1":"val1"}},"group2":{"children":["group1"],"vars":{"var2":"val2"}},"group3":{"children":["group2"],"vars":{"var3":"val3"}}}`

func TestGenInventory(t *testing.T) {
	db = dbMock{}

	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "Generate inventory",
			want:    []byte(inv),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenInventory()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}
