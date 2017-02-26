package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// Migration information
type Migration struct {
	ID   string
	Up   []string
	Down []string
}

// Direction (Up or Down)
type Direction int

const (
	// Up used on migration to forward
	Up Direction = iota

	// Down used on rollback
	Down
)

// NewMigration constructor
func NewMigration(dir string, id string) Migration {
	b, _ := ioutil.ReadFile(dir + "/" + id + ".sql")
	r := regexp.MustCompile(`(?m)-- \+migrate Up\n([\s\S]*)\n-- \+migrate Down\n([\s\S]*)`)
	sqls := r.FindSubmatch(b)
	up := splitSQL(string(sqls[1]))
	down := splitSQL(string(sqls[2]))
	return Migration{ID: id, Up: up, Down: down}
}

func splitSQL(src string) []string {
	sqls := []string{}
	for _, s := range strings.Split(src, ";\n") {
		// replace sql comment
		if strings.HasPrefix(s, "\n--") {
			exp := regexp.MustCompile(`\n--.*\n`)
			s = exp.ReplaceAllString(s, "")
		}
		if len(s) == 0 || s == "\n" {
			continue
		}
		sqls = append(sqls, s+";")
	}
	return sqls
}

// Exec migration
func (m *Migration) Exec(db *sql.DB, d Direction) {
	var sql []string
	if d == Up {
		sql = m.Up
	} else {
		sql = m.Down
	}

	log.Printf("Executing %v", m.ID)
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v\n", err)
	}

	for _, s := range sql {
		_, err = db.Exec(s)
		failIfError(s, err, tx)
		log.Println("Executed:")
		fmt.Printf("%s\n\n", s)
	}

	var migSQL string
	if d == Up {
		migSQL = "INSERT INTO migorate_migrations(id, migrated_at) VALUES(?, NOW())"
	} else {
		migSQL = "DELETE FROM migorate_migrations WHERE id = ?"
	}
	_, err = db.Exec(migSQL, m.ID)
	failIfError(migSQL, err, tx)

	tx.Commit()
}

func failIfError(s string, err error, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to execute SQL: %v\n%v\nRollback transaction.", s, err)
	}
}
