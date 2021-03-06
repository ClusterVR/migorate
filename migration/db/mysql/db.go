package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// RCFilePath configuration file path
var RCFilePath = ".migoraterc"

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
	buf, err := ioutil.ReadFile(RCFilePath)
	if err != nil {
		log.Fatalf("Failed to load %s: %v\n", RCFilePath, err)
	}
	m := mySQLRunCommand{}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		log.Fatalf("Failed to load %s as YAML: %v\n", RCFilePath, err)
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS migorate_migrations(id VARCHAR(64) PRIMARY KEY, migrated_at TIMESTAMP);")
	if err != nil {
		log.Fatalf("Failed to create migration management table: %v\n", err)
	}
	return db
}
