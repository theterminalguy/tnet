package model

type Job struct {
	ID int
}

func (j *Job) Hello() string {
	return "Hello from Job Struct"
}
