package migration

import (
	"fmt"
	"regexp"
	"strings"
)

func generateSQL(name string, cols []string) string {
	r := regexp.MustCompile("^create_(.+)$")
	if r.MatchString(name) {
		g := r.FindSubmatch([]byte(name))
		return generateCreateTable(string(g[1]), cols)
	}
	return `-- +migrate Up


-- +migrate Down

`
}

func generateCreateTable(table string, cols []string) string {
	template := `-- +migrate Up
CREATE TABLE %v(%v);

-- +migrate Down
DROP TABLE %v;
`
	var convertedCols string
	for _, c := range cols {
		col := strings.Split(c, ":")
		t := convertType(col[1])
		convertedCols += col[0] + " " + t + ", "
	}
	if len(convertedCols) > 0 {
		convertedCols = convertedCols[0 : len(convertedCols)-2]
	}
	return fmt.Sprintf(template, table, convertedCols, table)
}

func convertType(t string) string {
	switch t {
	case "id":
		return "INT PRIMARY KEY AUTO_INCREMENT"
	case "integer":
		return "INT"
	case "datetime":
		return "DATETIME"
	case "string":
		return "VARCHAR(255)"
	case "timestamp":
		return "TIMESTAMP"
	}
	return t
}
