package config

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

var (
	testDefaultConfig = DefaultsConfig{
		Domain:    "domain.local",
		Monitored: true,
		Enabled:   true,
	}

	testMariadbConfig = MariaDBConfig{
		User:     "root",
		Password: "local",
		Host:     "localhost:3306",
		DB:       "ansible",
	}

	testSQLiteConfig = SQLiteConfig{
		Path:   "admiral.sqlite",
		Memory: true,
	}

	testConfig = Config{
		SQLite:   testSQLiteConfig,
		MariaDB:  testMariadbConfig,
		Defaults: testDefaultConfig,
	}
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "new fom file",
			want: &testConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
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
			conf := &Config{
				Defaults: tt.defaults,
			}
			if got := conf.NewDefaultHost(); !reflect.DeepEqual(got, tt.want) {
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
			conf := &Config{
				Defaults: tt.defaults,
			}
			if got := conf.NewDefaultGroup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.newDefaultGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
