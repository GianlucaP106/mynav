package ui

import "mynav/pkg/core"

func SystemUpdate() bool {
	if Api().Core.IsConfigInitialized && !Api().Core.IsUpdateAsked() {
		Api().Core.SetUpdateAsked()
		update, newTag := Api().Core.DetectUpdate()
		if update {
			OpenConfirmationDialog(func(b bool) {
				if b {
					SetActionEnd(core.GetUpdateSystemCmd())
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}
