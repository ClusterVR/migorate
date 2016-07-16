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

	assert.True(t, len(rc.Mysql.Host) > 0)
	assert.Equal(t, "3306", rc.Mysql.Port)
	assert.Equal(t, "migorate", rc.Mysql.User)
	assert.Equal(t, "migorate", rc.Mysql.Password)
	assert.Equal(t, "migorate", rc.Mysql.Database)
}

func TestDatabase(t *testing.T) {
	createTestRc()
	defer removeTestRc()

	db := Database()
	assert.NotNil(t, db)

	// Table migorate_migrations should exist
	_, err := db.Exec("SELECT * FROM migorate_migrations;")
	assert.Nil(t, err)
	db.Close()
}

func createTestRc() {
	buf, _ := ioutil.ReadFile("../../../test/rc/mysql.yml")
	ioutil.WriteFile(".migoraterc", buf, os.ModePerm)
}

func removeTestRc() {
	os.Remove(".migoraterc")
}
