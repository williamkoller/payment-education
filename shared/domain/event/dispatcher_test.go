package shared_event_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	shared_event "github.com/williamkoller/system-education/shared/domain/event"
)

type fakeEvent struct {
	name string
}

func (f *fakeEvent) EventName() string {
	return f.name
}

func TestDispatcher_Register(t *testing.T) {
	d := shared_event.NewDispatcher()

	called := false
	handler := func(event interface{}) {
		called = true
	}

	d.Register("user.created", handler)

	d.Dispatch(&fakeEvent{name: "user.created"})

	time.Sleep(50 * time.Millisecond)

	assert.True(t, called, "Handler deveria ter sido chamado")
}

func TestDispatcher_Dispatch_NoHandler(t *testing.T) {
	d := shared_event.NewDispatcher()
	called := false

	handler := func(event interface{}) {
		called = true
	}

	d.Register("user.updated", handler)
	d.Dispatch(&fakeEvent{name: "user.created"})

	time.Sleep(50 * time.Millisecond)

	assert.False(t, called, "Handler não deveria ter sido chamado para evento diferente")
}

func TestDispatcher_MultipleHandlers(t *testing.T) {
	d := shared_event.NewDispatcher()
	var mu sync.Mutex
	callCount := 0

	handler1 := func(event interface{}) {
		mu.Lock()
		callCount++
		mu.Unlock()
	}
	handler2 := func(event interface{}) {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	d.Register("user.created", handler1)
	d.Register("user.created", handler2)

	d.Dispatch(&fakeEvent{name: "user.created"})

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 2, callCount, "Ambos handlers deveriam ter sido chamados")
}

func TestDispatcher_ConcurrentDispatch(t *testing.T) {
	d := shared_event.NewDispatcher()
	var mu sync.Mutex
	count := 0

	handler := func(event interface{}) {
		mu.Lock()
		count++
		mu.Unlock()
	}

	d.Register("user.created", handler)

	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Dispatch(&fakeEvent{name: "user.created"})
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 50, count, "Handler deveria ser chamado 50 vezes")
}
