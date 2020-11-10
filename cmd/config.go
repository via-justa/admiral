package cmd

import (
	"log"

	"github.com/via-justa/admiral/datastructs"

	"github.com/spf13/viper"
)

// Config configuration for admiral client
type config struct {
	Defaults DefaultsConfig
}

//DefaultsConfig specific configurations from file
type DefaultsConfig struct {
	Domain    string
	Monitored bool
	Enabled   bool
}

// NewConfig initialize new configuration
func newConfig() *config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/admiral")
	viper.SetConfigName(".admiral")
	viper.AddConfigPath("$HOME")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not read config")
	}

	conf := new(config)

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Could not unmarshal config")
	}

	return conf
}

func (conf *config) newDefaultHost() datastructs.Host {
	return datastructs.Host{
		Domain:    conf.Defaults.Domain,
		Monitored: conf.Defaults.Monitored,
		Enabled:   conf.Defaults.Enabled,
	}
}

func (conf *config) newDefaultGroup() datastructs.Group {
	return datastructs.Group{
		Monitored: conf.Defaults.Monitored,
		Enabled:   conf.Defaults.Enabled,
	}
}
