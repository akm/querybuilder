package querybuilder

import (
	"cloud.google.com/go/datastore"
)

type ValuedFilter struct {
	Name     string `json:"name"`
	IntValue int    `json:"value"`
}

func (vf *ValuedFilter) Call(q *datastore.Query) *datastore.Query {
	switch vf.Name {
	case "offset":
		return q.Offset(vf.IntValue)
	case "limit":
		return q.Limit(vf.IntValue)
	default:
		return q
	}
}
