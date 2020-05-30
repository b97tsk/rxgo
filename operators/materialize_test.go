package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestMaterialize(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.Materialize(), operators.Count()),
			rx.Throw(ErrTest).Pipe(operators.Materialize(), operators.Count()),
			rx.Just("A", "B", "C").Pipe(operators.Materialize(), operators.Count()),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(operators.Materialize(), operators.Count()),
		},
		[][]interface{}{
			{1, rx.Completed},
			{1, rx.Completed},
			{4, rx.Completed},
			{4, rx.Completed},
		},
	)
}
