package database

import (
	"fmt"
)

func InitDatabase(dbType string) error {
	// switch dbType {
	// case "sql":
	// 	return sql.InitSQLDatabase()
	// case "nosql":
	// 	return nosql.InitNoSQLDatabase()
	// default:
	// 	return fmt.Errorf("unsupported database type %s", dbType)
	// }
}

type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("operation %s failed: %v", e.Operation, e.Err)

}
