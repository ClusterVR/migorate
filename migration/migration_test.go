package migration

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mizoguche/migorate/migration/db/mysql"
)

func TestGenerate(t *testing.T) {
	os.Mkdir("tmp", 0777)
	name := "test_migration"

	Generate("tmp", name)

	files, _ := ioutil.ReadDir("tmp")
	assert.Equal(t, 1, len(files), "Expected 1 file generated.")

	r := regexp.MustCompile(`\d\d\d\d\d\d\d\d\d\d\d\d\d\d_` + name + ".sql")
	assert.True(t, r.MatchString(files[0].Name()), "Filename \"%v\" is not formatted.", files[0].Name())

	buf, _ := ioutil.ReadFile("tmp/" + files[0].Name())
	s := string(buf)
	r = regexp.MustCompile(`(?m)^-- \+migrate Up$`)
	assert.True(t, r.MatchString(s), "Generated file does not contains template \"-- +migrate Up\"")

	r = regexp.MustCompile(`(?m)^-- \+migrate Down$`)
	assert.True(t, r.MatchString(s), "Generated file does not contains template \"-- +migrate Down\"")

	os.RemoveAll("tmp")
}

func TestPlan(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	migrations := *Plan("../test/fixtures/1_two_migrations", Up)
	assert.Equal(t, 2, len(migrations), "Expect 2 migration found but %v found.", len(migrations))

	assertCreateUsersMigration(t, migrations[0])
	assertCreateBooksMigration(t, migrations[1])
}

func TestPlanWhenAlreadyMigratedLastFile(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092604_create_books', NOW());")

	migrations := *Plan("../test/fixtures/1_two_migrations", Up)
	assert.Equal(t, 1, len(migrations), "Expect 1 migration found but %v found.", len(migrations))
	assertCreateUsersMigration(t, migrations[0])

	db.Close()
}

func TestPlanWhenAlreadyMigrated(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	res, err := db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092556_create_users', NOW());")
	fmt.Print(res)
	fmt.Print(err)

	migrations := *Plan("../test/fixtures/1_two_migrations", Up)
	assert.Equal(t, 1, len(migrations), "Expect 1 migration found but %v found.", len(migrations))
	assertCreateBooksMigration(t, migrations[0])

	db.Close()
}

func assertCreateUsersMigration(t *testing.T, m Migration) {
	assert.Equal(t, "20160714092556_create_users", m.Id, "Migration id")
	assert.Equal(t, 2, len(m.Up), "%+v", m.Up)
	assert.Equal(t, 1, len(m.Down), "%+v", m.Down)
	assert.Equal(t, "CREATE TABLE users(id PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), email VARCHAR(255), created_at TIMESTAMP);", m.Up[0])
	assert.Equal(t, "ALTER TABLE users ADD INDEX index_users_email(email);", m.Up[1])
	assert.Equal(t, "DROP TABLE users;", m.Down[0])
}

func assertCreateBooksMigration(t *testing.T, m Migration) {
	assert.Equal(t, "20160714092604_create_books", m.Id, "Migration id")
	assert.Equal(t, 1, len(m.Up), "%+v", m.Up)
	assert.Equal(t, 1, len(m.Down), "%+v", m.Down)
	assert.Equal(t, "CREATE TABLE books(id PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255), author VARCHAR(255), created_at TIMESTAMP);", m.Up[0])
	assert.Equal(t, "DROP TABLE books;", m.Down[0])
}

func initDb() *sql.DB {
	buf, _ := ioutil.ReadFile("../test/rc/mysql.yml")
	ioutil.WriteFile(".migoraterc", buf, os.ModePerm)

	db := mysql.Database()
	db.Exec("DELETE FROM migorate_migrations")
	return db
}

func cleanupDb(db *sql.DB) {
	os.Remove(".migoraterc")
	db.Exec("DELETE FROM migorate_migrations")
}
