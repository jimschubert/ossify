package new

import "gopkg.in/urfave/cli.v1"

var Project = cli.Command{
	Name:  "project",
	Usage: "Create a new project",
	Action: func(c *cli.Context) error {
		return nil
	},
}
