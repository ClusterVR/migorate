package migration

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestGenerate(t *testing.T) {
	os.Mkdir("tmp", 0777)
	name := "test_migration"

	Generate("tmp", name)

	files, _ := ioutil.ReadDir("tmp")
	assert.Equal(t, 1, len(files), "Expected 1 file generated.")

	r := regexp.MustCompile(`\d\d\d\d\d\d\d\d\d\d\d\d\d\d_` + name + ".sql")
	assert.True(t, r.MatchString(files[0].Name()), "Filename \"%v\" is not formatted.", files[0].Name())

	buf, _ := ioutil.ReadFile("tmp/" + files[0].Name())
	s := string(buf)
	r = regexp.MustCompile(`(?m)^-- \+migrate up$`)
	assert.True(t, r.MatchString(s), "Generated file does not contains template \"-- +migrate up\"")

	r = regexp.MustCompile(`(?m)^-- \+migrate down$`)
	assert.True(t, r.MatchString(s), "Generated file does not contains template \"-- +migrate down\"")

	os.RemoveAll("tmp")
}
