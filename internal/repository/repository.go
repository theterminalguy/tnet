package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/theterminalguy/tenlog"
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/internal/database"
	"github.com/theterminalguy/tentn/util/osutil"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	// database connection
	dBConn *ent.Client

	// database query context
	dBContext context.Context

	ErrRecordDeleted error = errors.New("entity has been deleted")
)

func init() {
	// TODO: should we have this here or in main?
	osutil.CheckEnv()
	service := "dev"
	dsn := "tenlog.log"
	if os.Getenv("ENV") == "production" || os.Getenv("ENV") == "staging" {
		service = "sentry"
		dsn = os.Getenv("SENTRY_DSN")
	}
	tenlog.SetLogger(service, dsn)
	tenlog.SetAppName("tentn")
	tenlog.SetEnvName(os.Getenv("ENV"))

	var client *ent.Client
	var err error
	dBContext = context.Background()
	if os.Getenv("ENV") == "staging" || os.Getenv("ENV") == "test" {
		client, err = database.NewSQLite3InMemoryClient()
	} else {
		client, err = database.NewPostgresClient()
	}
	if err != nil {
		panic(fmt.Sprintf("Database Error %v", err))
	}
	dBConn = client
}

func slugify(title string, id uuid.UUID) string {
	return slug.Make(fmt.Sprintf("%v %v", title, id))
}

func ValidateParams(s interface{}, fields ...string) error {
	validate := validator.New()
	if len(fields) > 0 {
		if err := validate.StructPartial(s, fields...); err != nil {
			return err
		}
	} else {
		if err := validate.Struct(s); err != nil {
			return err
		}
	}
	return nil
}

func setNillableStringField(val string, cb func(v string) error) error {
	if val != "" {
		err := cb(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func setNillableJSONArrayField(vals []string, cb func(v []string) error) error {
	if len(vals) > 0 {
		err := cb(vals)
		if err != nil {
			return err
		}
	}
	return nil
}

func setNillableYearsOfExperience(val *float32, cb func(v *float32) error) error {
	if val != nil {
		err := cb(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func setNillableBoolField(val bool, cb func(v bool) error) error {
	err := cb(val)
	if err != nil {
		return err
	}
	return nil
}

func GetDBContext() context.Context {
	return dBContext
}

//TODO: This is a linear implementation and should be optimized
func LinearCheckElemArray(a, b []string) bool {
	if len(a) > len(b) {
		return false
	}
	m, n := 0, len(a)
	for i := 0; i < len(b); i++ {
		for j := 0; j < n; j++ {
			if a[j] == b[i] {
				m++
			}
		}
	}
	return m >= n
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
