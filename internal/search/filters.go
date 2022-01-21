package search

type Filter string

var (
	UUID_EQ  Filter = "uuid_eq"
	UUID_NEQ Filter = "uuid_neq"

	SLUG_EQ  Filter = "slug_eq"
	SLUG_NEQ Filter = "slug_neq"

	TITLE_EQ  Filter = "title_eq"
	TITLE_NEQ Filter = "title_neq"

	NAME_EQ  Filter = "name_eq"
	NAME_NEQ Filter = "name_eq"

	YEAR_EXP_EQ  Filter = "year_exp_eq"
	YEAR_EXP_NEQ Filter = "year_exp_neq"

	EMAIL    Filter = "email"
	EMAIL_EQ Filter = "email_eq"

	CITY    Filter = "city"
	CITY_EQ Filter = "city_eq"

	COUNTRY Filter = "country"

	TENTN_CODE_EQ Filter = "tentn_code_eq"

	STATUS_EQ  Filter = "status_eq"
	STATUS_NEQ Filter = "status_neq"

	REFRERRAL_CODE_EQ Filter = "referral_code_eq"

	YEARS_OF_EXPERIENCE    Filter = "years_of_experience"
	YEARS_OF_EXPERIENCE_EQ Filter = "years_of_experience_eq"

	SKILLS_IN Filter = "skills_in"

	PREFERRED_TITLE_LIKE Filter = "preferred_title_like"
)
