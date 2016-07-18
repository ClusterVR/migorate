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
	var fk string
	for _, c := range cols {
		col := strings.Split(c, ":")
		colname, t, f := convertType(col[0], col[1])
		fk += f
		convertedCols += colname + " " + t + ", "
	}
	if len(fk) > 0 {
		convertedCols += fk
		convertedCols = convertedCols[0 : len(convertedCols)-2]
	} else if len(convertedCols) > 0 {
		convertedCols = convertedCols[0 : len(convertedCols)-2]
	}
	return fmt.Sprintf(template, table, convertedCols, table)
}

func convertType(c string, t string) (col string, typ string, fk string) {
	switch t {
	case "id":
		return c, "INT PRIMARY KEY AUTO_INCREMENT", ""
	case "integer":
		return c, "INT", ""
	case "datetime":
		return c, "DATETIME", ""
	case "string":
		return c, "VARCHAR(255)", ""
	case "timestamp":
		return c, "TIMESTAMP", ""
	case "references":
		col = c + "_id"
		return col, "INT", "FOREIGN KEY (" + col + ") REFERENCES " + c + "s(id), "
	}
	return c, t, ""
}
