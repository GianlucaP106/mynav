package ui

func systemUpdate() bool {
	if getApi().LocalConfiguration.IsConfigInitialized && !getApi().GlobalConfiguration.IsUpdateAsked() {
		getApi().GlobalConfiguration.SetUpdateAsked()
		update, newTag := getApi().GlobalConfiguration.DetectUpdate()
		if update {
			openConfirmationDialog(func(b bool) {
				if b {
					runAction(func() {
						getApi().GlobalConfiguration.UpdateMynav()
					})
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}
