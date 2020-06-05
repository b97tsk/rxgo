package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/misc"
	"github.com/b97tsk/rx/x/queue"
)

type concatObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) (rx.Observable, error)
}

func (obs concatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	type X struct {
		Index  int
		Active int
		Buffer queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{Active: 1}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		var avoidRecursion misc.AvoidRecursion
		avoidRecursion.Do(func() {
			if x.Buffer.Len() == 0 {
				x.Active--
				if x.Active == 0 {
					sink.Complete()
				}
				return
			}

			sourceIndex := x.Index
			sourceValue := x.Buffer.PopFront()
			x.Index++

			obs, err := obs.Project(sourceValue, sourceIndex)
			if err != nil {
				sink.Error(err)
				return
			}

			obs.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
					return
				}
				if ctx.Err() != nil {
					return
				}
				avoidRecursion.Do(func() {
					x := <-cx
					doNextLocked(x)
					cx <- x
				})
			})
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.Active == 1 {
				x.Active++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.Active--
			if x.Active == 0 {
				sink(t)
			}
			cx <- x
		}
	})
}

// ConcatAll converts a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
//
// ConcatAll flattens an Observable-of-Observables by putting one inner
// Observable after the other.
func ConcatAll() rx.Operator {
	return ConcatMap(projectToObservable)
}

// ConcatMap projects each source value to an Observable which is merged in
// the output Observable, in a serialized fashion waiting for each one to
// complete before merging the next.
//
// ConcatMap maps each value to an Observable, then flattens all of these inner
// Observables using ConcatAll.
func ConcatMap(project func(interface{}, int) (rx.Observable, error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := concatObservable{source, project}
		return rx.Create(obs.Subscribe)
	}
}

// ConcatMapTo projects each source value to the same Observable which is
// merged multiple times in a serialized fashion on the output Observable.
//
// It's like ConcatMap, but maps each value always to the same inner Observable.
func ConcatMapTo(inner rx.Observable) rx.Operator {
	return ConcatMap(func(interface{}, int) (rx.Observable, error) { return inner, nil })
}
