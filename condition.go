package querybuilder

import (
	"google.golang.org/appengine/datastore"
)

type Condition struct {
	Field string      `json:"field"`
	Ope   Ope         `json:"ope"`
	Value interface{} `json:"value"`
}

func (c *Condition) Call(q *datastore.Query) *datastore.Query {
	return q.Filter(c.Field+c.Ope.String(), c.Value)
}
