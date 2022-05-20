package main

import (
	"fmt"
	"log"
	"os"

	"github.com/10hourlabs/tentn/internal/task"
)

func main() {
	taskName := os.Args[1]
	params := ""
	if len(os.Args) > 2 {
		params = os.Args[2]
	}
	fmt.Printf("Run task %s with params %s? [y/n] ", taskName, params)
	var answer string
	fmt.Scanln(&answer)
	if answer == "y" {
		fmt.Println("Runnig task...")
		executor := ""
		err := task.Run(taskName, params, executor)
		if err != nil {
			log.Fatalf("Error running task: %s", err)
		}
		fmt.Println("Task completed successfully")
	}
}
