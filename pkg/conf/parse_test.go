package conf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)
}

func TestShow(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)
	t.Log(Show(c))
}

func Test_replaceDSN(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)

	str := Show(c)

	fmt.Printf(replaceDSN(str))
}

func Test_replacePWD(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c)
	assert.NoError(t, err)

	var keywords []string
	keywords = append(keywords, `"dsn"`, `"password"`)
	str := Show(c)

	fmt.Printf(replacePWD(str, keywords...))
}

func Test_watchConfig(t *testing.T) {
	c := make(map[string]interface{})
	err := Parse("test.yml", &c, func() {
		t.Log("enable watch config file")
	})
	assert.NoError(t, err)

	watchConfig(c)
}
