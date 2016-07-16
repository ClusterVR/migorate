package migration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExec(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	migrations := Plan(testMigrationPath, Up)
	for _, m := range *migrations {
		m.Exec(db, Up)
	}

	rows, err := db.Query("SHOW TABLES LIKE 'users'")
	assert.True(t, rows.Next(), "Table users should exits")
	assert.Nil(t, err, "Error should NOT occurred: %+v", err)

	rows, err = db.Query("SHOW TABLES LIKE 'books'")
	assert.True(t, rows.Next(), "Table books should exits")
	assert.Nil(t, err, "Error should NOT occurred: %+v", err)
}
