package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

var _ time.Time

// todo generate the request and response struct to here
// delete the templates code start

// CreateUserExampleRequest 创建请求参数，所有字段是必须的，并且满足binding规则
// binding使用说明 https://github.com/go-playground/validator
type CreateUserExampleRequest struct {
	Name     string `json:"name" binding:"min=2"`         // 名称
	Email    string `json:"email" binding:"email"`        // 邮件
	Password string `json:"password" binding:"md5"`       // 密码
	Phone    string `form:"phone" binding:"e164"`         // 手机号码，e164表示<+国家编号><手机号码>
	Avatar   string `form:"avatar" binding:"min=5"`       // 头像
	Age      int    `form:"age" binding:"gt=0,lt=120"`    // 年龄
	Gender   int    `form:"gender" binding:"gte=0,lte=2"` // 性别，1:男，2:女
}

// UpdateUserExampleByIDRequest 更新请求参数，所有字段不是必须的，字段为非零值更新
type UpdateUserExampleByIDRequest struct {
	ID       uint64 `json:"id" binding:"-"`      // id
	Name     string `json:"name" binding:""`     // 名称
	Email    string `json:"email" binding:""`    // 邮件
	Password string `json:"password" binding:""` // 密码
	Phone    string `form:"phone" binding:""`    // 手机号码
	Avatar   string `form:"avatar" binding:""`   // 头像
	Age      int    `form:"age" binding:""`      // 年龄
	Gender   int    `form:"gender" binding:""`   // 性别，1:男，2:女
}

// GetUserExampleByIDRespond 返回数据
type GetUserExampleByIDRespond struct {
	ID        string    `json:"id"`         // id
	Name      string    `json:"name"`       // 名称
	Email     string    `json:"email"`      // 邮件
	Phone     string    `json:"phone"`      // 手机号码
	Avatar    string    `json:"avatar"`     // 头像
	Age       int       `json:"age"`        // 年龄
	Gender    int       `json:"gender"`     // 性别，1:男，2:女
	Status    int       `json:"status"`     // 账号状态
	LoginAt   int64     `json:"login_at"`   // 登录时间戳
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// delete the templates code end

// GetUserExamplesByIDsRequest request form ids
type GetUserExamplesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id列表
}

// GetUserExamplesRequest request form params
type GetUserExamplesRequest struct {
	query.Params // 查询参数
}

// ListUserExamplesRespond list data
type ListUserExamplesRespond []struct {
	GetUserExampleByIDRespond
}
