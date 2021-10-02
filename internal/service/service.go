package service

import (
	"github.com/go-playground/validator/v10"
)

func validateParams(s interface{}) error {
	if err := validator.New().Struct(s); err != nil {
		return err
	}
	return nil
}
