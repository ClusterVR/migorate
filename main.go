package main

import (
	"fmt"
	"github.com/mizoguche/migorate/migration"
	"github.com/urfave/cli"
	"os"
	"log"
	"github.com/mizoguche/migorate/migration/db/mysql"
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
				migrations := *migration.Plan(path, migration.Up, dest(c))
				count := len(migrations)
				if count == 0 {
					log.Printf("No migration planned.")
					return nil
				}

				log.Println("Planned migrations:")
				for i, m := range migrations {
					log.Printf("  %0"+fmt.Sprintf("%d", count/10+1)+"d: %+v\n", (i + 1), m.ID)
				}
				return nil
			},
		},
		{
			Name:    "exec",
			Usage:   "execute migration",
			Action: func(c *cli.Context) error { return migrate(migration.Up, dest(c)) },
		},
		{
			Name:    "rollback",
			Usage:   "rollback migration",
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