package ui

func SystemUpdate() bool {
	if Api().Configuration.IsConfigInitialized && !Api().Configuration.IsUpdateAsked() {
		Api().Configuration.SetUpdateAsked()
		update, newTag := Api().Configuration.DetectUpdate()
		if update {
			OpenConfirmationDialog(func(b bool) {
				if b {
					RunAction(func() {
						Api().Configuration.UpdateMynav()
					})
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}
