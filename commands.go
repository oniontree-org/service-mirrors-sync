package main

import (
	"github.com/urfave/cli/v2"
	"time"
)

func (a *Application) commands() {
	a.app = &cli.App{
		Name:    "service-mirrors",
		Version: Version,
		Usage:   "Sync service mirrors",
		Commands: cli.Commands{
			&cli.Command{
				Name:      "sync",
				Usage:     "Download mirrors and update an OnionTree repository",
				ArgsUsage: "<id>",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleSyncCommand(),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "replace old URLs",
					},
					&cli.BoolFlag{
						Name:  "no-verify-signature",
						Usage: "don't verify signed message",
					},
					&cli.DurationFlag{
						Name:  "timeout",
						Usage: "request timeout",
						Value: 15 * time.Second,
					},
				},
			},
		},
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "C",
				Value: ".",
				Usage: "change directory to",
			},
		},
	}
}
