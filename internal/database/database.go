package database

import (
	"fmt"

	"github.com/10hourlabs/tentn/ent"
)

type Databaser interface {
	Open() (*ent.Client, error)
	GetDSN() string
	RunMigration() error
}

type DBError struct {
	Summary string
	Err     error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("%s: %v", e.Summary, e.Err)
}

type CreateSchemaError = DBError

func NewCreateSchemaError(err error) *CreateSchemaError {
	return &CreateSchemaError{
		Summary: "failed creating schema resources",
		Err:     err,
	}
}

type ConnectionError = DBError

func NewConnectionError(err error) *ConnectionError {
	return &ConnectionError{
		Summary: "failed opening connection to database",
		Err:     err,
	}
}
