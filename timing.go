package rx

import (
	"context"
	"time"
)

// Ticker creates an Observable that emits time.Time values every specified
// interval of time.
func Ticker(duration time.Duration) Observable {
	return Create(
		func(ctx context.Context, sink Observer) {
			go func() {
				ticker := time.NewTicker(duration)
				defer ticker.Stop()
				done := ctx.Done()
				for {
					select {
					case <-done:
						return
					case t := <-ticker.C:
						sink.Next(t)
					}
				}
			}()
		},
	)
}

// Timer creates an Observable that emits only a time.Time value after
// a particular time span has passed.
func Timer(duration time.Duration) Observable {
	return Create(
		func(ctx context.Context, sink Observer) {
			go func() {
				timer := time.NewTimer(duration)
				defer timer.Stop()
				select {
				case <-ctx.Done():
				case t := <-timer.C:
					sink.Next(t)
					sink.Complete()
				}
			}()
		},
	)
}
