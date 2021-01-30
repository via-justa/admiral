package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/via-justa/admiral/datastructs"
)

// MariaDBConfig MariaDB specific configurations
type MariaDBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string
}

// SSHProxy SSH settings to proxy commends thru
type SSHProxy struct {
	User     string
	KeyPath  string `toml:"key-path" mapstructure:"key-path"`
	Host     string
	Port     int
	Password string
}

// SSH settings for ssh command
type SSH struct {
	User                  string
	KeyPath               string `toml:"key-path" mapstructure:"key-path"`
	Port                  int
	Password              string
	StrictHostKeyChecking bool `toml:"strict-host-key-checking" mapstructure:"strict-host-key-checking"`
	Proxy                 bool
}

// SQLiteConfig SQLite specific configurations
type SQLiteConfig struct {
	Path   string
	Memory bool // Used for tests, if set to true the database will run in memory and be discarded after each run.
}

// DefaultsConfig specific hosts and groups default configurations
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
	SSHProxy SSHProxy       `toml:"ssh-proxy" mapstructure:"ssh-proxy"`
	SSH      SSH            `toml:"ssh" mapstructure:"ssh"`
}

// NewConfig initialize new configuration
func NewConfig() *Config {
	var err error

	v := viper.New()
	// first check on local folder if config exists
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/admiral")
	// Override all for tests
	v.AddConfigPath("../fixtures")

	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// if couldn't find local config, search in user home folder
			v.SetConfigName(".admiral")
			v.AddConfigPath("$HOME")
			err = v.ReadInConfig()

			if err != nil {
				log.Fatal("Could not read config")
			}
		} else {
			// Config file was found but another error was produced
		}
	}

	conf := new(Config)

	err = v.Unmarshal(&conf)
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
		Variables: "{}",
	}
}

// NewDefaultGroup return group with defaults from config
func (conf *Config) NewDefaultGroup() datastructs.Group {
	return datastructs.Group{
		Monitored: conf.Defaults.Monitored,
		Enabled:   conf.Defaults.Enabled,
		Variables: "{}",
	}
}
