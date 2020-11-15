package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/via-justa/admiral/datastructs"
)

//MariaDBConfig MariaDB specific configurations
type MariaDBConfig struct {
	User     string
	Password string
	Host     string
	DB       string
}

//SQLiteConfig SQLite specific configurations
type SQLiteConfig struct {
	Path   string
	Memory bool // Used for tests, if set to true the database will run in memory and be discarded after each run.
}

//DefaultsConfig specific hosts and groups default configurations
type DefaultsConfig struct {
	Domain    string
	Monitored bool
	Enabled   bool
}

// Config database configuration for admiral client
type Config struct {
	SQLite   SQLiteConfig   `toml:"sqlite" mapstructure:"sqlite"`
	MariaDB  MariaDBConfig  `toml:"mariadb" mapstructure:"mariadb"`
	Defaults DefaultsConfig `toml:"defaults" mapstructure:"defaults"`
}

// NewConfig initialize new configuration
func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/admiral")
	viper.SetConfigName(".admiral")
	viper.AddConfigPath("$HOME")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not read config")
	}

	conf := new(Config)

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Could not unmarshal config")
	}

	return conf
}

// NewDefaultHost return host with defaults from config
func (conf *Config) NewDefaultHost() datastructs.Host {
	return datastructs.Host{
		Domain:    conf.Defaults.Domain,
		Monitored: conf.Defaults.Monitored,
		Enabled:   conf.Defaults.Enabled,
	}
}

// NewDefaultGroup return group with defaults from config
func (conf *Config) NewDefaultGroup() datastructs.Group {
	return datastructs.Group{
		Monitored: conf.Defaults.Monitored,
		Enabled:   conf.Defaults.Enabled,
	}
}
