package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTicker(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Ticker(Step(1)).Pipe(
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return idx
				},
			),
			operators.Take(3),
		),
		0, 1, 2, Completed,
	).TestAll()
}

func TestTimer(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Timer(Step(1)).Pipe(operators.MapTo(42)),
		42, Completed,
	).TestAll()
}
