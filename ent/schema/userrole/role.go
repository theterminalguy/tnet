package userrole

type Role string

const (
	Talent    Role = "talent"
	Recruiter Role = "recruiter"
	Developer Role = "developer"
	// Admin      Role = "admin" // TODO: an admin account is not yet implemented requires business discussions
	// SuperAdmin Role = "superadmin" // TODO: a superadmin account is not yet implemented, requires business discussions
)

func (Role) Values() (kinds []string) {
	for _, k := range []Role{Talent, Recruiter, Developer} {
		kinds = append(kinds, string(k))
	}
	return
}
