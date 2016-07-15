package mysql

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadRc(t *testing.T) {
	createTestRc()
	defer removeTestRc()

	rc := LoadRc()

	assert.Equal(t, &MysqlRunCommand{
		Mysql: RunCommand{
			Host:     "localhost",
			Port:     "3306",
			User:     "migorate",
			Password: "migorate",
			Database: "migorate",
		},
	}, rc)
}

func TestDatabase(t *testing.T) {
	createTestRc()
	defer removeTestRc()

	db := Database()
	assert.NotNil(t, db)

	// Table migorate_migrations should exist
	_, err := db.Exec("SELECT * FROM migorate_migrations;")
	assert.Nil(t, err)
}

func createTestRc() {
	ioutil.WriteFile(".migoraterc", []byte(`
mysql:
  host: localhost
  port: 3306
  user: migorate
  password: migorate
  database: migorate
`), os.ModePerm)
}

func removeTestRc() {
	os.Remove(".migoraterc")
}
