package app

import (
	"flag"
	"fmt"
	"mynav/pkg/ui"
)

func Main() {
	version := flag.Bool("version", false, "Version of mynav")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	for {
		action := ui.Start()
		if action == nil || action.Command == nil {
			break
		}

		handleAction(action)
	}
}
