package utils

import (
	"errors"
	. "reflect"
	"strconv"
)

var (
	NotPtrError   = errors.New("arg is not ptr")
	FieldSetError = errors.New("field type not support error")
)

func SetValue(fieldValue Value, headerValue string) error {
	var err error
	if len(headerValue) > 0 {
		switch fieldValue.Kind() {
		case Bool:
			var x bool
			if x, err = strconv.ParseBool(headerValue); err == nil {
				fieldValue.SetBool(x)
			}
		case Int, Int8, Int16, Int32, Int64:
			var x int64
			if x, err = strconv.ParseInt(headerValue, 10, 64); err == nil {
				fieldValue.SetInt(x)
			}
		case Uint, Uint8, Uint16, Uint32, Uint64:
			var x uint64
			if x, err = strconv.ParseUint(headerValue, 10, 64); err == nil {
				fieldValue.SetUint(x)
			}
		case Float32, Float64:
			var x float64
			if x, err = strconv.ParseFloat(headerValue, 64); err == nil {
				fieldValue.SetFloat(x)
			}
		case String:
			fieldValue.SetString(headerValue)
		default:
			return FieldSetError
		}
	}
	return err
}

func IsPtr(i interface{}) error {
	t := TypeOf(i)
	if t.Kind() != Ptr {
		return NotPtrError
	}
	return nil
}

func PtrRange(ptr interface{}, handler func(t StructField, v Value) error) error {
	if err := IsPtr(ptr); err != nil {
		return err
	}
	t := TypeOf(ptr).Elem()
	value := ValueOf(ptr).Elem()

	for i := 0; i < t.NumField(); i++ {
		var (
			currField = value.Field(i)
			tf        = t.Field(i)
		)
		if err := handler(tf, currField); err != nil {
			return err
		}
	}
	return nil
}
