package stream

import (
	"fmt"
	"reflect"
)

// ToList transform elements in stream into Slice or Array.
// If the receiver of Reduce is a Slice, elements will be appended to the tail of the Slice.
// If the receiver of Reduce is an Array, the array will be filled with elements in stream.
// The assigned length to Array is the minimum of Array length and element size in stream.
func ToList() ReduceFunc {
	return func(in []interface{}, out interface{}) error {
		if !isPtr(out) {
			return fmt.Errorf("%T is not assignable", out)
		}
		valOut := settableValue(out)
		if !isList(valOut) {
			return fmt.Errorf("%T is not slice or array type", out)
		}
		if valOut.Type().Kind() == reflect.Slice {
			for ix := range in {
				if err := appendSlice(reflect.ValueOf(in[ix]), valOut); err != nil {
					return err
				}
			}
			return nil
		}
		if valOut.Type().Kind() == reflect.Array {
			arrayLen := valOut.Type().Len()
			inputLen := len(in)
			for i := 0; i < arrayLen && i < inputLen; i++ {
				if err := setValue(reflect.ValueOf(in[i]), valOut.Index(i)); err != nil {
					return err
				}
			}
			return nil
		}
		return nil
	}
}

// First transform the first element in stream to the given type as the out param passed to Reduce.
// If stream is empty, out will be untouched.
func First() ReduceFunc {
	return func(in []interface{}, out interface{}) error {
		if !isPtr(out) {
			return fmt.Errorf("%T is not assignable", out)
		}
		valOut := settableValue(out)
		if isList(valOut) {
			return fmt.Errorf("%T should not be slice or array type", out)
		}
		if len(in) == 0 {
			return nil
		}
		return setValue(reflect.ValueOf(in[0]), valOut)
	}
}

// Last transform the first element in stream to the given type as the out param passed to Reduce.
// If stream is empty, out will be untouched.
func Last() ReduceFunc {
	return func(in []interface{}, out interface{}) error {
		if !isPtr(out) {
			return fmt.Errorf("%T is not assignable", out)
		}
		valOut := settableValue(out)
		if isList(valOut) {
			return fmt.Errorf("%T should not be slice or array type", out)
		}
		if len(in) == 0 {
			return nil
		}
		return setValue(reflect.ValueOf(in[len(in)-1]), valOut)
	}
}
