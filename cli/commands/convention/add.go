package convention

import (
	"gopkg.in/urfave/cli.v1"
)

var Add = cli.Command{
	Name:  "add",
	Usage: "Adds a new custom convention (local-only) to the list of known conventions.",
	Action: func(c *cli.Context) error {
		//config, err := config.LoadConfig()
		//if err != nil {
		//	return err
		//}
		//
		//// TODO: create config.ConventionsPath
		//
		//conventionsPath := config.ConventionsPath
		//if conventionsPath == "" {
		//	return errors.New("invalid conventions path: please update your configuration and try again")
		//}
		//
		//if err = os.MkdirAll(conventionsPath, 0700); err != nil {
		//	return err
		//}

		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "The identifier to be associated with your customized convention. This will take precedence over a built-in convention with the same id.",
		},
	},
}