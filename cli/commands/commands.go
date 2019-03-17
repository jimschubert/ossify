package commands

import (
	"github.com/jimschubert/ossify/cli/commands/convention"
	"github.com/jimschubert/ossify/cli/commands/license"
	"github.com/jimschubert/ossify/cli/commands/new"
	"gopkg.in/urfave/cli.v1"
)

func All() []cli.Command {
	return []cli.Command{
		new.Command(),
		license.Command(),
		convention.Command(),
	}
}
