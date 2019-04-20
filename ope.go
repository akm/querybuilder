package querybuilder

type Ope string

const (
	LT  Ope = "<"
	LTE Ope = "<="
	GT  Ope = ">"
	GTE Ope = ">="
	EQ  Ope = "="
)

func (ope Ope) String() string {
	return string(ope)
}
