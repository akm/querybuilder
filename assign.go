package querybuilder

import (
	"fmt"
	"reflect"
	"strings"
)

type AssignFunc func(interface{}) error
type AssignFuncs []AssignFunc

func (s AssignFuncs) Assign(entity interface{}) error {
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
		v := reflect.ValueOf(value)
		switch e.Type().Kind() {
		case reflect.Struct:
			err := ReflectWalkIn(&e, field, ".", func(f *reflect.Value) error {
				f.Set(v)
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("Entity type: %T is not a struct. %v", entity, entity)
		}
	}
}

func ReflectWalkIn(base *reflect.Value, field, sep string, f func(*reflect.Value) error) error {
	fields := strings.Split(field, sep)
	return ReflectWalkInImpl(base, fields, f)
}

func ReflectWalkInImpl(curr *reflect.Value, fields []string, f func(*reflect.Value) error) error {
	if len(fields) < 1 {
		return f(curr)
	}

	switch curr.Type().Kind() {
	case reflect.Slice, reflect.Array:
		l := curr.Len()
		for i := 0; i < l; i++ {
			v := curr.Index(i)
			err := ReflectWalkInImpl(&v, fields, f)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		field := curr.FieldByName(fields[0])
		if !field.IsValid() {
			return fmt.Errorf("%s has no field named %s (from Entity %s %v)", curr.String(), fields[0])
		}
		return ReflectWalkInImpl(&field, fields[1:], f)
	default:
		return fmt.Errorf("%s is not struct but %v", curr.String(), curr.Interface())
	}
	return nil
}
