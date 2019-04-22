package querybuilder

import (
	"fmt"
	"reflect"
)

type AssignFunc func(interface{}) error
type AssignFuncs []AssignFunc

func (s AssignFuncs) Assign(entity interface{}) error {
	fmt.Printf("entity %v\n", entity)
	for _, f := range s {
		if err := f(entity); err != nil {
			return err
		}
	}
	return nil
}

func (s AssignFuncs) AssignAll(entities interface{}) error {
	v := reflect.ValueOf(entities)
	switch v.Type().Kind() {
	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			if err := s.Assign(v.Index(i).Interface()); err != nil {
				return nil
			}
		}
	case reflect.Ptr:
		return s.AssignAll(v.Elem().Interface())
	default:
		return fmt.Errorf("Unsupported type of slice %T", entities)
	}
	return nil
}

func AssignFuncFor(field string, value interface{}) AssignFunc {
	return func(entity interface{}) error {
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
				return nil
			} else {
				return fmt.Errorf("Entity type: %T Field: %s not found. %v", entity, field, entity)
			}
		default:
			return fmt.Errorf("Entity type: %T is not a struct. %v", entity, entity)
		}
	}
}
