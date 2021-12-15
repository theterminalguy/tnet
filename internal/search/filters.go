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

	EMAIL_EQ Filter = "email_eq"

	TENTN_CODE_EQ Filter = "tentn_code_eq"

	STATUS_EQ  Filter = "status_eq"
	STATUS_NEQ Filter = "status_neq"

	REFRERRAL_CODE_EQ Filter = "referral_code_eq"
)
