package cmd

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

var testDefaultConfig = DefaultsConfig{
	Domain:    "domain.local",
	Monitored: true,
	Enabled:   true,
}

var testConf = config{
	Defaults: testDefaultConfig,
}

func Test_config_newDefaultHost(t *testing.T) {
	tests := []struct {
		name     string
		defaults DefaultsConfig
		want     datastructs.Host
	}{
		{
			name:     "new host from config",
			defaults: testDefaultConfig,
			want: datastructs.Host{
				Domain:    "domain.local",
				Monitored: true,
				Enabled:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{
				Defaults: tt.defaults,
			}
			if got := conf.newDefaultHost(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.newDefaultHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_config_newDefaultGroup(t *testing.T) {
	tests := []struct {
		name     string
		defaults DefaultsConfig
		want     datastructs.Group
	}{
		{
			name:     "new group from config",
			defaults: testDefaultConfig,
			want: datastructs.Group{
				Monitored: true,
				Enabled:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{
				Defaults: tt.defaults,
			}
			if got := conf.newDefaultGroup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.newDefaultGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
