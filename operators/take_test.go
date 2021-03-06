package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTake(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			operators.Take(0),
		),
		Completed,
	).Case(
		rx.Range(1, 9).Pipe(
			operators.Take(3),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Range(1, 3).Pipe(
			operators.Take(3),
		),
		1, 2, Completed,
	).Case(
		rx.Range(1, 1).Pipe(
			operators.Take(3),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Take(0),
		),
		Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Take(3),
		),
		1, 2, 3, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 3),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Take(3),
		),
		1, 2, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 1),
			rx.Throw(ErrTest),
		).Pipe(
			operators.Take(3),
		),
		ErrTest,
	).TestAll()
}
