package repository

// Scopes are three levels deep, separated by dots e.g. recruiter.talents.read.
// The last part of the scope is the action, e.g. recruiter.talents.read
// Actions can be plural or singular, e.g. recruiter.talents.read or recruiter.talents.read
// Singluar actions denotes a single resourece
// Plural actions denotes a collection of resources
// Example actions:
// 		search
// 		reads (plural)
// 		creates (plural)
// 		updates (plural)
// 		deletes (plural)
// 		read (singular)
// 		create (singular)
// 		update (singular)
// 		delete (singular)

type Oauth2Scope struct {
	Name    string
	Summary string
}

type Oauth2Scopes map[string]Oauth2Scope

var scopes = map[string]Oauth2Scope{
	"recruiter.talents.search": {
		Name:    "recruiter.talents.search",
		Summary: "Search for talents",
	},
}

// Oauth2ScopeRepository is a repository for Oauth2 scopes

type Oauth2ScopeRepository struct{}

func NewOauth2ScopeRepository() *Oauth2ScopeRepository {
	return &Oauth2ScopeRepository{}
}

func (repo *Oauth2ScopeRepository) GetAll() Oauth2Scopes {
	return scopes
}

func (repo *Oauth2ScopeRepository) IsValid(scope string) bool {
	if _, ok := scopes[scope]; ok {
		return true
	}
	return false
}
