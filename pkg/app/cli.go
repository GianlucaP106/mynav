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
		Args *CliArgs
	}
	CliArgs struct {
		Version *bool
		Path    *string
	}
)

func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) Run() {
	cli.ParseArgs()
	cli.handleVersionFlag()
	cli.handlePathFlag()
	NewApp().Start()
}

func (cli *Cli) ParseArgs() {
	version := flag.Bool("version", false, "Version of mynav")
	path := flag.String("path", ".", "Path to open mynav in")
	flag.Parse()
	cli.Args = &CliArgs{
		Version: version,
		Path:    path,
	}
}

func (cli *Cli) handleVersionFlag() {
	if *cli.Args.Version {
		fmt.Println(pkg.VERSION)
		os.Exit(0)
	}
}

func (cli *Cli) handlePathFlag() {
	if cli.Args.Path != nil && *cli.Args.Path != "" {
		if err := os.Chdir(*cli.Args.Path); err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func Main() {
	NewCli().Run()
}
