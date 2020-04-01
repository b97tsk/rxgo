package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/queue"
)

// An ExpandConfigure is a configure for Expand.
type ExpandConfigure struct {
	Project    func(interface{}, int) rx.Observable
	Concurrent int
}

// Use creates an Operator from this configure.
func (configure ExpandConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := expandObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type expandObservable struct {
	Source rx.Observable
	ExpandConfigure
}

func (obs expandObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	type X struct {
		Index           int
		ActiveCount     int
		SourceCompleted bool
		Buffer          queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		sourceIndex := x.Index
		sourceValue := x.Buffer.PopFront()
		x.Index++

		sink.Next(sourceValue)

		// calls obs.Project synchronously
		obs1 := obs.Project(sourceValue, sourceIndex)

		go obs1.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				x := <-cx
				x.Buffer.PushBack(t.Value)
				if x.ActiveCount != obs.Concurrent {
					x.ActiveCount++
					doNextLocked(x)
				}
				cx <- x

			case t.HasError:
				sink(t)

			default:
				x := <-cx
				if x.Buffer.Len() > 0 {
					doNextLocked(x)
				} else {
					x.ActiveCount--
					if x.ActiveCount == 0 && x.SourceCompleted {
						sink(t)
					}
				}
				cx <- x
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.ActiveCount != obs.Concurrent {
				x.ActiveCount++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.ActiveCount == 0 {
				sink(t)
			}
			cx <- x
		}
	})
}

// Expand recursively projects each source value to an Observable which is
// merged in the output Observable.
//
// It's similar to MergeMap, but applies the projection function to every
// source value as well as every output value. It's recursive.
func Expand(project func(interface{}, int) rx.Observable) rx.Operator {
	return ExpandConfigure{project, -1}.Use()
}