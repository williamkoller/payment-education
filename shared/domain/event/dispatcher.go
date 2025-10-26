package shared_event

import "sync"

type Handler func(event interface{})

type Dispatcher struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]Handler),
	}
}

func (d *Dispatcher) Register(eventName string, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *Dispatcher) Dispatch(event interface{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	switch e := event.(type) {
	case interface{ EventName() string }:
		if handlers, ok := d.handlers[e.EventName()]; ok {
			for _, h := range handlers {
				go h(e)
			}
		}
	}
}
