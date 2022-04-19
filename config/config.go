package config

import (
	c0 "github.com/isyscore/isc-gobase/config"
)

var MySQLConfig c0.StorageConnectionConfig
var SQLiteConfig c0.StorageConnectionConfig

func LoadDatabaseConfig() {
	_ = c0.GetValueObject("mysql", &MySQLConfig)
	_ = c0.GetValueObject("sqlite", &SQLiteConfig)
}
