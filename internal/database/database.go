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

type Error struct {
	Summary string
	Err     error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Summary, e.Err)
}

type CreateSchemaError error

func NewCreateSchemaError(err error) CreateSchemaError {
	return &Error{
		Summary: "failed creating schema resources",
		Err:     err,
	}
}

type ConnectionError error

func NewConnectionError(err error) ConnectionError {
	return &Error{
		Summary: "failed opening connection to database",
		Err:     err,
	}
}
