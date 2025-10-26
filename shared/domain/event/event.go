package shared_event

import "time"

type Event interface {
	EventName() string
	OccurredOn() time.Time
}