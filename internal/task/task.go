package task

import "log"

type Tasker interface {
	Run(params string) error
}

func Run(name, params string) {
	if task, ok := Lookup[name]; ok {
		if err := task.Run(params); err != nil {
			log.Println(err)
			return
		}
		log.Println("Task completed successfully")
		return
	}
	log.Printf("task %s not found\n", name)
}
