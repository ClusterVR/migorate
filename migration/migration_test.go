package migration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExec(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	migrations := Plan(testMigrationPath, Up, "")
	for _, m := range *migrations {
		m.Exec(db, Up)
	}

	rows, err := db.Query("SHOW TABLES LIKE 'users'")
	assert.True(t, rows.Next(), "Table users should exits")
	assert.Nil(t, err, "Error should NOT occurred: %+v", err)

	rows, err = db.Query("SELECT COUNT(id) FROM migorate_migrations WHERE id = ?", "20160714092556_create_users")
	assert.Equal(t, 1, count(rows), "Migration should be inserted")
	assert.Nil(t, err, "Error should NOT occurred: %+v", err)

	rows, err = db.Query("SHOW TABLES LIKE 'books'")
	assert.True(t, rows.Next(), "Table books should exits")
	assert.Nil(t, err, "Error should NOT occurred: %+v", err)

	rows, err = db.Query("SELECT COUNT(id) FROM migorate_migrations WHERE id = ?", "20160714092604_create_books")
	assert.Equal(t, 1, count(rows), "Migration should be inserted")
	assert.Nil(t, err, "Error should NOT occurred: %+v", err)
}
