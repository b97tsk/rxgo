package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestReduce(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}
		return val
	}
	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 7).Pipe(operators.Reduce(max)),
			rx.Just(42).Pipe(operators.Reduce(max)),
			rx.Empty().Pipe(operators.Reduce(max)),
			rx.Range(1, 7).Pipe(operators.Reduce(sum)),
			rx.Just(42).Pipe(operators.Reduce(sum)),
			rx.Empty().Pipe(operators.Reduce(sum)),
			rx.Throw(ErrTest).Pipe(operators.Reduce(sum)),
		},
		[][]interface{}{
			{6, rx.Complete},
			{42, rx.Complete},
			{rx.Complete},
			{21, rx.Complete},
			{42, rx.Complete},
			{rx.Complete},
			{ErrTest},
		},
	)
}

func TestFold(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}
		return val
	}
	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 7).Pipe(operators.Fold(-1, max)),
			rx.Just(42).Pipe(operators.Fold(-1, max)),
			rx.Empty().Pipe(operators.Fold(-1, max)),
			rx.Range(1, 7).Pipe(operators.Fold(-1, sum)),
			rx.Just(42).Pipe(operators.Fold(-1, sum)),
			rx.Empty().Pipe(operators.Fold(-1, sum)),
			rx.Throw(ErrTest).Pipe(operators.Fold(-1, sum)),
		},
		[][]interface{}{
			{6, rx.Complete},
			{42, rx.Complete},
			{-1, rx.Complete},
			{20, rx.Complete},
			{41, rx.Complete},
			{-1, rx.Complete},
			{ErrTest},
		},
	)
}