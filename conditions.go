package querybuilder

import (
	"google.golang.org/appengine/datastore"
)

type Conditions []*Condition

type ConditionPredict func(*Condition) bool

func (s Conditions) Call(q *datastore.Query) *datastore.Query {
	for _, i := range s {
		q = i.Call(q)
	}
	return q
}

func (s Conditions) IneqFields() Strings {
	r := Strings{}
	for _, i := range s {
		if i.Ope != EQ {
			r = append(r, i.Field)
		}
	}
	return r.Uniq()
}

func (s Conditions) HasMultipleIneqFields() bool {
	return len(s.IneqFields()) > 1
}
