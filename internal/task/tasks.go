package task

var Lookup = map[string]Tasker{
	"create-fake-talents": NewCreateFakeTalents(),
	"create-fake-jobs":    NewCreateFakeJob(),
	"create-fake-skills":  NewCreateFakeSkill(),
}
