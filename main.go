package main

import (
	"fmt"
	"github.com/mizoguche/migorate/migration"
	"github.com/mizoguche/migorate/migration/db/mysql"
	"github.com/urfave/cli"
	"log"
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
				return migration.Generate(path, c.Args().First())
			},
		},
		{
			Name:    "plan",
			Aliases: []string{"p"},
			Usage:   "plan migration",
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "rollback, r",
					Usage: "Show rollback migration",
				},
				cli.BoolFlag{
					Name: "verbose, v",
					Usage: "Show each SQL in migration",
				},
			},
			Action: func(c *cli.Context) error {
				d := migration.Up
				if c.Bool("rollback") {
					d = migration.Down
				}
				path := "db/migrations"
				migrations := *migration.Plan(path, d, dest(c))
				count := len(migrations)
				if count == 0 {
					log.Printf("No migration planned.")
					return nil
				}

				log.Println("Planned migrations:")
				for i, m := range migrations {
					log.Printf("  %0"+fmt.Sprintf("%d", count/10+1)+"d: %+v\n", (i + 1), m.ID)
					if c.Bool("verbose") && d == migration.Up {
						for _, s := range m.Up {
							log.Printf("        %+v\n", s)
						}
					}
					if c.Bool("verbose") && d == migration.Down {
						for _, s := range m.Down {
							log.Printf("        %+v\n", s)
						}
					}
				}
				return nil
			},
		},
		{
			Name:   "exec",
			Usage:  "execute migration",
			Action: func(c *cli.Context) error { return migrate(migration.Up, dest(c)) },
		},
		{
			Name:   "rollback",
			Usage:  "rollback migration",
			Action: func(c *cli.Context) error { return migrate(migration.Down, dest(c)) },
		},
	}

	app.Run(os.Args)
}

func dest(c *cli.Context) string {
	if c.NArg() > 0 {
		return c.Args()[0]
	}
	return ""
}

func migrate(d migration.Direction, dest string) error {
	path := "db/migrations"
	migrations := *migration.Plan(path, d, dest)
	if len(migrations) == 0 {
		log.Printf("No migration executed.")
		return nil
	}

	db := mysql.Database()
	defer db.Close()
	for _, m := range migrations {
		m.Exec(db, d)
	}
	return nil
}
