package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestCombineLatest(t *testing.T) {
	subscribe(
		t,
		CombineLatest(
			Just("A", "B").Pipe(addLatencyToValue(3, 5)),
			Just("C", "D").Pipe(addLatencyToValue(2, 4)),
			Just("E", "F").Pipe(addLatencyToValue(1, 3)),
		).Pipe(toString),
		"[A C E]", "[A C F]", "[A D F]", "[B D F]", Complete,
	)
}
