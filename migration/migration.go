package migration

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// Migration information
type Migration struct {
	ID   string
	Up   []string
	Down []string
}

type MigrationDirection int

const (
	Up MigrationDirection = iota
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
	raw := strings.Split(src, ";")
	sqls := make([]string, 0, len(raw))
	for _, s := range strings.Split(src, ";\n") {
		sql := strings.Replace(s, "\n", "", -1) + ";"
		sqls = append(sqls, sql)
	}
	sqls = sqls[:len(sqls)-1]
	return sqls
}
