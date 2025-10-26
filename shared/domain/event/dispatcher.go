package shared_event

import "sync"

type Handler func(event interface{})

type Dispatcher struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[string][]Handler)}
}

func (d *Dispatcher) Register(eventName string, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *Dispatcher) Dispatch(event interface{}) {
	if event == nil {
		return
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	ev, ok := event.(Event)
	if !ok {
		return
	}

	handlers, ok := d.handlers[ev.EventName()]
	if !ok || len(handlers) == 0 {
		return
	}

	for _, h := range handlers {
		go func(handler Handler, e interface{}) {
			defer func() { recover() }()
			handler(e)
		}(h, event)
	}

}

func (d *Dispatcher) DispatchSync(event interface{}) {
	if event == nil {
		return
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	ev, ok := event.(Event)
	if !ok {
		return
	}

	handlers, ok := d.handlers[ev.EventName()]
	if !ok || len(handlers) == 0 {
		return
	}

	for _, h := range handlers {
		func(handler Handler, e interface{}) {
			defer func() { recover() }()
			handler(e)
		}(h, event)
	}

}
