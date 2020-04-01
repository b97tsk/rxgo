package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type findObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs findObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		sourceIndex = -1
		observer    rx.Observer
	)

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if obs.Predicate(t.Value, sourceIndex) {
				observer = rx.NopObserver
				sink(t)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func Find(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := findObservable{source, predicate}
		return rx.Create(obs.Subscribe)
	}
}