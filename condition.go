package querybuilder

import (
	"reflect"

	"cloud.google.com/go/datastore"
)

type Condition struct {
	Field string      `json:"field"`
	Ope   Ope         `json:"ope"`
	Value interface{} `json:"value"`
}

func (c *Condition) Call(q *datastore.Query) *datastore.Query {
	return q.Filter(c.Field+c.Ope.String(), c.OriginalTypeValue())
}

var primitiveTypeMap = map[reflect.Kind]reflect.Type{
	reflect.Bool:    reflect.TypeOf(bool(false)),
	reflect.Int:     reflect.TypeOf(int(0)),
	reflect.Int8:    reflect.TypeOf(int8(0)),
	reflect.Int16:   reflect.TypeOf(int16(0)),
	reflect.Int32:   reflect.TypeOf(int32(0)),
	reflect.Int64:   reflect.TypeOf(int64(0)),
	reflect.Uint:    reflect.TypeOf(uint(0)),
	reflect.Uint8:   reflect.TypeOf(uint8(0)),
	reflect.Uint16:  reflect.TypeOf(uint16(0)),
	reflect.Uint32:  reflect.TypeOf(uint32(0)),
	reflect.Uint64:  reflect.TypeOf(uint64(0)),
	reflect.Uintptr: reflect.TypeOf(uintptr(0)),
	reflect.Float32: reflect.TypeOf(float32(0)),
	reflect.Float64: reflect.TypeOf(float64(0)),
}

func (c *Condition) OriginalTypeValue() interface{} {
	v := reflect.ValueOf(c.Value)
	pt, ok := primitiveTypeMap[v.Type().Kind()]
	if ok && pt != nil {
		return v.Convert(pt).Interface()
	} else {
		return c.Value
	}
}
