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

	COUNTRY_CODE_EQ  Filter = "country_code_eq"
	COUNTRY_CODE_NEQ Filter = "country_code_neq"

	TENTN_CODE_EQ Filter = "tentn_code_eq"

	CREATED_AT_IN  Filter = "created_at_in"
	CREATED_NOT_IN Filter = "created_not_in"

	START Filter = "start"
	END   Filter = "end"

	COMPANY_NAME_EQ  Filter = "company_name_eq"
	COMPANY_NAME_NEQ Filter = "company_name_eq"

	LOCATION_EQ  Filter = "location_eq"
	LOCATION_NEQ Filter = "location_neq"

	JOB_TITLE_EQ  Filter = "job_title_eq"
	JOB_TITLE_NEQ Filter = "job_title_neq"

	PRIMARY_TECH_CONT Filter = "primary_tech_cont" // translates fully to primary_technologu_contains

	STATUS_EQ   Filter = "status_eq"
	STATUS_NEQ Filter = "status_neq"

	DEGREE_EQ   Filter = "degree_eq"
	DEGREE_NEQ Filter = "degree_neq"

	INST_NAME_EQ   Filter = "inst_name_eq" // where inst_name is a shorthand for INSTITUTION_NAME
	INST_NAME_NEQ Filter = "inst_name_neq"

	PROGRAM_EQ   Filter = "program_eq"
	PROGRAM_NEQ Filter = "program_neq"

	PREFRRED_EQ   Filter = "program_eq"
	PREFERRED_NEQ Filter = "program_neq"

	REFRERRAL_CODE_EQ   Filter = "referral_code_eq"
)
