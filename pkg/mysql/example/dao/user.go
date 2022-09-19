package dao

import (
	"context"
	"errors"

	"github.com/zhufuyi/sponge/pkg/mysql/example/model"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

var _ UserDao = (*userDao)(nil)

// UserDao 定义dao接口
type UserDao interface {
	Create(ctx context.Context, table *model.User) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]model.User, int64, error)
}

type userDao struct {
	*Dao
}

// NewUserDao 创建dao接口
func NewUserDao(dao *Dao) UserDao {
	return &userDao{dao}
}

// Create 创建一条记录，插入记录后，id值被回写到table中
func (d *userDao) Create(ctx context.Context, table *model.User) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID 根据id删除一条记录
func (d *userDao) DeleteByID(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

// Deletes 根据id删除多条记录
func (d *userDao) Deletes(ctx context.Context, ids []uint64) error {
	return d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.User{}).Error
}

// UpdateByID 根据id更新记录
func (d *userDao) UpdateByID(ctx context.Context, table *model.User) error {
	if table.ID < 1 {
		return errors.New("id cannot be less than 0")
	}

	update := map[string]interface{}{}
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Password != "" {
		update["password"] = table.Password
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Phone != "" {
		update["phone"] = table.Phone
	}
	if table.Avatar != "" {
		update["avatar"] = table.Avatar
	}
	if table.Age > 0 {
		update["age"] = table.Age
	}
	if table.Gender > 0 {
		update["gender"] = table.Gender
	}
	if table.LoginAt > 0 {
		update["login_at"] = table.LoginAt
	}

	return d.db.WithContext(ctx).Model(table).Where("id = ?", table.ID).Updates(update).Error
}

// GetByID 根据id获取一条记录
func (d *userDao) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	table := &model.User{}

	err := d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
	if err != nil {
		return nil, err
	}
	return table, nil
}

// GetByColumns 根据分页和列信息筛选多条记录
// params 包括分页参数和查询参数
// 分页参数(必须):

//	page: 页码，从0开始
//	size: 每页行数
//	sort: 排序字段，默认是id倒叙，可以在字段前添加-号表示倒序，没有-号表示升序，多个字段用逗号分隔
//
// 查询参数(非必须):
//
//	name: 列名
//	exp: 表达式，有=、!=、>、>=、<、<=、like七种类型，值为空时默认是=
//	value: 列值
//	logic: 表示逻辑类型，有&(and)、||(or)两种类型，值为空时默认是and
//
// 示例: 查询年龄大于20的男性
//
//	params = &query.Params{
//	    Page: 0,
//	    Size: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "男",
//		},
//	}
func (d *userDao) GetByColumns(ctx context.Context, params *query.Params) ([]model.User, int64, error) {
	query, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = d.db.WithContext(ctx).Model(&model.User{}).Where(query, args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, total, nil
	}

	tables := []model.User{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(query, args...).Find(&tables).Error
	if err != nil {
		return nil, 0, err
	}

	return tables, total, err
}
