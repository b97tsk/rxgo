package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestDelay(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 5).Pipe(operators.Delay(Step(3))),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(
				AddLatencyToNotifications(0, 3),
				operators.Delay(Step(1)),
			),
		},
		[][]interface{}{
			{1, 2, 3, 4, rx.Complete},
			{1, 2, 3, 4, ErrTest},
		},
	)
}
