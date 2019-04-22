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
		e := reflect.ValueOf(entity)
		if e.Type().Kind() == reflect.Ptr {
			e = e.Elem()
		}
		switch e.Type().Kind() {
		case reflect.Struct:
			v := reflect.ValueOf(value)
			f := e.FieldByName(field)
			if f.IsValid() {
				f.Set(v)
			} else {
				panic(fmt.Sprintf("Entity type: %T Field: %s not found. %v", entity, field, entity))
			}
		default:
			panic(fmt.Sprintf("Entity type: %T is not a struct. %v", entity, entity))
		}
	}
}
