package querybuilder

import (
	"google.golang.org/appengine/datastore"
)

type QueryBuilder struct {
	Fields     Strings         `json:"fields,omitempty"`
	Ignored    Strings         `json:"ignored,omitempty"`
	SortFields Strings         `json:"sort_fields,omitempty"`
	Conditions []*Condition    `json:"conditions,omitempty"`
	Filters    []*ValuedFilter `json:"filters,omitempty"`
	Assigns    Assigners       `json:"assigns,omitempty"`
}

func New(fields ...string) *QueryBuilder {
	return &QueryBuilder{Fields: fields}
}

func (qb *QueryBuilder) AddCondition(field string, ope Ope, value interface{}) *QueryBuilder {
	qb.Conditions = append(qb.Conditions, &Condition{Field: field, Ope: ope, Value: value})
	return qb
}

func (qb *QueryBuilder) AddIntFilter(name string, value int) *QueryBuilder {
	qb.Filters = append(qb.Filters, &ValuedFilter{Name: name, IntValue: value})
	return qb
}

func (qb *QueryBuilder) Eq(field string, value interface{}) *QueryBuilder {
	qb.AddCondition(field, EQ, value)
	qb.Assigns = append(qb.Assigns, AssignerFor(field, value))
	qb.Ignored = append(qb.Ignored, field)
	return qb
}

func (qb *QueryBuilder) Lt(field string, value interface{}) *QueryBuilder {
	return qb.Ineq(LT, field, value)
}

func (qb *QueryBuilder) Lte(field string, value interface{}) *QueryBuilder {
	return qb.Ineq(LTE, field, value)
}

func (qb *QueryBuilder) Gt(field string, value interface{}) *QueryBuilder {
	return qb.Ineq(GT, field, value)
}

func (qb *QueryBuilder) Gte(field string, value interface{}) *QueryBuilder {
	return qb.Ineq(GTE, field, value)
}

func (qb *QueryBuilder) Ineq(ope Ope, field string, value interface{}) *QueryBuilder {
	qb.AddCondition(field, ope, value)
	if !qb.SortFields.Has(field) {
		qb.SortFields = append([]string{field}, qb.SortFields...)
	}
	return qb
}

const utf8LastChar = "\xef\xbf\xbd"

func (qb *QueryBuilder) Starts(field, value string) *QueryBuilder {
	return qb.Gte(field, value).Lte(field, value+utf8LastChar)
}

func (qb *QueryBuilder) Asc(field string) *QueryBuilder {
	return qb.AddSort(field)
}

func (qb *QueryBuilder) Desc(field string) *QueryBuilder {
	return qb.AddSort("-" + field)
}

func (qb *QueryBuilder) AddSort(field string) *QueryBuilder {
	qb.SortFields = append(qb.SortFields, field)
	return qb
}

func (qb *QueryBuilder) Offset(v int) *QueryBuilder {
	return qb.AddIntFilter("offset", v)
}

func (qb *QueryBuilder) Limit(v int) *QueryBuilder {
	return qb.AddIntFilter("limit", v)
}

func (qb *QueryBuilder) ProjectFields() Strings {
	return qb.Fields.Except(qb.Ignored)
}

func (qb *QueryBuilder) BuildForCount(q *datastore.Query) *datastore.Query {
	for _, f := range qb.Conditions {
		q = f.Call(q)
	}
	return q
}

func (qb *QueryBuilder) BuildForList(q *datastore.Query) (*datastore.Query, Assigners) {
	for _, f := range qb.SortFields {
		q = q.Order(f)
	}
	{
		fields := qb.ProjectFields()
		if len(fields) > 0 {
			q = q.Project(fields...)
		}
	}
	for _, f := range qb.Filters {
		q = f.Call(q)
	}
	return q, qb.Assigns
}

func (qb *QueryBuilder) Build(q *datastore.Query) (*datastore.Query, Assigners) {
	return qb.BuildForList(qb.BuildForCount(q))
}
