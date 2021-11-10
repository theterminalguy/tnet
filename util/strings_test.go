package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitlelizeUnderscore(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		want         string
	}{
		{
			name:         "hello_world example",
			resourceName: "hello_world",
			want:         "HelloWorld",
		},
		{
			name:         "random_word example",
			resourceName: "random_word",
			want:         "RandomWord",
		},
		{
			name:         "job_applicant example",
			resourceName: "job_applicant",
			want:         "JobApplicant",
		},
		{
			name:         "single word 'job' example",
			resourceName: "job",
			want:         "Job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TitlelizeUnderscore(tt.resourceName)
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestRemoveUnderscore(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		want         string
	}{
		{
			name:         "hello_world example",
			resourceName: "hello_world",
			want:         "helloworld",
		},
		{
			name:         "random_word example",
			resourceName: "random_word",
			want:         "randomword",
		},
		{
			name:         "job_applicant example",
			resourceName: "job_applicant",
			want:         "jobapplicant",
		},
		{
			name:         "single word 'job' example",
			resourceName: "job",
			want:         "job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveUnderscore(tt.resourceName)
			assert.Equal(t, tt.want, got)
		})
	}
}
