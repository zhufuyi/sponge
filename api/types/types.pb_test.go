package types

import (
	"testing"

	"google.golang.org/protobuf/runtime/protoimpl"
)

func TestColumn(t *testing.T) {
	obj := &Column{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Name:          "foo",
		Exp:           "=",
		Value:         "bar",
		Logic:         "AND",
	}

	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetName()
	obj.GetExp()
	obj.GetValue()
	obj.GetLogic()
}

func TestParams(t *testing.T) {
	obj := &Params{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Page:          0,
		Limit:         10,
		Sort:          "name",
		Columns: []*Column{{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			Name:          "foo",
			Exp:           "=",
			Value:         "bar",
			Logic:         "AND",
		},
		},
	}

	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetPage()
	obj.GetLimit()
	obj.GetSort()
	obj.GetColumns()
}
