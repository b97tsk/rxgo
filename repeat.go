package rx

import (
	"context"
)

type repeatObservable struct {
	Source Observable
	Count  int
}

func (obs repeatObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		count          = obs.Count
		observer       Observer
		avoidRecursive avoidRecursiveCalls
	)

	subscribe := func() {
		obs.Source.Subscribe(ctx, observer)
	}

	observer = func(t Notification) {
		switch {
		case t.HasValue || t.HasError:
			sink(t)
		default:
			if count == 0 {
				sink(t)
			} else {
				if count > 0 {
					count--
				}
				avoidRecursive.Do(subscribe)
			}
		}
	}

	avoidRecursive.Do(subscribe)

	return ctx, cancel
}

// Repeat creates an Observable that repeats the stream of items emitted by the
// source Observable at most count times.
func (Operators) Repeat(count int) OperatorFunc {
	return func(source Observable) Observable {
		if count == 0 {
			return Empty()
		}
		if count == 1 {
			return source
		}
		if count > 0 {
			count--
		}
		return repeatObservable{source, count}.Subscribe
	}
}
