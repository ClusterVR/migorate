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
	buf, _ := ioutil.ReadFile("../../../test/rc/mysql.yml")
	ioutil.WriteFile(".migoraterc", buf, os.ModePerm)
}

func removeTestRc() {
	os.Remove(".migoraterc")
}
