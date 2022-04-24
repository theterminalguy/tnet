package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathToScope(t *testing.T) {
	tests := []struct {
		path string
		verb string
		want string
	}{
		{
			path: "/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060",
			verb: "GET",
			want: "recruiter.talents.read",
		},
		{
			path: "/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060/jobs",
			verb: "GET",
			want: "recruiter.talents.read",
		},
		{
			path: "/v1/recruiter/talents/4565b3fd-ff30-4ce4-b278-e019ef298060/jobs/1323",
			verb: "GET",
			want: "recruiter.talents.read",
		},
		{
			path: "/v1/recruiter/talent-collections/4565b3fd-ff30-4ce4-b278-e019ef298060",
			verb: "POST",
			want: "recruiter.talent-collections.create",
		},
		{
			path: "/v1/recruiter/talent-collections/4565b3fd-ff30-4ce4-b278-e019ef298060",
			verb: "PUT",
			want: "recruiter.talent-collections.update",
		},
		{
			path: "/v1/recruiter/talent-collections/4565b3fd-ff30-4ce4-b278-e019ef298060",
			verb: "DELETE",
			want: "recruiter.talent-collections.delete",
		},
		{
			path: "/v1/recruiter/talent-collections/4565b3fd-ff30-4ce4-b278-e019ef298060/jobs",
			verb: "GET",
			want: "recruiter.talent-collections.read",
		},
		{
			path: "/v1/recruiter/talent-collections/search",
			verb: "GET",
			want: "recruiter.talent-collections.search",
		},
		{
			path: "/v1/recruiter/talent-collections/dummy",
			verb: "POST",
			want: "recruiter.talent-collections.dummy",
		},
		{
			path: "/v1/recruiter/talent-collections",
			verb: "GET",
			want: "recruiter.talent-collections.reads",
		},
		{
			path: "/v1/recruiter/talent-collections",
			verb: "POST",
			want: "recruiter.talent-collections.creates",
		},
		{
			path: "/v1/recruiter/talent-collections",
			verb: "PUT",
			want: "recruiter.talent-collections.updates",
		},
		{
			path: "/v1/recruiter/talent-collections",
			verb: "DELETE",
			want: "recruiter.talent-collections.deletes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := PathToScope(tt.path, tt.verb)
			assert.Equal(t, tt.want, got)
		})
	}
}
