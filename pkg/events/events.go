package events

var emitter *EventEmitter = NewEventEmitter()

func AddEventListener(name string, listener func()) {
	emitter.Event(name).AddListener(listener)
}

func RemoveEventListener(name string, listenerId string) {
	emitter.Event(name).RemoveListener(listenerId)
}

func EmitEvent(name string) {
	emitter.Event(name).Emit()
}
