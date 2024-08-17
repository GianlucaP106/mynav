package ui

import "mynav/pkg/tui"

func systemUpdate() bool {
	if getApi().Configuration.IsConfigInitialized && !getApi().Configuration.IsUpdateAsked() {
		getApi().Configuration.SetUpdateAsked()
		update, newTag := getApi().Configuration.DetectUpdate()
		if update {
			openConfirmationDialog(func(b bool) {
				if b {
					tui.RunAction(func() {
						getApi().Configuration.UpdateMynav()
					})
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}
