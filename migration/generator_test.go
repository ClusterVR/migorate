package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSQL(t *testing.T) {
	sql := generateSQL("create_users", []string{"id:id", "name:string", "login_count:integer", "last_login_at:datetime", "created_at:timestamp", "comment:references", "book:references"})
	assert.Equal(t, `-- +migrate Up
CREATE TABLE users(id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), login_count INT, last_login_at DATETIME, created_at TIMESTAMP, comment_id INT, book_id INT, FOREIGN KEY (comment_id) REFERENCES comments(id), FOREIGN KEY (book_id) REFERENCES books(id));

-- +migrate Down
DROP TABLE users;
`, sql)
}
