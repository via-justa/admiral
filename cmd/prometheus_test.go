// nolint:
package cmd

import (
	"reflect"
	"testing"
)

var promSD = `[
    {
        "targets": [
            "host1.local"
        ],
        "labels": {
            "group": "group1",
            "inherited_groups": ""
        }
    },
    {
        "targets": [
            "host2.local"
        ],
        "labels": {
            "group": "group2",
            "inherited_groups": ""
        }
    },
    {
        "targets": [
            "host3.local"
        ],
        "labels": {
            "group": "group3",
            "inherited_groups": "group4,group5"
        }
    }
]`

func Test_genPrometheusSDFile(t *testing.T) {
	db = dbMock{}

	tests := []struct {
		name           string
		wantPromSDFile []byte
		wantErr        bool
	}{
		{
			name:           "Generate prometheus static file",
			wantPromSDFile: []byte(promSD),
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPromSDFile, err := genPrometheusSDFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("genPrometheusSDFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPromSDFile, tt.wantPromSDFile) {
				t.Errorf("genPrometheusSDFile() = %s, want %s", gotPromSDFile, tt.wantPromSDFile)
			}
		})
	}
}
