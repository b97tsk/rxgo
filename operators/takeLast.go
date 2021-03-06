package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// TakeLast emits only the last count values emitted by the source.
func TakeLast(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count <= 0 {
			return rx.Empty()
		}

		return takeLastObservable{source, count}.Subscribe
	}
}

type takeLastObservable struct {
	Source rx.Observable
	Count  int
}

func (obs takeLastObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var queue queue.Queue

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			if queue.Len() == obs.Count {
				queue.Pop()
			}

			queue.Push(t.Value)

		case t.HasError:
			sink(t)

		default:
			for i, j := 0, queue.Len(); i < j; i++ {
				if ctx.Err() != nil {
					return
				}

				sink.Next(queue.At(i))
			}

			sink(t)
		}
	})
}
