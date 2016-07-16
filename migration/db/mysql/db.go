package mysql

import (
	"database/sql"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql" // Use mysql driver
)

// RunCommand is configuration for database connection
type RunCommand struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// MysqlRunCommand is configuration for MySQL connection
type mySQLRunCommand struct {
	Mysql RunCommand
}

func loadRc() *mySQLRunCommand {
	buf, err := ioutil.ReadFile(".migoraterc")
	if err != nil {
		log.Fatalf("Failed to load .migoraterc: %v\n", err)
	}
	m := mySQLRunCommand{}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		log.Fatalf("Failed to load .migoraterc as YAML: %v\n", err)
	}
	return &m
}

// Database which is connected to MySQL
func Database() *sql.DB {
	rc := loadRc()
	uri := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", rc.Mysql.User, rc.Mysql.Password, rc.Mysql.Host, rc.Mysql.Port, rc.Mysql.Database)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS migorate_migrations(id VARCHAR(255) PRIMARY KEY, migrated_at TIMESTAMP);")
	if err != nil {
		log.Fatalf("Failed to create migration management table: %v\n", err)
	}
	return db
}
