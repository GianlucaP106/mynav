package events

import (
	"mynav/pkg/tasks"
)

var emitter *EventEmitter = NewEventEmitter()

func AddEventListener(name string, listener func(string)) {
	emitter.Event(name).AddListener(listener)
}

func RemoveEventListener(name string, listenerId string) {
	emitter.Event(name).RemoveListener(listenerId)
}

func Emit(name string) {
	tasks.QueueTask(func() {
		emitter.Event(name).Emit()
	})
}

const (
	WorkspaceChangeEvent           = "WorkspaceChange"
	TmuxSessionChangeEvent         = "TmuxSessionChangeEvent"
	TmuxWindowChangeEvent          = "TmuxWindowChangeEvent"
	TmuxPreviewChangeEvent         = "TmuxPreviewChangeEvent"
	TmuxPaneChangeEvent            = "TmuxPaneChangeEvent"
	TopicChangeEvent               = "TopicChangeEvent"
	PortChangeEvent                = "PortChangeEvent"
	GithubPrsChangesEvent          = "GithubPrsChangesEvent"
	GithubReposChangesEvent        = "GithubReposChangesEvent"
	GithubDeviceAuthenticatedEvent = "GithubDeviceAuthenticatedEvent"
)
