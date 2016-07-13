package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func Generate(dir string, name string) error {
	t := time.Now()
	id := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	filepath := fmt.Sprintf("%s/%s_%s.sql", dir, id, name)
	content := []byte(`-- +migrate up


-- +migrate down

`)
	err := ioutil.WriteFile(filepath, content, os.ModePerm)
	if err != nil {
		log.Printf("Failed to generate file\n%v", err)
		return err
	}

	log.Printf("Generated: %v", filepath)
	return nil
}
