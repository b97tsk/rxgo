package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMergeSync(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(
			operators.MergeSyncAll(1),
		),
		"A", "B", "C", "D", "E", "F", Completed,
	).Case(
		rx.Timer(Step(1)).Pipe(
			operators.MergeSyncMapTo(rx.Just("A"), 1),
		),
		"A", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.MergeSyncAll(1),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.MergeSyncAll(1),
		),
		ErrTest,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.Just("A", "B", "C", "D", "E").Pipe(
		operators.MergeSyncMap(
			func(val interface{}, idx int) rx.Observable {
				return rx.Just(val).Pipe(DelaySubscription(2))
			},
			3,
		),
		operators.DoOnComplete(
			func() { t.Fatal("should not happen") },
		),
	).Subscribe(ctx, rx.Noop)
}
