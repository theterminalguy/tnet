package task

type Task string

const (
	// tokGen is used to generate a new token for a user
	tokGen Task = "tokgen"

	// createFakeTalents is used to create fake talents locally
	createFakeTalents Task = "create-fake-talents"

	// createFakeJobs is used to create fake jobs locally
	createFakeJobs Task = "create-fake-jobs"

	// createFakeSkills is used to create fake skills locally
	createFakeSkills Task = "create-fake-skills"

	// importTalents is used to import talents from a Google Cloud Storage bucket
	importTalents Task = "import-talents"

	// updateTalentPP is used to update talent profile pictures
	updateTalentPP Task = "update-talent-pp"

	// approveClient is used to approve an OAuth client
	approveClient Task = "approve-client"

	// updateClientScope is used to update an OAuth client's scope
	updateClientScope Task = "update-client-scope"

	// makeClientInternal is used to make an OAuth client internal
	makeClientInternal Task = "make-client-internal"
)

var Lookup = map[Task]Tasker{
	tokGen:             NewTaskTokGen(),
	createFakeTalents:  NewTaskCreateFakeTalents(),
	createFakeJobs:     NewTaskCreateFakeJob(),
	createFakeSkills:   NewTaskCreateFakeSkill(),
	importTalents:      NewTaskImportTalents(),
	updateTalentPP:     NewTaskUpdateProfilePicture(),
	approveClient:      NewTaskApproveClient(),
	updateClientScope:  NewTaskUpdateClientScope(),
	makeClientInternal: NewTaskMakeClientInternal(),
}
