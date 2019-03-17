package new

import (
	"gopkg.in/urfave/cli.v1"
)

// This command aims to differ form the 'license' command in that we can evaluate project-level licenses and perform a validation here.
// That target may prove to be too difficult, and may only ever end up being recommendation-only.
var License = cli.Command{
	Name:  "license",
	Usage: "Create or evaluate a new license for your repository.",
	Action: func(c *cli.Context) error {

		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "The license type (any OSI id under SPDX schema). See license --list for details.",
			Value: cwd(),
		},
		cli.StringFlag{
			Name:  "list",
			Usage: "The license type (any OSI id under SPDX schema). See --list for details.",
			Value: cwd(),
		},
	},
}
