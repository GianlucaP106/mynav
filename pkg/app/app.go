package app

import (
	"flag"
	"fmt"
	"mynav/pkg"
	"mynav/pkg/ui"
	"os"
)

func Main() {
	version := flag.Bool("version", false, "Version of mynav")
	path := flag.String("path", ".", "Path to open mynav in")
	flag.Parse()

	if *version {
		fmt.Println(pkg.VERSION)
		return
	}

	if path != nil && *path != "" {
		if err := os.Chdir(*path); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	for {
		action := ui.Start()
		if action == nil || action.Command == nil {
			break
		}

		handleAction(action)
	}
}
