package shared_event

import "time"

type Event interface {
	EventName() string
	OccurredOn() time.Time
}

type AggregateRoot struct {
	domainEvents []Event
}

func (a *AggregateRoot) AddDomainEvent(e Event) {
	a.domainEvents = append(a.domainEvents, e)
}

func (a *AggregateRoot) PullDomainEvents() []Event {
	if a == nil {
		return nil
	}
	events := a.domainEvents
	a.domainEvents = []Event{}
	return events
}
