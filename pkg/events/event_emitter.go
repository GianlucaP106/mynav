package events

import (
	"sync"

	"github.com/google/uuid"
)

type (
	EventListener func(listenerId string)
	Event         struct {
		listeners map[string]EventListener
		mu        *sync.RWMutex
		name      string
	}
	EventEmitter struct {
		events map[string]*Event
		mu     *sync.RWMutex
	}
)

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(map[string]*Event),
		mu:     &sync.RWMutex{},
	}
}

func (ee *EventEmitter) Event(name string) *Event {
	ee.mu.RLock()
	defer ee.mu.RUnlock()
	event := ee.events[name]
	if event != nil {
		return event
	}

	event = NewEvent(name, ee.mu)
	ee.events[name] = event

	return event
}

func NewEvent(name string, mu *sync.RWMutex) *Event {
	return &Event{
		name:      name,
		listeners: make(map[string]EventListener),
		mu:        mu,
	}
}

func (e *Event) AddListener(listener func(string)) (listenerId string) {
	id := uuid.New().String()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[id] = listener
	return id
}

func (e *Event) RemoveListener(listenerId string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.listeners, listenerId)
}

func (e *Event) Emit() {
	e.mu.RLock()
	listeners := make(map[string]EventListener)
	for id, el := range e.listeners {
		listeners[id] = el
	}
	e.mu.RUnlock()

	for id, l := range listeners {
		l(id)
	}
}
