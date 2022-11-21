package benchmark

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestNew(t *testing.T) {
	_, err := New("localhost", "testProto/test.proto", "Create", nil, 100)
	assert.NoError(t, err)

	_, err = New("localhost", "testProto/test2.proto", "Create", nil, 100)
	assert.Error(t, err)

	_, err = New("localhost", "testProto/test3.proto", "Create", nil, 100)
	assert.Error(t, err)

	_, err = New("localhost", "testProto/test4.proto", "Create", nil, 100)
	assert.Error(t, err)
}

func Test_params_Run(t *testing.T) {
	req := &pluginpb.CodeGeneratorRequest{}
	opts := protogen.Options{}
	gen, err := opts.New(req)
	o1 := gen.Response()

	b, err := New("localhost", "testProto/test.proto", "Create", o1, 2)
	assert.NoError(t, err)

	err = b.Run()
	t.Log(err) //  replace github.com/golang/protobuf/proto --> google.golang.org/protobuf/proto
	time.Sleep(time.Second)
}
