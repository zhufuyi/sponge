package jy2struct

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseError(t *testing.T) {
	testData := `foo:bar`
	r := strings.NewReader(testData)
	_, err := ParseJSON(r)
	assert.Error(t, err)

	testData = `	foo: bar`
	r = strings.NewReader(testData)
	_, err = ParseYaml(r)
	assert.Error(t, err)

	_, err = jyParse(r, ParseYaml, "", "", nil, false, false)
	assert.Error(t, err)

	v := FmtFieldName("")
	v = lintFieldName(v)
	assert.Equal(t, "_", v)

	v = stringifyFirstChar("2foo")
	assert.Equal(t, "two_foo", v)
}

func Test_convertKeysToStrings(t *testing.T) {
	testData := map[interface{}]interface{}{"foo": "bar"}
	v := convertKeysToStrings(testData)
	assert.NotNil(t, v)
}

func Test_mergeElements(t *testing.T) {
	testData := "foo"
	v := mergeElements(testData)
	assert.Equal(t, testData, v)

	testData2 := []interface{}{}
	v = mergeElements(testData2)
	assert.Empty(t, v)
	testData2 = []interface{}{"foo", "bar"}
	v = mergeElements(testData2)
	assert.Equal(t, testData2[0], v.([]interface{})[0])
}

func Test_mergeObjects(t *testing.T) {
	var (
		o1 = []interface{}{"foo", "bar"}
		o2 = map[string]interface{}{"foo": "bar"}
		o3 = map[interface{}]interface{}{"foo": "bar"}
	)
	v := mergeObjects(nil, o2)
	assert.Equal(t, o2, v)
	v = mergeObjects(o1, nil)
	assert.Equal(t, o1, v)
	v = mergeObjects(o1, o2)
	assert.Nil(t, v)
	v = mergeObjects("foo", "bar")
	assert.Equal(t, "foo", v)

	v = mergeObjects(o1, o1)
	t.Log(v)
	v = mergeObjects(o2, o2)
	t.Log(v)
	v = mergeObjects(o3, o3)
	t.Log(v)
}
