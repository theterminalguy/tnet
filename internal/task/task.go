package task

import (
	"fmt"
	"os"
	"strings"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/util/collection"
)

type Tasker interface {
	Run(params string) error
}

var AllowedExecutors = []string{
	"sp@10hourlabs.com",
	"abiodun.solomon@10hourlabs.com",
}

func Run(name, params, executor string) error {
	// executor is only required in production
	if os.Getenv("ENV") == "production" {
		if strings.Contains(name, "fake") {
			// do not run fake tasks in production
			return nil
		}
		if executor == "" {
			return fmt.Errorf("executor is required")
		}
		if !collection.Contains(AllowedExecutors, executor) {
			return fmt.Errorf("executor %s is not allowed", executor)
		}
	}
	if task, ok := Lookup[name]; ok {
		if err := task.Run(params); err != nil {
			tenlog.Error(fmt.Sprintf("Error running task %s: %s", name, err))
			return err
		}
		tenlog.Info(fmt.Sprintf("Task %s completed successfully", name))
		return nil
	}
	tenlog.Error(fmt.Sprintf("Task %s not found", name))
	return fmt.Errorf("task %s not found", name)
}
