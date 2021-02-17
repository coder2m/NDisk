package utils

import (
	"reflect"
)

type (
	Header interface {
		GetHeader(key string) string
	}
)

const (
	HeaderTag = "header"
)

func BindHttpHeader(c Header, ptr interface{}) error {
	return PtrRange(ptr, func(t reflect.StructField, v reflect.Value) error {
		var (
			name string
		)
		if name = t.Tag.Get(HeaderTag); len(name) == 0 {
			name = SnakeName(t.Name)
		}
		if name == "-" {
			return nil
		}
		return SetValue(v, c.GetHeader(name))
	})
}
