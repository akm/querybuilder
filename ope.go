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

var Operators = []Ope{LT, LTE, GT, GTE, EQ}
var OperatorMap = BuildOperatorMap(Operators)

func BuildOperatorMap(opes []Ope) map[string]Ope {
	r := map[string]Ope{}
	for _, ope := range opes {
		r[string(ope)] = ope
	}
	return r
}
