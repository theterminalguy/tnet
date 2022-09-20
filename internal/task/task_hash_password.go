package task

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type TaskHashPassword struct {
}

func NewTaskHashPassword() *TaskHashPassword {
	return &TaskHashPassword{}
}

func (t *TaskHashPassword) Run(password string) error {
	if password == "" {
		return fmt.Errorf("password is empty")
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	fmt.Println(string(hashedSecret))
	return nil
}
