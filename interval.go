package rx

import (
	"context"
	"time"
)

type intervalObservable struct {
	InitialDelay time.Duration
	Period       time.Duration
	HasPeriod    bool
}

func (obs intervalObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	if obs.HasPeriod {
		if obs.InitialDelay != obs.Period {
			scheduleOnce(ctx, obs.InitialDelay, func() {
				index := 0
				wait := make(chan struct{})
				schedule(ctx, obs.Period, func() {
					<-wait
					sink.Next(index)
					index++
				})
				sink.Next(index)
				index++
				close(wait)
			})
		} else {
			index := 0
			schedule(ctx, obs.Period, func() {
				sink.Next(index)
				index++
			})
		}
	} else {
		scheduleOnce(ctx, obs.InitialDelay, func() {
			sink.Next(0)
			sink.Complete()
		})
	}

	return ctx, cancel
}

// Interval creates an Observable that emits sequential integers every
// specified interval of time.
func Interval(period time.Duration) Observable {
	return intervalObservable{period, period, true}.Subscribe
}

// Timer creates an Observable that starts emitting after an initialDelay and
// emits ever increasing integers after each period of time thereafter.
//
// Its like Interval, but you can specify when should the emissions start.
//
// If period is not specified, the output Observable emits only one value,
// zero. Otherwise, it emits an infinite sequence.
func Timer(initialDelay time.Duration, period ...time.Duration) Observable {
	obs := intervalObservable{InitialDelay: initialDelay}
	switch len(period) {
	case 0:
	case 1:
		obs.Period = period[0]
		obs.HasPeriod = true
	default:
		panic("too many parameters")
	}
	return obs.Subscribe
}
