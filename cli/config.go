package cli

import (
	"log"

	"github.com/spf13/viper"
	"github.com/via-justa/admiral/database"
)

// Config configuration for admiral client
type Config struct {
	Database database.DatabaseConfig
}

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
