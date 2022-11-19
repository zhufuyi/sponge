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
	Phone    string `form:"phone" binding:"e164"`         // phone number, e164 means <+ country code> <cell phone number>.
	Avatar   string `form:"avatar" binding:"min=5"`       // avatar
	Age      int    `form:"age" binding:"gt=0,lt=120"`    // age
	Gender   int    `form:"gender" binding:"gte=0,lte=2"` // gender, 1:Male, 2:Female, other values:unknown
}

// UpdateUserExampleByIDRequest update request parameters, all fields are not required, fields are updated with non-zero values
type UpdateUserExampleByIDRequest struct {
	ID       uint64 `json:"id" binding:"-"`      // id
	Name     string `json:"name" binding:""`     // username
	Email    string `json:"email" binding:""`    // email
	Password string `json:"password" binding:""` // password
	Phone    string `form:"phone" binding:""`    // phone number
	Avatar   string `form:"avatar" binding:""`   // avatar
	Age      int    `form:"age" binding:""`      // age
	Gender   int    `form:"gender" binding:""`   // gender, 1:Male, 2:Female, other values:unknown
}

// GetUserExampleByIDRespond response data
type GetUserExampleByIDRespond struct {
	ID        string    `json:"id"`         // id
	Name      string    `json:"name"`       // username
	Email     string    `json:"email"`      // email
	Phone     string    `json:"phone"`      // phone number
	Avatar    string    `json:"avatar"`     // avatar
	Age       int       `json:"age"`        // age
	Gender    int       `json:"gender"`     // gender, 1:Male, 2:Female, other values:unknown
	Status    int       `json:"status"`     // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt   int64     `json:"login_at"`   // login timestamp
	CreatedAt time.Time `json:"created_at"` // create time
	UpdatedAt time.Time `json:"updated_at"` // update time
}

// delete the templates code end

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
