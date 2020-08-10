package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// A ThrottleConfigure is a configure for Throttle.
type ThrottleConfigure struct {
	DurationSelector func(interface{}) rx.Observable
	Leading          bool
	Trailing         bool
}

// Make creates an Operator from this configure.
func (configure ThrottleConfigure) Make() rx.Operator {
	if configure.DurationSelector == nil {
		panic("Throttle: DurationSelector is nil")
	}
	return func(source rx.Observable) rx.Observable {
		return throttleObservable{source, configure}.Subscribe
	}
}

// A ThrottleTimeConfigure is a configure for ThrottleTime.
type ThrottleTimeConfigure struct {
	Duration time.Duration
	Leading  bool
	Trailing bool
}

// Make creates an Operator from this configure.
func (configure ThrottleTimeConfigure) Make() rx.Operator {
	obsTimer := rx.Timer(configure.Duration)
	durationSelector := func(interface{}) rx.Observable { return obsTimer }
	return ThrottleConfigure{
		DurationSelector: durationSelector,
		Leading:          configure.Leading,
		Trailing:         configure.Trailing,
	}.Make()
}

type throttleObservable struct {
	Source rx.Observable
	ThrottleConfigure
}

func (obs throttleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Trailing struct {
			Value    interface{}
			HasValue bool
		}
	}

	var (
		doThrottle  func(interface{})
		throttleCtx context.Context
	)

	doThrottle = func(val interface{}) {
		ctx, cancel := context.WithCancel(ctx)
		throttleCtx = ctx

		var observer rx.Observer
		observer = func(t rx.Notification) {
			observer = rx.Noop
			defer cancel()
			if obs.Trailing || t.HasError {
				if critical.Enter(&x.Section) {
					if t.HasError {
						critical.Close(&x.Section)
						sink(t)
						return
					}
					if x.Trailing.HasValue {
						sink.Next(x.Trailing.Value)
						x.Trailing.HasValue = false
						doThrottle(x.Trailing.Value)
					}
					critical.Leave(&x.Section)
				}
			}
		}

		obs := obs.DurationSelector(val)
		go obs.Subscribe(throttleCtx, observer.Sink)
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Trailing.Value = t.Value
				x.Trailing.HasValue = true
				if throttleCtx == nil || throttleCtx.Err() != nil {
					doThrottle(t.Value)
					if obs.Leading {
						sink(t)
						x.Trailing.HasValue = false
					}
				}
				critical.Leave(&x.Section)
			default:
				critical.Close(&x.Section)
				if x.Trailing.HasValue {
					sink.Next(x.Trailing.Value)
				}
				sink(t)
			}
		}
	})
}

// Throttle creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration determined
// by another Observable, then repeats this process.
//
// It's like ThrottleTime, but the silencing duration is determined by a second
// Observable.
func Throttle(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return ThrottleConfigure{
		DurationSelector: durationSelector,
		Leading:          true,
		Trailing:         false,
	}.Make()
}

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func ThrottleTime(d time.Duration) rx.Operator {
	return ThrottleTimeConfigure{
		Duration: d,
		Leading:  true,
		Trailing: false,
	}.Make()
}
