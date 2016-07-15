package mysql

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
