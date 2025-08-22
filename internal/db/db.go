package db

import (
	"charity/internal/config"
	"fmt"
)

func GetDBSource(config *config.Config, dbName string) string {
	sslMode := ""
	if config.SSLMode == "disable" {
		sslMode = "?sslmode=disable"
	} else {
		sslMode = "?sslmode=require"
	}
	// return the structure "postgres://root:secret@localhost:5432/${db_name}?sslmode=disable"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, dbName, sslMode)
}
