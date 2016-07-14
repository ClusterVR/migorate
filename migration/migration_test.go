package migration

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
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
	migrations := *Plan("../test/fixtures/1_two_migrations", Up)
	assert.Equal(t, 2, len(migrations), "Expect 2 migration found but %v found.", len(migrations))

	assert.Equal(t, "20160714092556_create_users", migrations[0].Id, "Migration id")
	assert.Equal(t, 2, len(migrations[0].Up), "%+v", migrations[0].Up)
	assert.Equal(t, 1, len(migrations[0].Down), "%+v", migrations[0].Down)
	assert.Equal(t, "CREATE TABLE users(id PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), email VARCHAR(255), created_at TIMESTAMP);", migrations[0].Up[0])
	assert.Equal(t, "ALTER TABLE users ADD INDEX index_users_email(email);", migrations[0].Up[1])
	assert.Equal(t, "DROP TABLE users;", migrations[0].Down[0])

	assert.Equal(t, "20160714092604_create_books", migrations[1].Id, "Migration id")
	assert.Equal(t, 1, len(migrations[1].Up), "%+v", migrations[1].Up)
	assert.Equal(t, 1, len(migrations[1].Down), "%+v", migrations[1].Down)
	assert.Equal(t, "CREATE TABLE books(id PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255), author VARCHAR(255), created_at TIMESTAMP);", migrations[1].Up[0])
	assert.Equal(t, "DROP TABLE books;", migrations[1].Down[0])
}