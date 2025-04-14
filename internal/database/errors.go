package database

import (
	"fmt"
)

type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("operation %s failed: %v", e.Operation, e.Err)

}
