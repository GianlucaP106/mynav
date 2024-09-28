package app

import (
	"flag"
	"fmt"
	"log"
	"mynav/pkg"
	"os"
)

type (
	Cli struct {
		args *CliArgs
	}
	CliArgs struct {
		version *bool
		path    *string
	}
)

func newCli() *Cli {
	return &Cli{}
}

func (cli *Cli) run() {
	cli.parseArgs()
	cli.handleVersionFlag()
	cli.handlePathFlag()
}

func (cli *Cli) parseArgs() {
	version := flag.Bool("version", false, "Version of mynav")
	path := flag.String("path", ".", "Path to open mynav in")
	flag.Parse()
	cli.args = &CliArgs{
		version: version,
		path:    path,
	}
}

func (cli *Cli) handleVersionFlag() {
	if *cli.args.version {
		fmt.Println(pkg.VERSION)
		os.Exit(0)
	}
}

func (cli *Cli) handlePathFlag() {
	if cli.args.path != nil && *cli.args.path != "" {
		if err := os.Chdir(*cli.args.path); err != nil {
			log.Fatalln(err.Error())
		}
	}
}
