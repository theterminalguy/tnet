package task

import (
	"errors"

	"github.com/theterminalguy/tentn/util/collection"
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
	"sp@theterminalguy.com": {
		Email:        "sp@theterminalguy.com",
		HashedSecret: "$2a$10$D8h9YIEiliWZ2TgEcUlme.zirTu46RW.8yBDemM9uqaq1eJM07Gwu",
		AllowedTasks: []string{
			"*",
		},
	},
	"abiodun.solomon@theterminalguy.com": {
		Email:        "abiodun.solomon@theterminalguy.com",
		HashedSecret: "$2a$10$Pd/xmOzFTu1TKfGHUFoKc.kJ4DeWjgkwoWGgBHfmZhDTTWUM8112.",
		AllowedTasks: []string{
			"tokgen",
			"import-talents",
			"update-talent-pp",
		},
	},
	"fortune.nwankwo@theterminalguy.com": {
		Email:        "fortune.nwankwo@theterminalguy.com",
		HashedSecret: "$2a$12$fgS7DkQWSe8eWI4S1d.hm.mEzZ1kreatSXucv1eZv10KLrZc0BcjW",
		AllowedTasks: []string{
			"update-talent-pp",
		},
	},
	"drey.olawaye@theterminalguy.com": {
		Email:        "drey.olawaye@theterminalguy.com",
		HashedSecret: "$2a$10$/O5vZY7Q5LxcTFjo7sGtH.WsPPwV1Irx80.dyQEoEBNU159uKDXJO",
		AllowedTasks: []string{
			"import-talents",
			"update-talent-pp",
		},
	},
	"ahbot-slack@theterminalguy.com": {
		Email:        "ahbot-slack@theterminalguy.com",
		HashedSecret: "$2a$10$ka2YjqAXzfRy2nZIPJE6LuDES2inxlYlYFpSHMBKU4AdE9PJasFie",
		AllowedTasks: []string{},
	},
}
