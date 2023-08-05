// Package types define the structure of request parameters and return results in this package
package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

var _ time.Time

// todo generate the request and response struct to here
// delete the templates code start

// CreateUserExampleRequest create request parameters, all fields are mandatory and meet binding rules
// binding instructions for use https://github.com/go-playground/validator
type CreateUserExampleRequest struct {
	Name     string `json:"name" binding:"min=2"`         // username
	Email    string `json:"email" binding:"email"`        // email
	Password string `json:"password" binding:"md5"`       // password
	Phone    string `json:"phone" binding:"e164"`         // phone number, e164 means <+ country code> <cell phone number>.
	Avatar   string `json:"avatar" binding:"min=5"`       // avatar
	Age      int    `json:"age" binding:"gt=0,lt=120"`    // age
	Gender   int    `json:"gender" binding:"gte=0,lte=2"` // gender, 1:Male, 2:Female, other values:unknown
}

// UpdateUserExampleByIDRequest update request parameters, all fields are not required, fields are updated with non-zero values
type UpdateUserExampleByIDRequest struct {
	ID       uint64 `json:"id" binding:"-"`      // id
	Name     string `json:"name" binding:""`     // username
	Email    string `json:"email" binding:""`    // email
	Password string `json:"password" binding:""` // password
	Phone    string `json:"phone" binding:""`    // phone number
	Avatar   string `json:"avatar" binding:""`   // avatar
	Age      int    `json:"age" binding:""`      // age
	Gender   int    `json:"gender" binding:""`   // gender, 1:Male, 2:Female, other values:unknown
}

// GetUserExampleByIDRespond response data
type GetUserExampleByIDRespond struct {
	ID        string    `json:"id"`        // id
	Name      string    `json:"name"`      // username
	Email     string    `json:"email"`     // email
	Phone     string    `json:"phone"`     // phone number
	Avatar    string    `json:"avatar"`    // avatar
	Age       int       `json:"age"`       // age
	Gender    int       `json:"gender"`    // gender, 1:Male, 2:Female, other values:unknown
	Status    int       `json:"status"`    // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt   int64     `json:"loginAt"`   // login timestamp
	CreatedAt time.Time `json:"createdAt"` // create time
	UpdatedAt time.Time `json:"updatedAt"` // update time
}

// delete the templates code end

// DeleteUserExamplesByIDsRequest request form ids
type DeleteUserExamplesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUserExamplesByIDsRequest request form ids
type GetUserExamplesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUserExamplesRequest request form params
type GetUserExamplesRequest struct {
	query.Params // query parameters
}

// ListUserExamplesRespond list data
type ListUserExamplesRespond []struct {
	GetUserExampleByIDRespond
}
