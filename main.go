package main

import (
	"github.com/jimschubert/ossify/cli/commands"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ossify"
	app.Usage = "Give some structure to your open-source software projects."

	app.EnableBashCompletion = true
	app.Commands = commands.All()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
