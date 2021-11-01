package util

import "testing"

func TestTitlelizeUnderscore(t *testing.T) {
	tests := []struct{
		name string
		resourceName string
		want string
	}{
		{
			name : "hello_world example",
			resourceName: "hello_world",
			want: "HelloWorld",
		},
		{
			name : "random_word example",
			resourceName: "random_word",
			want: "RandomWord",
		},
		{
			name : "job_applicant example",
			resourceName: "job_applicant",
			want: "JobApplicant",
		},
		{
			name : "single word 'job' example",
			resourceName: "job",
			want: "Job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t * testing.T) {
			if got := TitlelizeUnderscore(tt.resourceName); got != tt.want {
				t.Errorf("TitlelizeUnderscore(%s) = %s; want %s", tt.resourceName, got, tt.want)
			}
		})
	}
}

func TestRemoveUnderscore(t *testing.T) {
	tests := []struct{
		name string
		resourceName string
		want string
	}{
		{
			name : "hello_world example",
			resourceName: "hello_world",
			want: "helloworld",
		},
		{
			name : "random_word example",
			resourceName: "random_word",
			want: "randomword",
		},
		{
			name : "job_applicant example",
			resourceName: "job_applicant",
			want: "jobapplicant",
		},
		{
			name : "single word 'job' example",
			resourceName: "job",
			want: "job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t * testing.T) {
			if got := RemoveUnderscore(tt.resourceName); got != tt.want {
				t.Errorf("TitlelizeUnderscore(%s) = %s; want %s", tt.resourceName, got, tt.want)
			}
		})
	}
}