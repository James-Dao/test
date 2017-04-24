package config

import (
	"github.com/cloudfoundry/gosteno"
	"os"
)

type Config struct {
	LogLevelString    string `json:"LogLevelString"`
	Index             string `json:"Index"`
	DB                string `json:"DB"`
	Influxdb_Host     string `json:"Influxdb_Host"`
	Influxdb_Username string `json:"Influxdb_Username"`
	Influxdb_Password string `json:"Influxdb_Password"`
	Influxdb_Database string `json:"Influxdb_Database"`
}

func defaults() Config {
	return Config{
		LogLevelString: "INFO",
	}
}

func FromEnv() (*Config, error) {
	log_level := os.Getenv("LogLevelString")
	db := os.Getenv("DB")
	index := os.Getenv("Index")
	Influxdb_Host := os.Getenv("Influxdb_Host")
	Influxdb_Username := os.Getenv("Influxdb_Username")
	Influxdb_Password := os.Getenv("Influxdb_Password")
	Influxdb_Database := os.Getenv("Influxdb_Database")

	c := &Config{
		LogLevelString:    log_level,
		DB:                db,
		Index:             index,
		Influxdb_Host:     Influxdb_Host,
		Influxdb_Username: Influxdb_Username,
		Influxdb_Password: Influxdb_Password,
		Influxdb_Database: Influxdb_Database,
	}
	return c, nil
}
func (conf *Config) LogLevel() gosteno.LogLevel {
	switch conf.LogLevelString {
	case "INFO":
		return gosteno.LOG_INFO
	case "DEBUG":
		return gosteno.LOG_DEBUG
	default:
		return gosteno.LOG_INFO
	}
}
