package main

import (
	"fmt"

	"github.com/10hourlabs/jobsapi/internal/job"
)

func main() {
	j := &job.Job{}
	fmt.Println(j.Hello())
}
