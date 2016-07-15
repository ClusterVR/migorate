package mysql

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadRc(t *testing.T) {
	ioutil.WriteFile(".migoraterc", []byte(`
mysql:
  host: localhost
  port: 3306
  user: migorate
  password: migoratepassword
  database: migorate_database
`), os.ModePerm)
	defer os.Remove(".migoraterc")

	rc := LoadRc()

	assert.Equal(t, &MysqlRunCommand{
		Mysql: RunCommand{
			Host:     "localhost",
			Port:     "3306",
			User:     "migorate",
			Password: "migoratepassword",
			Database: "migorate_database",
		},
	}, rc)
}
