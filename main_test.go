package assert

import (
	"errors"
	"testing"
)

func TestEqual(t *testing.T) {
	Equal(1, 1, t)
	Equal(1.0, 1.0, t)
	Equal(nil, nil, t)
	Equal([]uint32{1, 2, 3}, []uint32{1, 2, 3}, t)
}

func TestEqualWithin(t *testing.T) {
	EqualWithin(1.0, 1.0+1e-16, 1e-15, t)
	EqualWithin([]float64{1, 1 + 1e-16}, []float64{1, 1}, 1e-15, t)
}

func TestSuccess(t *testing.T) {
	Success(nil, t)
}

func TestFailure(t *testing.T) {
	Failure(errors.New("had a bad day"), t)
}
