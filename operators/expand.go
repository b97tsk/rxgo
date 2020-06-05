package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/queue"
)

// An ExpandConfigure is a configure for Expand.
type ExpandConfigure struct {
	Project    func(interface{}, int) (rx.Observable, error)
	Concurrent int
}

// Use creates an Operator from this configure.
func (configure ExpandConfigure) Use() rx.Operator {
	if configure.Project == nil {
		panic("Expand: nil Project")
	}
	if configure.Concurrent == 0 {
		configure.Concurrent = -1
	}
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
		Active          int
		Buffer          queue.Queue
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		sourceIndex := x.Index
		sourceValue := x.Buffer.PopFront()
		x.Index++

		sink.Next(sourceValue)

		obs1, err := obs.Project(sourceValue, sourceIndex)
		if err != nil {
			sink.Error(err)
			return
		}

		go obs1.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				x := <-cx
				x.Buffer.PushBack(t.Value)
				if x.Active != obs.Concurrent {
					x.Active++
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
					x.Active--
					if x.Active == 0 && x.SourceCompleted {
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
			if x.Active != obs.Concurrent {
				x.Active++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.Active == 0 {
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
func Expand(project func(interface{}, int) (rx.Observable, error)) rx.Operator {
	return ExpandConfigure{project, -1}.Use()
}
