package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// DoAfter mirrors the source, but performs a side effect after each emission.
func DoAfter(tap rx.Observer) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				defer tap(t)
				sink(t)
			})
		}
	}
}

// DoAfterNext mirrors the source, but performs a side effect after each value.
func DoAfterNext(f func(interface{})) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					defer f(t.Value)
				}

				sink(t)
			})
		}
	}
}

// DoAfterError mirrors the source and, when the source throws an error,
// performs a side effect after mirroring this error.
func DoAfterError(f func(error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if t.HasError {
					defer f(t.Error)
				}

				sink(t)
			})
		}
	}
}

// DoAfterComplete mirrors the source and, when the source completes, performs
// a side effect after mirroring this completion.
func DoAfterComplete(f func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if !t.HasValue && !t.HasError {
					defer f()
				}

				sink(t)
			})
		}
	}
}

// DoAfterErrorOrComplete mirrors the source and, when the source throws an
// error or completes, performs a side effect after mirroring this error or
// completion.
func DoAfterErrorOrComplete(f func()) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				if !t.HasValue {
					defer f()
				}

				sink(t)
			})
		}
	}
}
