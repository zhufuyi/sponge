package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestGetServices(t *testing.T) {
	method := &protogen.Method{
		GoName: "Create",
		Input: &protogen.Message{
			GoIdent: protogen.GoIdent{GoName: "CreateRequest"},
		},
		Output: &protogen.Message{
			GoIdent: protogen.GoIdent{GoName: "CreateReply"},
		},
	}

	service := &protogen.Service{
		Desc:     nil,
		GoName:   "Greeter",
		Methods:  []*protogen.Method{method},
		Location: protogen.Location{},
		Comments: protogen.CommentSet{},
	}
	var file = &protogen.File{
		Services: []*protogen.Service{service},
	}
	pss := GetServices("greeter", file)
	for _, s := range pss {
		t.Logf("%+v", s)
		for _, m := range s.Methods {
			t.Logf("%+v", m)
		}
	}
}

func TestMethods(t *testing.T) {
	sm := &ServiceMethod{}
	t.Log(sm.AddOne(1))

	s := PbService{}
	t.Log(s.RandNumber())

	typeNames := []string{"bool", "int32", "float", "string", "unknown", "repeated uint64"}
	field := &RequestField{}
	for _, v := range typeNames {
		field.FieldType = v
		t.Log(field.GoTypeZero())
	}
}

func Test_getRequestFields(t *testing.T) {
	fields := []*protogen.Field{
		{
			Desc:     &desc{},
			GoName:   "foo1",
			Comments: protogen.CommentSet{},
		},
	}

	requestFields := getRequestFields(fields)
	assert.Equal(t, 1, len(requestFields))
}

type desc struct {
	protoreflect.FieldDescriptor
}

func (d desc) Kind() protoreflect.Kind {
	return protoreflect.Int32Kind
}

func (d desc) Cardinality() protoreflect.Cardinality {
	return protoreflect.Repeated
}
