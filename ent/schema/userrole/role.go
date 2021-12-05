package userrole

type Role string

const (
	Talent     Role = "talent"
	Recruiter  Role = "recruiter"
	Admin      Role = "admin"
	SuperAdmin Role = "superadmin"
)

func (Role) Values() (kinds []string) {
	for _, k := range []Role{Talent, Recruiter, Admin, SuperAdmin} {
		kinds = append(kinds, string(k))
	}
	return
}
