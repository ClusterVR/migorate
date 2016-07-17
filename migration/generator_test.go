package migration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSQL(t *testing.T) {
	sql := generateSQL("create_users", []string{"id:id", "name:string", "login_count:integer", "last_login_at:datetime", "created_at:timestamp"})
	assert.Equal(t, `-- +migrate Up
CREATE TABLE users(id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), login_count INT, last_login_at DATETIME, created_at TIMESTAMP);

-- +migrate Down
DROP TABLE users;
`, sql)
}
