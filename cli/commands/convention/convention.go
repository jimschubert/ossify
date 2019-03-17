package convention

import (
	"gopkg.in/urfave/cli.v1"
)

func Command() cli.Command {
	return cli.Command{
		Name: "convention",
		Usage: "Manage structure/file convention.",
		Category: "Manage",
		Subcommands: []cli.Command{
			List,
			Add,
		},
		Action: func(c *cli.Context) error {
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "name",
				Usage: "The name of the convention, to be used by commands supporting conventional naming, generation, or validation.",
			},
		},
	}
}
