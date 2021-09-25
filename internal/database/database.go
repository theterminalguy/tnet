package database

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
)

func Client() (*ent.Client, error) {
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed opening connection to sqlite: %v", err))
	}
	ctx := context.Background()

	// Run the automatic migration tool to create all schema resources
	if err := client.Schema.Create(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("failed creating schema resources: %v", err))
	}
	return client, nil
}
