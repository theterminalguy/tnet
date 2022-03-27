package task

var Lookup = map[string]Tasker{
	"tokgen":              NewTokGen(),
	"create-fake-talents": NewCreateFakeTalents(),
	"create-fake-jobs":    NewCreateFakeJob(),
	"create-fake-skills":  NewCreateFakeSkill(),
}
