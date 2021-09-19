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
	WeHave                string = "wehave"
	YouHave               string = "youHave"
	Summary               string = "summary"
	Pronoun               string = "pronoun"
	Location              string = "location"
	Category              string = "category"
	LastName              string = "last_name"
	CreatedAt             string = "created_at"
	UpdatedAt             string = "updated_at"
	DeletedAt             string = "deleted_at"
	Thumbnail             string = "thumbnail"
	Preferred             string = "preferred"
	FirstName             string = "first_name"
	TenTNCode             string = "tentn_code"
	Employment            string = "employment"
	CoutryCode            string = "country_code"
	ReferralCode          string = "referral_code"
	Requirements          string = "requirements"
	PreferredName         string = "preferred_name"
	JoinedTenTNAt         string = "joined_tentn_at"
	YearsOfExperience     string = "years_of_experience"
	PreferredJobTitle     string = "preferred_job_title"
	ProfessionalStartDate string = "professional_start_date"

	// Database References
	Skills         string = "skills"
	Referrer       string = "referrer"
	Referees       string = "referees"
	Applicant      string = "applicant"
	PortfolioLinks string = "portfoliolinks"

	// Foreign Keys
	ReferrerID  string = "referrer_id"
	ApplicantID string = "applicant_id"

	// Default field values
	NULL string = "NULL"
	RemoteEarth = "Remote, Earth"
)
