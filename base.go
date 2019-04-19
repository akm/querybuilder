package querybuilder

import (
	"reflect"

	"google.golang.org/appengine/datastore"
)

type Ope string

const (
	LT  Ope = "<"
	LTE Ope = "<="
	GT  Ope = ">"
	GTE Ope = ">="
	EQ  Ope = "="
)

type FilterFunc func(*datastore.Query) *datastore.Query
type AssignFunc func(interface{})

type Strings []string

func (s Strings) Has(v string) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}

func (s Strings) Except(v Strings) Strings {
	r := Strings{}
	for _, i := range s {
		if !v.Has(i) {
			r = append(r, i)
		}
	}
	return r
}

type QueryBuilder struct {
	Fields     Strings
	ignored    Strings
	sortFields Strings
	filters    []FilterFunc
	assigns    []AssignFunc
}

func (qb *QueryBuilder) Eq(field string, value interface{}) *QueryBuilder {
	qb.filters = append(qb.filters, func(q *datastore.Query) *datastore.Query {
		return q.Filter(field+EQ, value)
	})
	qb.assigns = append(qb.assigns, func(entiry interface{}) {
		e := reflect.ValueOf(entity)
		v := reflect.ValueOf(value)
		f := e.FieldByName(field)
		f.Set(v)
	})
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
		return q.Filter(field+EQ, value)
	})
	qb.sortFields = append([]string{field}, qb.sortFields...)
	return qb
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

func (qb *QueryBuilder) ProjectFields() Strings {
	return qb.Fields.Except(ignored)
}

func (qb *QueryBuilder) Build(q *datastore.Query) (*datastore.Query, AssignFunc) {
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
	return q, func(entity interface{}) {
		for _, f := range qb.assigns {
			f(entity)
		}
	}
}
