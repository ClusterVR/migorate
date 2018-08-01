package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql" // Use mysql driver

	"github.com/mizoguche/migorate/migration"
	"github.com/mizoguche/migorate/migration/db/mysql"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Usage: "migrations files directory path",
			Value: "db/migrations",
		},
		cli.StringFlag{
			Name:  "rc",
			Usage: "migrations configuration file directory path",
			Value: mysql.RCFilePath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate migration file",
			Action: func(c *cli.Context) error {
				path := c.GlobalString("path")
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					return err
				}
				return migration.Generate(path, c.Args().First(), c.Args()[1:c.NArg()])
			},
		},
		{
			Name:    "plan",
			Aliases: []string{"p"},
			Usage:   "plan migration",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "rollback, r",
					Usage: "Show rollback migration",
				},
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Show each SQL in migration",
				},
			},
			Action: func(c *cli.Context) error {
				d := migration.Up
				if c.Bool("rollback") {
					d = migration.Down
				}
				path := c.GlobalString("path")
				mysql.RCFilePath = c.GlobalString("rc")
				execFile := dest(c)
				if len(execFile) == 0 {
					log.Println("Please specify the file_name or 'all' to execute. ")
					return nil
				}
				migrations := *migration.Plan(path, d, execFile)
				count := len(migrations)
				if count == 0 {
					log.Println("No migration planned.")
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
			Name:  "exec",
			Usage: "execute migration",
			Action: func(c *cli.Context) error {
				path := c.GlobalString("path")
				mysql.RCFilePath = c.GlobalString("rc")
				execFile := dest(c)
				if len(execFile) == 0 {
					log.Println("Please specify the file_name or 'all' to execute. ")
					return nil
				}
				return migrate(path, migration.Up, execFile)
			},
		},
		{
			Name:  "rollback",
			Usage: "rollback migration",
			Action: func(c *cli.Context) error {
				path := c.GlobalString("path")
				mysql.RCFilePath = c.GlobalString("rc")
				execFile := dest(c)
				if len(execFile) == 0 {
					log.Println("Please specify the file_name or 'all' to execute. ")
					return nil
				}
				return migrate(path, migration.Down, execFile)
			},
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

func migrate(path string, d migration.Direction, dest string) error {
	migrations := *migration.Plan(path, d, dest)
	if len(migrations) == 0 {
		log.Println("No migration executed.")
		return nil
	}

	db := mysql.Database()
	defer db.Close()
	for _, m := range migrations {
		m.Exec(db, d)
	}
	return nil
}
