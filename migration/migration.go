package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"database/sql"
	"github.com/mizoguche/migorate/migration/db/mysql"
)

type Migration struct {
	Id   string
	Up   []string
	Down []string
}

type MigrationDirection int

const (
	Up MigrationDirection = iota
	Down
)

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

func Plan(dir string, direction MigrationDirection) *[]Migration {
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
			count := count(rows)
			if count == 0 {
				sqls = append(sqls, NewMigration(dir, id))
			}
		}
	}

	return &sqls
}

func count(r *sql.Rows) (count int) {
	for r.Next() {
		r.Scan(&count)
	}
	return count
}

func NewMigration(dir string, id string) Migration {
	b, _ := ioutil.ReadFile(dir + "/" + id + ".sql")
	r := regexp.MustCompile(`(?m)-- \+migrate Up\n([\s\S]*)\n-- \+migrate Down\n([\s\S]*)`)
	sqls := r.FindSubmatch(b)
	up := splitSql(string(sqls[1]))
	down := splitSql(string(sqls[2]))
	return Migration{Id: id, Up: up, Down: down}
}

func splitSql(src string) []string {
	raw := strings.Split(src, ";")
	sqls := make([]string, 0, len(raw))
	for _, s := range strings.Split(src, ";\n") {
		sql := strings.Replace(s, "\n", "", -1) + ";"
		sqls = append(sqls, sql)
	}
	sqls = sqls[:len(sqls)-1]
	return sqls
}
