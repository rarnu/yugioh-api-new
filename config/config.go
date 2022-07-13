package config

import (
	c0 "github.com/isyscore/isc-gobase/config"
)

type RushStorageConfig struct {
	Jp string `yaml:"jp"`
	Cn string `yaml:"cn"`
}

var MySQLConfig c0.StorageConnectionConfig
var SQLiteConfig c0.StorageConnectionConfig
var RushDuelConfig RushStorageConfig

func LoadDatabaseConfig() {
	_ = c0.GetValueObject("mysql", &MySQLConfig)
	_ = c0.GetValueObject("sqlite", &SQLiteConfig)
	_ = c0.GetValueObject("rushduel", &RushDuelConfig)
}
