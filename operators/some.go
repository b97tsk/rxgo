package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Some emits whether or not any item of the source satisfies a specified
// predicate.
//
// Some emits true or false, then completes.
func Some(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return someObservable{source, predicate}.Subscribe
	}
}

type someObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs someObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if obs.Predicate(t.Value, sourceIndex) {
				observer = rx.Noop

				sink.Next(true)
				sink.Complete()
			}

		case t.HasError:
			sink(t)

		default:
			sink.Next(false)
			sink.Complete()
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}
