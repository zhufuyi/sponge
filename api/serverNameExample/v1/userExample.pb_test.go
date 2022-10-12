package v1

import (
	"errors"
	"testing"

	"github.com/zhufuyi/sponge/api/types"

	"google.golang.org/protobuf/runtime/protoimpl"
)

func TestGenderType(t *testing.T) {
	obj := GenderType(1)
	obj.Enum()
	obj.String()
	obj.Descriptor()
	obj.Type()
	obj.Number()
	obj.EnumDescriptor()
}

func TestCreateUserExampleRequest(t *testing.T) {
	obj := &CreateUserExampleRequest{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Name:          "foo",
		Email:         "foo@bar.com",
		Password:      "foo",
		Phone:         "16000000000",
		Avatar:        "http://foo.com/1.jpg",
		Age:           11,
		Gender:        1,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj._validateHostname("localhost")
	obj._validateEmail(obj.Email)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetName()
	obj.GetEmail()
	obj.GetPassword()
	obj.GetPhone()
	obj.GetAvatar()
	obj.GetAge()
	obj.GetGender()
}

func TestCreateUserExampleRequestMultiError(t *testing.T) {
	obj := CreateUserExampleRequestMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestCreateUserExampleRequestValidationError(t *testing.T) {
	obj := CreateUserExampleRequestValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestCreateUserExampleReply(t *testing.T) {
	obj := &CreateUserExampleReply{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            10,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetId()
}

func TestCreateUserExampleReplyMultiError(t *testing.T) {
	obj := CreateUserExampleReplyMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestCreateUserExampleReplyValidationError(t *testing.T) {
	obj := CreateUserExampleReplyValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestDeleteUserExampleByIDRequest(t *testing.T) {
	obj := &DeleteUserExampleByIDRequest{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            10,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetId()
}

func TestDeleteUserExampleByIDRequestMultiError(t *testing.T) {
	obj := DeleteUserExampleByIDRequestMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestDeleteUserExampleByIDRequestValidationError(t *testing.T) {
	obj := DeleteUserExampleByIDRequestValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestDeleteUserExampleByIDReply(t *testing.T) {
	obj := &DeleteUserExampleByIDReply{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
}

func TestDeleteUserExampleByIDReplyMultiError(t *testing.T) {
	obj := DeleteUserExampleByIDReplyMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestDeleteUserExampleByIDReplyValidationError(t *testing.T) {
	obj := DeleteUserExampleByIDReplyValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestUpdateUserExampleByIDRequest(t *testing.T) {
	obj := &UpdateUserExampleByIDRequest{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            10,
		Name:          "foo",
		Email:         "foo@bar.com",
		Password:      "foo",
		Phone:         "16000000000",
		Avatar:        "http://foo.com/1.jpg",
		Age:           11,
		Gender:        1,
		Status:        1,
		LoginAt:       1661556775,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetId()
	obj.GetName()
	obj.GetEmail()
	obj.GetPassword()
	obj.GetPhone()
	obj.GetAvatar()
	obj.GetAge()
	obj.GetGender()
	obj.GetStatus()
	obj.GetLoginAt()
}

func TestUpdateUserExampleByIDRequestMultiError(t *testing.T) {
	obj := UpdateUserExampleByIDRequestMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestUpdateUserExampleByIDRequestValidationError(t *testing.T) {
	obj := UpdateUserExampleByIDRequestValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestUpdateUserExampleByIDReply(t *testing.T) {
	obj := &UpdateUserExampleByIDReply{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
}

func TestUpdateUserExampleByIDReplyMultiError(t *testing.T) {
	obj := UpdateUserExampleByIDReplyMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestUpdateUserExampleByIDReplyValidationError(t *testing.T) {
	obj := UpdateUserExampleByIDReplyValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestUserExample(t *testing.T) {
	obj := &UserExample{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            10,
		Name:          "foo",
		Email:         "foo@bar.com",
		Phone:         "16000000000",
		Avatar:        "http://foo.com/1.jpg",
		Age:           11,
		Gender:        1,
		Status:        1,
		LoginAt:       1661556775,
		CreatedAt:     1661556775,
		UpdatedAt:     1661556775,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetId()
	obj.GetName()
	obj.GetEmail()
	obj.GetPhone()
	obj.GetAvatar()
	obj.GetAge()
	obj.GetGender()
	obj.GetStatus()
	obj.GetLoginAt()
	obj.GetCreatedAt()
	obj.GetUpdatedAt()
}

func TestUserExampleMultiError(t *testing.T) {
	obj := UserExampleMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestUserExampleValidationError(t *testing.T) {
	obj := UserExampleValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestGetUserExampleByIDRequest(t *testing.T) {
	obj := &GetUserExampleByIDRequest{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            10,
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetId()
}

func TestGetUserExampleByIDRequestMultiError(t *testing.T) {
	obj := GetUserExampleByIDRequestMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestGetUserExampleByIDRequestValidationError(t *testing.T) {
	obj := GetUserExampleByIDRequestValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestGetUserExampleByIDReply(t *testing.T) {
	obj := &GetUserExampleByIDReply{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		UserExample: &UserExample{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			Id:            10,
			Name:          "foo",
			Email:         "foo@bar.com",
			Phone:         "16000000000",
			Avatar:        "http://foo.com/1.jpg",
			Age:           11,
			Gender:        1,
			Status:        1,
			LoginAt:       1661556775,
			CreatedAt:     1661556775,
			UpdatedAt:     1661556775,
		},
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetUserExample()
}

func TestGetUserExampleByIDReplyMultiError(t *testing.T) {
	obj := GetUserExampleByIDReplyMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestGetUserExampleByIDReplyValidationError(t *testing.T) {
	obj := GetUserExampleByIDReplyValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestListUserExampleByIDsRequest(t *testing.T) {
	obj := &ListUserExampleByIDsRequest{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Ids:           []uint64{10},
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetIds()
}

func TestListUserExampleByIDsRequestMultiError(t *testing.T) {
	obj := ListUserExampleByIDsRequestMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestListUserExampleByIDsRequestValidationError(t *testing.T) {
	obj := ListUserExampleByIDsRequestValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestListUserExampleByIDsReply(t *testing.T) {
	obj := &ListUserExampleByIDsReply{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		UserExamples: []*UserExample{{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			Id:            10,
			Name:          "foo",
			Email:         "foo@bar.com",
			Phone:         "16000000000",
			Avatar:        "http://foo.com/1.jpg",
			Age:           11,
			Gender:        1,
			Status:        1,
			LoginAt:       1661556775,
			CreatedAt:     1661556775,
			UpdatedAt:     1661556775,
		}},
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetUserExamples()
}

func TestListUserExampleByIDsReplyMultiError(t *testing.T) {
	obj := ListUserExampleByIDsReplyMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestListUserExampleByIDsReplyValidationError(t *testing.T) {
	obj := ListUserExampleByIDsReplyValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestListUserExampleRequest(t *testing.T) {
	obj := &ListUserExampleRequest{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Params: &types.Params{
			Page:  0,
			Limit: 10,
			Sort:  "name",
			Columns: []*types.Column{{
				Name:  "foo",
				Exp:   "=",
				Value: "bar",
				Logic: "AND",
			},
			},
		},
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetParams()
}

func TestListUserExampleRequestMultiError(t *testing.T) {
	obj := ListUserExampleRequestMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestListUserExampleRequestValidationError(t *testing.T) {
	obj := ListUserExampleRequestValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}

// ------------------------------------------------------------------------------------------

func TestListUserExampleReply(t *testing.T) {
	obj := &ListUserExampleReply{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Total:         1,
		UserExamples: []*UserExample{{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			Id:            10,
			Name:          "foo",
			Email:         "foo@bar.com",
			Phone:         "16000000000",
			Avatar:        "http://foo.com/1.jpg",
			Age:           11,
			Gender:        1,
			Status:        1,
			LoginAt:       1661556775,
			CreatedAt:     1661556775,
			UpdatedAt:     1661556775,
		}},
	}

	obj.Validate()
	obj.ValidateAll()
	obj.validate(true)
	obj.Reset()
	obj.String()
	obj.ProtoMessage()
	obj.ProtoReflect()
	obj.Descriptor()
	obj.GetTotal()
	obj.GetUserExamples()
}

func TestListUserExampleReplyMultiError(t *testing.T) {
	obj := ListUserExampleReplyMultiError{errors.New("mock error")}
	obj.Error()
	obj.AllErrors()
}

func TestListUserExampleReplyValidationError(t *testing.T) {
	obj := ListUserExampleReplyValidationError{
		field:  "",
		reason: "",
		cause:  nil,
		key:    false,
	}

	obj.Field()
	obj.Reason()
	obj.Cause()
	obj.Key()
	obj.ErrorName()
	obj.Error()
}
