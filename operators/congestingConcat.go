package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// CongestingConcatAll creates an Observable that flattens a higher-order
// Observable into a first-order Observable by concatenating the inner
// Observables in order.
//
// It's like ConcatAll, but it congests the source.
func CongestingConcatAll() rx.Operator {
	return CongestingConcatMap(projectToObservable)
}

// CongestingConcatMap creates an Observable that converts the source
// Observable into a higher-order Observable, by projecting each source
// value to an Observable, and flattens it into a first-order Observable
// by concatenating the inner Observables in order.
//
// It's like ConcatMap, but it congests the source.
func CongestingConcatMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return congestingConcatObservable{source, project}.Subscribe
	}
}

// CongestingConcatMapTo creates an Observable that converts the source
// Observable into a higher-order Observable, by projecting each source
// value to the same Observable, and flattens it into a first-order
// Observable by concatenating the inner Observables in order.
//
// It's like CongestingConcatMap, but maps each value always to the same
// inner Observable.
//
// It's like ConcatMapTo, but it congests the source.
func CongestingConcatMapTo(inner rx.Observable) rx.Operator {
	return CongestingConcatMap(func(interface{}, int) rx.Observable { return inner })
}

type congestingConcatObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs congestingConcatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			obs1 := obs.Project(t.Value, sourceIndex)
			err := obs1.BlockingSubscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
				}
			})
			if err != nil {
				observer = rx.Noop
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}
