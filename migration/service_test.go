package migration

import (
	"database/sql"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mizoguche/migorate/migration/db/mysql"
	"github.com/stretchr/testify/assert"
)

const testMigrationPath = "../test/fixtures/success_migrations"

func TestGenerate(t *testing.T) {
	os.Mkdir("tmp", 0777)
	name := "test_migration"

	Generate("tmp", name, []string{})

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

	migrations := *Plan(testMigrationPath, Up, "all")
	assert.Equal(t, 3, len(migrations), "Expect 2 migration found but %v found.", len(migrations))

	assertCreateUsersMigration(t, migrations[0])
	assertCreateBooksMigration(t, migrations[1])
	assertCreateAuthorsMigration(t, migrations[2])
}

func TestPlanToBooks(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	migrations := *Plan(testMigrationPath, Up, "20160714092604_create_books")
	assert.Equal(t, 2, len(migrations), "Expect 2 migration found but %v found.", len(migrations))

	assertCreateUsersMigration(t, migrations[0])
	assertCreateBooksMigration(t, migrations[1])
}

func TestPlanWhenAlreadyMigratedBooks(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092604_create_books', NOW());")

	migrations := *Plan(testMigrationPath, Up, "all")
	assert.Equal(t, 2, len(migrations), "Expect 1 migration found but %v found.", len(migrations))
	assertCreateUsersMigration(t, migrations[0])
	assertCreateAuthorsMigration(t, migrations[1])
}

func TestPlanWhenAlreadyMigrated(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092556_create_users', NOW());")

	migrations := *Plan(testMigrationPath, Up, "all")
	assert.Equal(t, 2, len(migrations), "Expect 1 migration found but %v found.", len(migrations))
	assertCreateBooksMigration(t, migrations[0])
	assertCreateAuthorsMigration(t, migrations[1])
}

func TestPlanRollback(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092556_create_users', NOW());")
	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092604_create_books', NOW());")
	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160716102604_create_authors', NOW());")

	migrations := *Plan(testMigrationPath, Down, "all")
	assert.Equal(t, 3, len(migrations), "Expect 2 migration found but %v found.", len(migrations))

	assertCreateAuthorsMigration(t, migrations[0])
	assertCreateBooksMigration(t, migrations[1])
	assertCreateUsersMigration(t, migrations[2])
}

func TestPlanRollbackToBooks(t *testing.T) {
	db := initDb()
	defer cleanupDb(db)

	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092556_create_users', NOW());")
	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160714092604_create_books', NOW());")
	db.Exec("INSERT INTO migorate_migrations(id, migrated_at) VALUES('20160716102604_create_authors', NOW());")

	migrations := *Plan(testMigrationPath, Down, "20160714092604_create_books")
	assert.Equal(t, 2, len(migrations), "Expect 2 migration found but %v found.", len(migrations))

	assertCreateAuthorsMigration(t, migrations[0])
	assertCreateBooksMigration(t, migrations[1])
}

func assertCreateUsersMigration(t *testing.T, m Migration) {
	assert.Equal(t, "20160714092556_create_users", m.ID, "Migration id")
	assert.Equal(t, 2, len(m.Up), "%+v", m.Up)
	assert.Equal(t, 1, len(m.Down), "%+v", m.Down)
	assertEqualStripQuery(t, "CREATE TABLE users ( id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), email VARCHAR(255), created_at TIMESTAMP);", m.Up[0])
	assertEqualStripQuery(t, "ALTER TABLE users ADD INDEX index_users_email(email);", m.Up[1])
	assertEqualStripQuery(t, "DROP TABLE users;", m.Down[0])
}

func assertCreateBooksMigration(t *testing.T, m Migration) {
	assert.Equal(t, "20160714092604_create_books", m.ID, "Migration id")
	assert.Equal(t, 1, len(m.Up), "%+v", m.Up)
	assert.Equal(t, 1, len(m.Down), "%+v", m.Down)
	assertEqualStripQuery(t, "CREATE TABLE books ( id INT PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255), author VARCHAR(255), created_at TIMESTAMP);", m.Up[0])
	assertEqualStripQuery(t, "DROP TABLE books;", m.Down[0])
}

func assertCreateAuthorsMigration(t *testing.T, m Migration) {
	assert.Equal(t, "20160716102604_create_authors", m.ID, "Migration id")
	assert.Equal(t, 4, len(m.Up), "%+v", m.Up)
	assert.Equal(t, 3, len(m.Down), "%+v", m.Down)
	// up
	assertEqualStripQuery(t, "CREATE TABLE authors ( id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), created_at TIMESTAMP);", m.Up[0])
	assertEqualStripQuery(t, "ALTER TABLE books DROP COLUMN author;", m.Up[1])
	assertEqualStripQuery(t, "ALTER TABLE books ADD (author_id INT NOT NULL);", m.Up[2])
	assertEqualStripQuery(t, "ALTER TABLE books ADD CONSTRAINT fk_books_author_id FOREIGN KEY (author_id) REFERENCES authors (id);", m.Up[3])
	// down
	assertEqualStripQuery(t, "ALTER TABLE books DROP FOREIGN KEY fk_books_author_id;", m.Down[0])
	assertEqualStripQuery(t, "DROP TABLE authors;", m.Down[1])
	assertEqualStripQuery(t, "ALTER TABLE books ADD (author VARCHAR(255));", m.Down[2])
}

func assertEqualStripQuery(t *testing.T, expected, received string) {
	assert.Equal(t, stripQuery(expected), stripQuery(received))
}

func initDb() *sql.DB {
	buf, _ := ioutil.ReadFile("../test/rc/mysql.yml")
	ioutil.WriteFile(mysql.RCFilePath, buf, os.ModePerm)

	db := mysql.Database()
	cleanUpTables(db)
	return db
}

func cleanupDb(db *sql.DB) {
	os.Remove(mysql.RCFilePath)
	cleanUpTables(db)
	db.Close()
}

func cleanUpTables(db *sql.DB) {
	db.Exec("DELETE FROM migorate_migrations")
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE books")
	db.Exec("DROP TABLE authors")
}

// testHelper

var re = regexp.MustCompile("\\s+")

// strip out new lines and trim spaces
func stripQuery(q string) (s string) {
	return strings.TrimSpace(re.ReplaceAllString(q, " "))
}
