package oneword

const (
	// Database Fields
	URL                   string = "url"
	UUID                  string = "uuid"
	City                  string = "city"
	Name                  string = "name"
	Note                  string = "note"
	Slug                  string = "slug"
	Email                 string = "email"
	Phone                 string = "phone"
	Title                 string = "title"
	Hiring                string = "hiring"
	Status                string = "status"
	WeHave                string = "we_have"
	YouHave               string = "you_have"
	Summary               string = "summary"
	Pronoun               string = "pronoun"
	Location              string = "location"
	Category              string = "category"
	LastName              string = "last_name"
	CreatedAt             string = "created_at"
	UpdatedAt             string = "updated_at"
	DeletedAt             string = "deleted_at"
	Thumbnail             string = "thumbnail"
	Screening             string = "screening"
	Preferred             string = "preferred"
	FirstName             string = "first_name"
	TenTNCode             string = "tentn_code"
	Employment            string = "employment"
	CoutryCode            string = "country_code"
	ReferralCode          string = "referral_code"
	Requirements          string = "requirements"
	PreferredName         string = "preferred_name"
	ReferralSource        string = "referral_source"
	YearsOfExperience     string = "years_of_experience"
	PreferredJobTitle     string = "preferred_job_title"
	ProfessionalStartDate string = "professional_start_date"

	// Database References
	Job             string = "job"
	Skills          string = "skills"
	Referrer        string = "referrer"
	Referees        string = "referees"
	Applications    string = "applications"
	JobApplications string = "job_applications"
	PortfolioLinks  string = "portfoliolinks"
	Talent          string = "talent"
	Talents         string = "talents"

	// Foreign Keys
	JobID      string = "job_id"
	ReferrerID string = "referrer_id"
	TalentID   string = "talent_id"

	// Default field values
	NULL        string = "NULL"
	RemoteEarth string = "Remote, Earth"

	// Others
	Jobs string = "jobs"

	Favorite = "favorite"

	Claims           = "claims"
	CurrentUser      = "currentUser"
	CurrentTalent    = "currentTalent"
	CurrentRecruiter = "currentRecruiter"
	CurrentDeveloper = "currentDeveloper"
	ClientID         = "clientID"
)
