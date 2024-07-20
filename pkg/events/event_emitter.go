package events

import "github.com/google/uuid"

type (
	EventListener func()
	Event         struct {
		listeners map[string]EventListener
		name      string
	}
	EventEmitter struct {
		events map[string]*Event
	}
)

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(map[string]*Event),
	}
}

func (ee *EventEmitter) Event(name string) *Event {
	event := ee.events[name]
	if event != nil {
		return event
	}

	event = NewEvent(name)
	ee.events[name] = event

	return event
}

func NewEvent(name string) *Event {
	return &Event{
		name:      name,
		listeners: make(map[string]EventListener),
	}
}

func (e *Event) AddListener(listener func()) (listenerId string) {
	id := uuid.New().String()
	e.listeners[id] = listener
	return id
}

func (e *Event) RemoveListener(listenerId string) {
	delete(e.listeners, listenerId)
}

func (e *Event) Emit() {
	for _, l := range e.listeners {
		l()
	}
}
