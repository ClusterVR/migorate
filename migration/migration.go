package migration

import (
	"fmt"
	"log"
	"os"
	"time"
	"io/ioutil"
)

func Generate(dir string, name string) {
	t := time.Now()
	id := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	filepath := fmt.Sprintf("%s/%s_%s.sql", dir, id, name)
	content := []byte(`-- +migrate up


-- +migrate down

`)
	err := ioutil.WriteFile(filepath, content, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to generate file\n%v", err)
	}
	log.Printf("Generated: %v", filepath)
}
