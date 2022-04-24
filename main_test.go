package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// assert.Equal(t, "recruiter.talents.search", PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "GET"))
// 	assert.Equal(t, "recruiter.talents.read", PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "POST"))
// 	assert.Equal(t, "recruiter.talents.update", PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "PUT"))
// 	assert.Equal(t, "recruiter.talents.delete", PathToScope("/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060", "DELETE"))
// 	assert.Equal(t, "recruiter.talents.reads", PathToScope("/v1/recruiter/talents", "GET"))
// 	assert.Equal(t, "recruiter.talents.creates", PathToScope("/v1/recruiter/talents", "POST"))
// 	assert.Equal(t, "recruiter.talents.updates", PathToScope("/v1/recruiter/talents", "PUT"))
// 	assert.Equal(t, "recruiter.talents.deletes", PathToScope("/v1/recruiter/talents", "DELETE"))
// 	assert.Equal(t, "recruiter.talents.search", PathToScope("/v1/recruiter/talents/search", "GET"))
// 	assert.Equal(t, "recruiter.talents.read", PathToScope("/v1/recruiter/talents/search", "POST"))
func TestPathToScope(t *testing.T) {
	tests := []struct {
		path    string
		verb    string
		want    string
	}{
		{
			path:    "/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060",
			verb:    "GET",
			want:    "recruiter.talents.read",
		},
		{
			path:    "/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060/jobs",
			verb:    "GET",
			want:    "recruiter.talents.read",
		},
		{
			path:    "/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060/jobs/1",
			verb:    "GET",
			want:    "recruiter.talents.read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := PathToScope(tt.path, tt.verb)
			assert.Equal(t, tt.want, got)
		})
	}
}
