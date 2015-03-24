// Package assert provides a number of functions for making assertions in
// tests.
package assert

import (
	"math"
	"reflect"
	"runtime"
	"testing"
)

// Equal asserts that two objects are equal.
func Equal(actual, expected interface{}, t *testing.T) {
	atype, etype := reflect.TypeOf(actual), reflect.TypeOf(expected)
	if atype == nil && etype == nil {
		return
	} else if (atype == nil) != (etype == nil) {
		raise(t, "got %T instead of %T", actual, expected)
	}

	kind := atype.Kind()
	if kind != etype.Kind() {
		raise(t, "got %T instead of %T", actual, expected)
	}

	switch kind {
	case reflect.Slice, reflect.Struct, reflect.Ptr:
		if !reflect.DeepEqual(actual, expected) {
			raise(t, "got %v instead of %v", actual, expected)
		}
	default:
		if actual != expected {
			raise(t, "got %v instead of %v", actual, expected)
		}
	}
}

// EqualWithin asserts that the distance between two scalars or the uniform
// distance between two vectors is less than the given value.
func EqualWithin(actual, expected interface{}, ε interface{}, t *testing.T) {
	typo := reflect.TypeOf(actual)
	if typo != reflect.TypeOf(expected) {
		raise(t, "got %v instead of %v", actual, expected)
	}

	kind := typo.Kind()
	if kind != reflect.TypeOf(expected).Kind() {
		raise(t, "got %T instead of %T", actual, expected)
	}

	avalue, evalue := reflect.ValueOf(actual), reflect.ValueOf(expected)

	if kind == reflect.Slice {
		kind = typo.Elem().Kind()
	} else {
		avalue = reflect.Append(reflect.MakeSlice(reflect.SliceOf(typo), 0, 1), avalue)
		evalue = reflect.Append(reflect.MakeSlice(reflect.SliceOf(typo), 0, 1), evalue)
	}

	if reflect.TypeOf(ε).Kind() != kind {
		raise(t, "got %T instead of %T", actual, ε)
	}

	if avalue.Len() != evalue.Len() {
		raise(t, "got %v instead of %v", actual, expected)
	}

	actual, expected = avalue.Interface(), evalue.Interface()

	switch kind {
	case reflect.Float64:
		actual, expected, ε := actual.([]float64), expected.([]float64), ε.(float64)
		for i := range actual {
			if Δ := math.Abs(actual[i] - expected[i]); Δ > ε {
				raise(t, "got %v instead of %v (delta %v)", actual[i], expected[i], Δ)
			}
		}
	default:
		panic("the type is not supported")
	}
}

// Success asserts that the error is nil.
func Success(err error, t *testing.T) {
	if err != nil {
		raise(t, "got an error '%v'", err)
	}
}

// Success asserts that the error is not nil.
func Failure(err error, t *testing.T) {
	if err == nil {
		raise(t, "expected an error")
	}
}

func raise(t *testing.T, format string, arguments ...interface{}) {
	if _, file, line, ok := runtime.Caller(2); ok {
		t.Errorf("%s:%d", file, line)
	}
	t.Fatalf(format, arguments...)
}
