package search

type Operator string

var (
	EQ   Operator = "eq"
	LT   Operator = "lt"
	GT   Operator = "gt"
	NEQ  Operator = "neq"
	GTEQ Operator = "gteq"
	LTEQ Operator = "lteq"
)
