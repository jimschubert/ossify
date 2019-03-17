package new

import (
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
)

func cwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func Command() cli.Command {
	return cli.Command{
		Name:     "new",
		Aliases:  []string{"n"},
		Usage:    "Create a new open-source project or other resource.",
		HelpName: "New",
		Category: "Create",
		Subcommands: []cli.Command{
			Project,
			Repository,
			License,
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "directory, d",
				Usage: "Path to the target `DIRECTORY`",
				Value: cwd(),
			},
		},
	}
}
