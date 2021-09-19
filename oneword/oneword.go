package oneword

const (
	// Database Fields
	URL                   string = "url"
	UUID                  string = "uuid"
	City                  string = "city"
	Name                  string = "name"
	Note                  string = "note"
	Email                 string = "email"
	Phone                 string = "phone"
	Pronoun               string = "pronoun"
	LastName              string = "last_name"
	CreatedAt             string = "created_at"
	UpdatedAt             string = "updated_at"
	DeletedAt             string = "deleted_at"
	Preferred             string = "preferred"
	FirstName             string = "first_name"
	TenTNCode             string = "tentn_code"
	CoutryCode            string = "country_code"
	ReferralCode          string = "referral_code"
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
)
