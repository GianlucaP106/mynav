package app

import (
	"flag"
	"fmt"
	"mynav/pkg/api"
	"mynav/pkg/ui"
)

func Main() {
	version := flag.Bool("version", false, "Version of mynav")
	flag.Parse()

	if *version {
		fmt.Println(api.VERSION)
		return
	}

	for {
		action := ui.Start()
		if action.Command == nil {
			break
		}
		handleAction(action)
	}
}
