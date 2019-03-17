package license

import (
	"errors"
	"fmt"
	"github.com/jimschubert/ossify/config"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var Add = cli.Command{
	Name:  "add",
	Usage: "Adds a new custom license (local-only) to the list of known licenses.",
	Action: func(c *cli.Context) error {
		config, err := config.ConfigManager.Load()
		if err != nil {
			return err
		}

		licensePath := config.LicensePath
		if licensePath == "" {
			return errors.New("invalid license path: please update your configuration and try again")
		}

		if err = os.MkdirAll(licensePath, 0700); err != nil {
			return err
		}

		template := c.String("template")
		if template == "" {
			template = c.Args()[0]

			if template == "" {
				return errors.New("invalid template: you must provide a template value")
			}
		}

		id := c.String("id")

		data, _ := ioutil.ReadFile(template)

		if err = os.MkdirAll(licensePath, 0700); err != nil {
			return err
		}

		targetFile := path.Join(licensePath, id)

		// TODO: Document how this allows users to specify default text for a license
		if err = ioutil.WriteFile(targetFile, data, 0644); err == nil {
			log.Println(fmt.Sprintf("Saved license with id %s to %s", id, targetFile))
			return nil
		}

		return err
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "The identifier to be associated with your customized license. This will take precedence over 'public' ids.",
		},

		cli.StringFlag{
			Name:  "template",
			Usage: "The template to add for the given identifier.",
		},
	},
}
