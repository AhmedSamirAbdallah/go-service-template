package database

import (
	"fmt"
	"go-service-template/internal/database/nosql"
	"go-service-template/internal/database/sql"
)

func InitDatabase(dbType string) error {
	switch dbType {
	case "sql":
		return sql.InitSQLDatabase()
	case "nosql":
		return nosql.InitNoSQLDatabase()
	default:
		return fmt.Errorf("unsupported database type %s", dbType)
	}
}
