package cli

import (
	"log"
	"reflect"
	"testing"
)

var promSD = `[{"targets":["host1.local"],"labels":{"groups":["group1","group2","group3"]}}]`

func TestGenPrometheusSDFile(t *testing.T) {
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
			gotPromSDFile, err := GenPrometheusSDFile()
			log.Printf("%s", gotPromSDFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenPrometheusSDFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPromSDFile, tt.wantPromSDFile) {
				t.Errorf("GenPrometheusSDFile() = %v, want %v", gotPromSDFile, tt.wantPromSDFile)
			}
		})
	}
}
