package querybuilder

import (
	"fmt"
	"reflect"
)

type AssignFunc func(interface{})
type AssignFuncs []AssignFunc

func (s AssignFuncs) Assign(entity interface{}) {
	fmt.Printf("entity %v\n", entity)
	for _, f := range s {
		f(entity)
	}
}

func (s AssignFuncs) AssignAll(entities interface{}) error {
	v := reflect.ValueOf(entities)
	switch v.Type().Kind() {
	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			s.Assign(v.Index(i).Interface())
		}
	case reflect.Ptr:
		return s.AssignAll(v.Elem().Interface())
	default:
		return fmt.Errorf("Unsupported type of slice %T", entities)
	}
	return nil
}

func AssignFuncFor(field string, value interface{}) AssignFunc {
	return func(entity interface{}) {
		e := reflect.Indirect(reflect.ValueOf(entity))
		v := reflect.ValueOf(value)
		f := e.FieldByName(field)
		f.Set(v)
	}
}
