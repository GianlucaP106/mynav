package app

import (
	"mynav/pkg/ui"
)

func Main() {
	for {
		action := ui.Start()
		if action == nil {
			break
		}
		handleAction(action)
	}
}
