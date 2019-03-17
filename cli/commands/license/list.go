package license

import (
	licenseUtil "github.com/jimschubert/ossify/licenses"
	"gopkg.in/urfave/cli.v1"
)

var List = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "Presents a list of known licenses.",
	Action: func(c *cli.Context) error {
		// TODO: load custom licenses from file system and merge with "official" list?
		licenses, err := licenseUtil.LoadLicenses()
		if err != nil {
			return err
		}
		for _, license := range *licenses {
			_ = license.Print()
		}
		return nil
	},
	Flags: []cli.Flag{
	},
}
