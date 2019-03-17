package new

import "gopkg.in/urfave/cli.v1"

var Repository = cli.Command{
	Name:    "repository",
	Aliases: []string{"repo"},
	Usage:   "Create a new remote repository.",
	Action: func(c *cli.Context) error {
		return nil
	},
}
