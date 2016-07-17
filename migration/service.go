package migration

import (
	"database/sql"
	"fmt"
	"github.com/mizoguche/migorate/migration/db/mysql"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

// Generate migration sql file
func Generate(dir string, name string) error {
	t := time.Now()
	id := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	filepath := fmt.Sprintf("%s/%s_%s.sql", dir, id, name)
	content := []byte(`-- +migrate Up


-- +migrate Down

`)
	err := ioutil.WriteFile(filepath, content, os.ModePerm)
	if err != nil {
		log.Printf("Failed to generate file\n%v\n", err)
		return err
	}

	log.Printf("Generated: %v", filepath)
	return nil
}

// Plan migration according to migrated information in database
func Plan(dir string, direction Direction, dest string) *[]Migration {
	db := mysql.Database()
	defer db.Close()

	files, _ := ioutil.ReadDir(dir)
	r := regexp.MustCompile(`(\d\d\d\d\d\d\d\d\d\d\d\d\d\d_.+)\.sql`)
	sqls := make([]Migration, 0, len(files))
	for _, f := range files {
		if r.MatchString(f.Name()) {
			g := r.FindSubmatch([]byte(f.Name()))
			id := string(g[1])
			rows, err := db.Query("SELECT COUNT(*) FROM migorate_migrations WHERE id = ?", id)
			if err != nil {
				log.Fatalf("Failed to query: %v", err)
			}
			sqls = *add(direction, rows, dir, id, &sqls)
		}
	}

	return &sqls
}

func add(d Direction, rows *sql.Rows, dir string, id string, sqls *[]Migration) *[]Migration{
	count := count(rows)
	if d == Up {
		if count == 0 {
			s := append(*sqls, NewMigration(dir, id))
			return &s
		}
		return sqls
	}

	if count == 1 {
		s := append(*sqls, Migration{})
		copy(s[1:], s[0:])
		s[0] = NewMigration(dir, id)
		return &s
	}
	return sqls
}

func count(r *sql.Rows) (count int) {
	for r.Next() {
		r.Scan(&count)
	}
	return count
}
