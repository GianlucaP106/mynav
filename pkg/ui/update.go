package ui

import "mynav/pkg/tui"

func systemUpdate() bool {
	if getApi().LocalConfiguration.IsConfigInitialized && !getApi().GlobalConfiguration.IsUpdateAsked() {
		getApi().GlobalConfiguration.SetUpdateAsked()
		update, newTag := getApi().GlobalConfiguration.DetectUpdate()
		if update {
			openConfirmationDialog(func(b bool) {
				if b {
					tui.RunAction(func() {
						getApi().GlobalConfiguration.UpdateMynav()
					})
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}
