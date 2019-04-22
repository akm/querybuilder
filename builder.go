package querybuilder

import (
	"google.golang.org/appengine/datastore"
)

type FilterFunc func(*datastore.Query) *datastore.Query

type QueryBuilder struct {
	Fields     Strings
	ignored    Strings
	sortFields Strings
	filters    []FilterFunc
	assigns    AssignFuncs
}

func New(fields ...string) *QueryBuilder {
	return &QueryBuilder{Fields: fields}
}

func (qb *QueryBuilder) Eq(field string, value interface{}) *QueryBuilder {
	qb.filters = append(qb.filters, func(q *datastore.Query) *datastore.Query {
		return q.Filter(field+EQ.String(), value)
	})
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
	qb.filters = append(qb.filters, func(q *datastore.Query) *datastore.Query {
		return q.Filter(field+ope.String(), value)
	})
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

func (qb *QueryBuilder) ProjectFields() Strings {
	return qb.Fields.Except(qb.ignored)
}

func (qb *QueryBuilder) Build(q *datastore.Query) (*datastore.Query, AssignFuncs) {
	for _, f := range qb.filters {
		q = f(q)
	}
	for _, f := range qb.sortFields {
		q = q.Order(f)
	}
	{
		fields := qb.ProjectFields()
		if len(fields) > 0 {
			q = q.Project(fields...)
		}
	}
	return q, qb.assigns
}
