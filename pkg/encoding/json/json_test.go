package json

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
	"testing"

	"github.com/stretchr/testify/assert"
)

type obj struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func TestJSON(t *testing.T) {
	c := codec{}

	name := c.Name()
	assert.Equal(t, Name, name)

	data, err := c.Marshal(&obj{ID: 1, Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, data)

	o := new(obj)
	err = c.Unmarshal(data, o)
	assert.NoError(t, err)
	assert.Equal(t, "foo", o.Name)
}

type obj2 struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (o obj2) MarshalJSON() ([]byte, error) {
	return []byte("test data"), nil
}

func TestJSON2(t *testing.T) {
	c := codec{}
	b, err := c.Marshal(&obj2{})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, b)

	err = c.Unmarshal(b, &obj2{})
	assert.NotNil(t, err)
}

type obj3 struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (o obj3) ProtoReflect() protoreflect.Message {
	req := &pluginpb.CodeGeneratorRequest{}
	opts := protogen.Options{}
	gen, _ := opts.New(req)

	return gen.Response().ProtoReflect()
}

func TestJSON3(t *testing.T) {
	c := codec{}
	b, err := c.Marshal(&obj3{})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, b)

	err = c.Unmarshal(b, &obj3{})
	assert.NoError(t, err)
}
