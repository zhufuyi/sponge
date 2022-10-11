package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestProto(t *testing.T) {
	c := codec{}

	name := c.Name()
	assert.Equal(t, Name, name)

	req := &pluginpb.CodeGeneratorRequest{}
	opts := protogen.Options{}
	gen, err := opts.New(req)
	o1 := gen.Response()

	b, err := c.Marshal(o1)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, b)

	o2 := new(pluginpb.CodeGeneratorRequest)
	err = c.Unmarshal(b, o2)
	assert.NoError(t, err)
}

func TestProtoError(t *testing.T) {
	c := codec{}

	_, err := c.Marshal(nil)
	assert.Error(t, err)
	err = c.Unmarshal(nil, nil)
	assert.Error(t, err)
}
