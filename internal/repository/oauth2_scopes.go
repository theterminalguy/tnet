package repository

import "github.com/10hourlabs/tentn/ent/schema/userrole"

type Oauth2ScopeCategory string

const (
	RecruiterJob              Oauth2ScopeCategory = "Recruiter Jobs"
	RecruiterTalent           Oauth2ScopeCategory = "Recruiter Talent Search"
	RecruiterTalentCollection Oauth2ScopeCategory = "Recruiter Talent Collections"
	RecruiterJobApplication   Oauth2ScopeCategory = "Recruiter Job Applications"
)

type Oauth2Scope struct {
	ID          string              `json:"name"`
	Description string              `json:"description"`
	Category    Oauth2ScopeCategory `json:"category"`
}

var OauthScopes = map[userrole.Role][]Oauth2Scope{
	userrole.Developer: {
		{
			ID:          "recruiters/jobs.read.own",
			Description: "Allow developers to read their own jobs",
			Category:    RecruiterJob,
		},
		{
			ID:          "recruiters/jobs.write.own",
			Description: "Allow developers to write their own jobs",
			Category:    RecruiterJob,
		},
		{
			ID:          "recruiter/talents.read.all",
			Description: "Allow developers to read all talents",
			Category:    RecruiterTalent,
		},
		{
			ID:          "recruiter/talent-collections.read.own",
			Description: "Allow developers to read their own talent collections",
			Category:    RecruiterTalentCollection,
		},
		{
			ID:          "recruiter/talent-collections.write.own",
			Description: "Allow developers to write their own talent collections",
			Category:    RecruiterTalentCollection,
		},
		{
			ID:          "recruiter/job-applications.read.own",
			Description: "Allow developers to read their own job applications",
			Category:    RecruiterJobApplication,
		},
		{
			ID:          "recruiter/job-applications.write.own",
			Description: "Allow developers to write their own job applications",
			Category:    RecruiterJobApplication,
		},
	},
}
