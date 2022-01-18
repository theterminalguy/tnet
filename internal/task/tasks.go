package task

var Lookup = map[string]Tasker{
	"create_fake_talents": NewCreateFakeTalents(),
	"create_fake_jobs": NewCreateFakeJob(),
}


