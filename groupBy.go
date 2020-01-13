package rx

import (
	"context"
)

// GroupedObservable is an Observable type used by GroupBy.
type GroupedObservable struct {
	Observable
	Key interface{}
}

type groupByObservable struct {
	Source         Observable
	KeySelector    func(interface{}) interface{}
	SubjectFactory func() Subject
}

func (obs groupByObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var groups = make(map[interface{}]Subject)
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			key := obs.KeySelector(t.Value)
			group, exists := groups[key]
			if !exists {
				group = obs.SubjectFactory()
				groups[key] = group
				sink.Next(GroupedObservable{group.Observable, key})
			}
			t.Observe(group.Observer)

		default:
			for _, group := range groups {
				t.Observe(group.Observer)
			}
			sink(t)
		}
	})
}

// GroupBy creates an Observable that groups the items emitted by the source
// Observable according to a specified criterion, and emits these grouped
// items as GroupedObservables, one GroupedObservable per group.
func (Operators) GroupBy(keySelector func(interface{}) interface{}, subjectFactory func() Subject) Operator {
	return func(source Observable) Observable {
		return groupByObservable{source, keySelector, subjectFactory}.Subscribe
	}
}
