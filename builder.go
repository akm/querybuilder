package querybuilder

import (
	"google.golang.org/appengine/datastore"
)

type FilterFunc func(*datastore.Query) *datastore.Query

type Condition struct {
	Field string      `json:"field"`
	Ope   Ope         `json:"ope"`
	Value interface{} `json:"value"`
}

func (c *Condition) Call(q *datastore.Query) *datastore.Query {
	return q.Filter(c.Field+c.Ope.String(), c.Value)
}

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

type QueryBuilder struct {
	Fields     Strings
	ignored    Strings
	sortFields Strings
	conditions []*Condition
	filters    []*ValuedFilter
	assigns    AssignFuncs
}

func New(fields ...string) *QueryBuilder {
	return &QueryBuilder{Fields: fields}
}

func (qb *QueryBuilder) AddCondition(field string, ope Ope, value interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, &Condition{Field: field, Ope: ope, Value: value})
	return qb
}

func (qb *QueryBuilder) AddIntFilter(name string, value int) *QueryBuilder {
	qb.filters = append(qb.filters, &ValuedFilter{Name: name, IntValue: value})
	return qb
}

func (qb *QueryBuilder) Eq(field string, value interface{}) *QueryBuilder {
	qb.AddCondition(field, EQ, value)
	qb.assigns = append(qb.assigns, AssignFuncFor(field, value))
	qb.ignored = append(qb.ignored, field)
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
	if !qb.sortFields.Has(field) {
		qb.sortFields = append([]string{field}, qb.sortFields...)
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
	qb.sortFields = append(qb.sortFields, field)
	return qb
}

func (qb *QueryBuilder) SortFields() Strings {
	return qb.sortFields
}

func (qb *QueryBuilder) Offset(v int) *QueryBuilder {
	return qb.AddIntFilter("offset", v)
}

func (qb *QueryBuilder) Limit(v int) *QueryBuilder {
	return qb.AddIntFilter("limit", v)
}

func (qb *QueryBuilder) ProjectFields() Strings {
	return qb.Fields.Except(qb.ignored)
}

func (qb *QueryBuilder) BuildForCount(q *datastore.Query) *datastore.Query {
	for _, f := range qb.conditions {
		q = f.Call(q)
	}
	return q
}

func (qb *QueryBuilder) BuildForList(q *datastore.Query) (*datastore.Query, AssignFuncs) {
	for _, f := range qb.sortFields {
		q = q.Order(f)
	}
	{
		fields := qb.ProjectFields()
		if len(fields) > 0 {
			q = q.Project(fields...)
		}
	}
	for _, f := range qb.filters {
		q = f.Call(q)
	}
	return q, qb.assigns
}

func (qb *QueryBuilder) Build(q *datastore.Query) (*datastore.Query, AssignFuncs) {
	return qb.BuildForList(qb.BuildForCount(q))
}
