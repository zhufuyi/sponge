package parse

import (
	"google.golang.org/protobuf/compiler/protogen"
	"testing"
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
	pss := GetServices(file)
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
}
