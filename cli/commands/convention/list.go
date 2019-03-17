package convention

import (
	"github.com/jimschubert/ossify/config"
	"github.com/jimschubert/ossify/config/conventions"
	"gopkg.in/urfave/cli.v1"
)

var List = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "Presents a list of known conventions.",
	Action: func(c *cli.Context) error {
		conf, err := config.ConfigManager.Load()
		if err != nil {
			return err
		}
		conventions.ConventionPath = conf.ConventionPath
		conventions, err := conventions.Load()
		if err == nil {
			for _, c := range *conventions {
				err = c.Print()
				if err != nil {
					break
				}
			}
		}
		return err
	},
	Flags: []cli.Flag{
	},
}

