package stream

import (
	"fmt"
	"reflect"
)

func isPtr(v interface{}) bool {
	return reflect.ValueOf(v).Type().Kind() == reflect.Ptr
}

func settableValue(v interface{}) reflect.Value {
	return reflect.ValueOf(v).Elem()
}

func isList(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Slice || v.Type().Kind() == reflect.Array
}

func setValue(in reflect.Value, receiver reflect.Value) error {
	if !in.Type().AssignableTo(receiver.Type()) {
		return fmt.Errorf("%v is not assignable to %v", in.Type(), receiver.Type())
	}
	receiver.Set(in)
	return nil
}

func appendSlice(in reflect.Value, slice reflect.Value) error {
	if slice.Type().Kind() != reflect.Slice {
		return fmt.Errorf("%v is not slice", slice.Type())
	}
	if !in.Type().AssignableTo(slice.Type().Elem()) {
		return fmt.Errorf("%v is not assignable to %v", in.Type(), slice.Elem().Type())
	}
	slice.Set(reflect.Append(slice, in))
	return nil
}
