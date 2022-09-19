// todo generate model codes to here
// delete the templates code start

package model

import (
	"github.com/zhufuyi/sponge/pkg/mysql"
)

// UserExample object fields mapping table
type UserExample struct {
	mysql.Model `gorm:"embedded"`

	Name     string `gorm:"column:name;NOT NULL" json:"name"`         // 用户名
	Password string `gorm:"column:password;NOT NULL" json:"password"` // 密码
	Email    string `gorm:"column:email;NOT NULL" json:"email"`       // 邮件
	Phone    string `gorm:"column:phone;NOT NULL" json:"phone"`       // 手机号码
	Avatar   string `gorm:"column:avatar;NOT NULL" json:"avatar"`     // 头像
	Age      int    `gorm:"column:age;NOT NULL" json:"age"`           // 年龄
	Gender   int    `gorm:"column:gender;NOT NULL" json:"gender"`     // 性别，1:男，2:女，其他值:未知
	Status   int    `gorm:"column:status;NOT NULL" json:"status"`     // 账号状态，1:未激活，2:已激活，3:封禁
	LoginAt  int64  `gorm:"column:login_at;NOT NULL" json:"login_at"` // 登录时间戳
}

// TableName get table name
func (table *UserExample) TableName() string {
	return "user_example"
}

// delete the templates code end
