package port_event

import shared_event "github.com/williamkoller/system-education/shared/domain/event"

type Dispacther interface {
	Dispatch(event interface{})
	Register(eventName string, handler shared_event.Handler)
}
