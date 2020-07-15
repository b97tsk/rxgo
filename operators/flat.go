package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type flatObservable struct {
	Source  rx.Observable
	Flat    func(observables ...rx.Observable) rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs flatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	var (
		observables []rx.Observable
		observer    rx.Observer
	)

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			obs1 := obs.Project(t.Value, sourceIndex)
			observables = append(observables, obs1)

		case t.HasError:
			sink(t)

		default:
			obs := obs.Flat(observables...)
			obs.Subscribe(ctx, sink)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// FlatAll creates an Observable that flattens a higher-order Observable
// into a first-order Observable, by applying a flat function to the inner
// Observables, and starts subscribing to it.
func FlatAll(flat func(observables ...rx.Observable) rx.Observable) rx.Operator {
	return FlatMap(flat, projectToObservable)
}

// FlatMap creates an Observable that converts the source Observable into a
// higher-order Observable, by projecting each source value to an Observable,
// and flattens it into a first-order Observable, by applying a flat function
// to the inner Observables, and starts subscribing to it.
func FlatMap(
	flat func(observables ...rx.Observable) rx.Observable,
	project func(interface{}, int) rx.Observable,
) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return flatObservable{source, flat, project}.Subscribe
	}
}
