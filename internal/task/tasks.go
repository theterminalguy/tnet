package task

var Lookup = map[string]Tasker{
	"tokgen":              NewTokGen(),
	"create-fake-talents": NewCreateFakeTalents(),
	"create-fake-jobs":    NewCreateFakeJob(),
	"create-fake-skills":  NewCreateFakeSkill(),
	"import-talents":      NewImportTalents(),
	"update-talent-pp":    NewUpdateProfilePicture(),
	"approve-client":      NewTaskApproveClient(),
}
