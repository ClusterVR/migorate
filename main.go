package main

import ()
import (
	"github.com/urfave/cli"
	"github.com/mizoguche/migorate/migration"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate migration file",
			Action: func(c *cli.Context) error {
				path := "db/migrations"
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					return err
				}
				migration.Generate(path, c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
