package task

var Lookup = map[string]Tasker{
	"create_fake_talents": NewCreateFakeTalents(),
}
