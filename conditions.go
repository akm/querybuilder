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
