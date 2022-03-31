package config

import (
	c0 "github.com/isyscore/isc-gobase/config"
)

var MySQLConfig c0.StorageConnectionConfig
var SQLiteConfig c0.StorageConnectionConfig

func LoadDatabaseConfig() {
	MySQLConfig = c0.StorageConnectionConfig{
		Host:     c0.GetValueStringDefault("mysql.host", "127.0.0.1"),
		Port:     c0.GetValueIntDefault("mysql.port", 3306),
		User:     c0.GetValueStringDefault("mysql.user", "root"),
		Password: c0.GetValueStringDefault("mysql.password", "root"),
	}
	SQLiteConfig = c0.StorageConnectionConfig{
		Host: c0.GetValueStringDefault("sqlite.host", "OmegaDB.cdb"),
	}
}
