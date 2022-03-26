package task

var Lookup = map[string]Tasker{
	"create-fake-users":   NewCreateFakeUsers(),
	"create-fake-talents": NewCreateFakeTalents(),
	"create-fake-jobs":    NewCreateFakeJob(),
	"create-fake-skills":  NewCreateFakeSkill(),
}
