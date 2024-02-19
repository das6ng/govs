package main

import (
	"github.com/urfave/cli/v2"
)

var goRootDir string

var app = &cli.App{
	Usage: "Golang version switcher.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "root",
			Usage:       "go sdk parent dir",
			Destination: &goRootDir,
			Aliases:     []string{"p"},
			EnvVars:     []string{"GO_ROOT_PARENT_DIR"},
			Value:       currGoRoot(),
		},
	},
	Action: runSwitcher,
}
