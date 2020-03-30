package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestZipAll(t *testing.T) {
	delay := DelaySubscription(1)
	observables := [...]rx.Observable{
		rx.Just(rx.Just("A", "B"), rx.Range(1, 4)),
		rx.Just(rx.Just("A", "B", "C"), rx.Range(1, 4)),
		rx.Just(rx.Just("A", "B", "C", "D"), rx.Range(1, 4)),
		rx.Just(rx.Just("A", "B"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
		rx.Just(rx.Just("A", "B", "C"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
		rx.Just(rx.Just("A", "B", "C", "D"), rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(delay)),
	}
	for i, obs := range observables {
		observables[i] = obs.Pipe(operators.ZipAll(), ToString())
	}
	SubscribeN(
		t,
		observables[:],
		[][]interface{}{
			{"[A 1]", "[B 2]", rx.Complete},
			{"[A 1]", "[B 2]", "[C 3]", rx.Complete},
			{"[A 1]", "[B 2]", "[C 3]", rx.Complete},
			{"[A 1]", "[B 2]", rx.Complete},
			{"[A 1]", "[B 2]", "[C 3]", rx.Complete},
			{"[A 1]", "[B 2]", "[C 3]", ErrTest},
		},
	)
}
