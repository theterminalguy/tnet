package task

import (
	"fmt"
	"os"
	"strings"

	"github.com/10hourlabs/tenlog"
)

type Tasker interface {
	Run(params string) error
}

func Run(taskName, params, executor, password string) error {
	// executor is only required in production
	if os.Getenv("ENV") == "production" {
		if strings.Contains(taskName, "fake") {
			return fmt.Errorf("error running task: %s", taskName)
		}
		if executor == "" {
			return fmt.Errorf("missing executor")
		}
		ex, ok := AllowedExecutors[executor]
		if !ok {
			return fmt.Errorf("excutor %s not found", executor)
		}
		if err := ex.Authenticate(password); err != nil {
			return err
		}
		if !ex.CanRunTask(taskName) {
			return fmt.Errorf("%s is not allowed to run task %s", ex.Email, taskName)
		}
	}
	if task, ok := Lookup[Task(taskName)]; ok {
		if err := task.Run(params); err != nil {
			tenlog.Error(fmt.Sprintf("Error running task %s: %s", taskName, err))
			return err
		}
		tenlog.Info(fmt.Sprintf("Task %s completed successfully", taskName))
		return nil
	}
	tenlog.Error(fmt.Sprintf("Task %s not found", taskName))
	return fmt.Errorf("task %s not found", taskName)
}
