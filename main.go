package main

import ()
import (
	"github.com/mizoguche/migorate/migration"
	"github.com/urfave/cli"
	"os"
	"fmt"
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
				return migration.Generate(path, c.Args().First())
			},
		},
		{
			Name:    "plan",
			Aliases: []string{"p"},
			Usage:   "plan migration",
			Action: func(c *cli.Context) error {
				path := "db/migrations"
				migrations := *migration.Plan(path, migration.Up)
				count := len(migrations)
				fmt.Println("Planned migrations:")
				for i, m := range migrations {
					fmt.Printf("  %0" + fmt.Sprintf("%d", count / 10 + 1) + "d: %+v\n", (i + 1), m.ID)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
