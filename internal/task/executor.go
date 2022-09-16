package task

import (
	"errors"

	"github.com/10hourlabs/tentn/util/collection"
	"golang.org/x/crypto/bcrypt"
)

type Executor struct {
	Email        string
	HashedSecret string

	// AllowedTasks is a list of tasks that this executor is allowed to run
	// "*" means all tasks are allowed by this executor
	AllowedTasks []string
}

func (e *Executor) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(e.HashedSecret), []byte(password)); err != nil {
		return errors.New("invalid username or password")
	}
	return nil
}

func (e *Executor) CanRunTask(taskName string) bool {
	if len(e.AllowedTasks) == 0 {
		return false
	}
	if e.AllowedTasks[0] == "*" {
		return true
	}
	return collection.Contains(e.AllowedTasks, taskName)
}

var AllowedExecutors = map[string]*Executor{
	"sp@10hourlabs.com": {
		Email:        "sp@10hourlabs.com",
		HashedSecret: "$2a$10$D8h9YIEiliWZ2TgEcUlme.zirTu46RW.8yBDemM9uqaq1eJM07Gwu",
		AllowedTasks: []string{
			"*",
		},
	},
	"abiodun.solomon@10hourlabs.com": {
		Email:        "abiodun.solomon@10hourlabs.com",
		HashedSecret: "$2a$10$Pd/xmOzFTu1TKfGHUFoKc.kJ4DeWjgkwoWGgBHfmZhDTTWUM8112.",
		AllowedTasks: []string{
			"tokgen",
			"import-talents",
		},
	},
	"fortune.nwankwo@10hourlabs.com": {
		Email:        "fortune.nwankwo@10hourlabs.com",
		HashedSecret: "$2a$10$duQjyZJFPluWQGvug1JUGu5uKifXTxT6eHEHesQ2cQEK7mwbCcTxC",
		AllowedTasks: []string{
			"update-talent-pp",
		},
	},
	"onwunma.anuli@10hourlabs.com": {
		Email:        "onwunma.anuli@10hourlabs.com",
		HashedSecret: "$2a$10$M1C/qWkl5rVE7q18PzUa7.CV7BYExHPl23cEc/1VHyHgp7FJJRFOS",
		AllowedTasks: []string{
			"import-talents",
		},
	},
}
