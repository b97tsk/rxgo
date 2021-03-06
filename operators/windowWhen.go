package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
	"github.com/b97tsk/rx/internal/norec"
)

// WindowWhen branches out the source values as a nested Observable using
// a factory function of closing Observables to determine when to start
// a new window.
//
// It's like BufferWhen, but emits a nested Observable instead of a slice.
func WindowWhen(closingSelector func() rx.Observable) rx.Operator {
	return WindowWhenConfigure{ClosingSelector: closingSelector}.Make()
}

// A WindowWhenConfigure is a configure for WindowWhen.
type WindowWhenConfigure struct {
	ClosingSelector func() rx.Observable
	WindowFactory   rx.DoubleFactory
}

// Make creates an Operator from this configure.
func (configure WindowWhenConfigure) Make() rx.Operator {
	if configure.ClosingSelector == nil {
		panic("WindowWhen: ClosingSelector is nil")
	}

	if configure.WindowFactory == nil {
		configure.WindowFactory = rx.Multicast
	}

	return func(source rx.Observable) rx.Observable {
		return windowWhenObservable{source, configure}.Subscribe
	}
}

type windowWhenObservable struct {
	Source rx.Observable
	WindowWhenConfigure
}

func (obs windowWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Window rx.Observer
	}

	window := obs.WindowFactory()
	x.Window = window.Observer
	sink.Next(window.Observable)

	var openWindow func()

	openWindow = norec.Wrap(func() {
		if ctx.Err() != nil {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			observer = rx.Noop

			cancel()

			if critical.Enter(&x.Section) {
				switch {
				case t.HasValue:
					x.Window.Complete()

					window := obs.WindowFactory()
					x.Window = window.Observer
					sink.Next(window.Observable)

					critical.Leave(&x.Section)

					openWindow()

				case t.HasError:
					critical.Close(&x.Section)

					x.Window.Sink(t)
					sink(t)

				default:
					critical.Leave(&x.Section)
				}
			}
		}

		closingNotifier := obs.ClosingSelector()

		closingNotifier.Subscribe(ctx, observer.Sink)
	})

	openWindow()

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Window.Sink(t)

				critical.Leave(&x.Section)

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				x.Window.Sink(t)
				sink(t)
			}
		}
	})
}
