package mysql

import (
	"database/sql"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type RunCommand struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type MysqlRunCommand struct {
	Mysql RunCommand
}

func LoadRc() *MysqlRunCommand {
	buf, err := ioutil.ReadFile(".migoraterc")
	if err != nil {
		log.Fatalf("Failed to load .migoraterc: %v", err)
	}
	m := MysqlRunCommand{}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		log.Fatalf("Failed to load .migoraterc as YAML: %v", err)
	}
	return &m
}

func Database() *sql.DB {
	rc := LoadRc()
	uri := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", rc.Mysql.User, rc.Mysql.Password, rc.Mysql.Host, rc.Mysql.Port, rc.Mysql.Database)
	fmt.Printf("connect to %v", uri)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS migorate_migrations(id VARCHAR(255) PRIMARY KEY, migrated_at TIMESTAMP);")
	if err != nil {
		log.Fatalf("Failed to create migration management table: %v", err)
	}
	return db
}
