package gotest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestNewDao(t *testing.T) {
	now := time.Now()
	testData := &User{
		ID:        1,
		Name:      "foo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	d := NewDao(nil, testData)
	d.IDao = newUserDao(d.DB)
	defer d.Close()
}

func TestDao_GetAnyArgs(t *testing.T) {
	now := time.Now()
	testData := &User{
		ID:        1,
		Name:      "foo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	d := NewDao(nil, testData)
	d.IDao = newUserDao(d.DB)
	defer d.Close()

	t.Log(d.GetAnyArgs(testData))

	// test error
	defer func() {
		recover()
	}()
	d.GetAnyArgs(make(chan string))
}

func TestAnyTime_Match(t *testing.T) {
	d := NewDao(nil, &User{ID: 1})
	defer d.Close()
	t.Log(d.AnyTime.Match(time.Now()), d.AnyTime.Match("test"))
}

// -------------------------- 增删改查示例 --------------------------------

type User struct {
	ID        uint64    `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	Name      string    `gorm:"column:name;NOT NULL" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL"`
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL"`
}

type userDao struct {
	db *gorm.DB
}

func newUserDao(db *gorm.DB) *userDao {
	return &userDao{db: db}
}

// Create 创建一条记录，插入记录后，id值被回写到table中
func (d *userDao) Create(ctx context.Context, table *User) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID 根据id删除一条记录
func (d *userDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		return nil
	}

	return nil
}

// UpdateByID 根据id更新记录
func (d *userDao) UpdateByID(ctx context.Context, table *User) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}
	if table.Name != "" {
		update["name"] = table.Name
	}

	err := d.db.WithContext(ctx).Model(table).Updates(update).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByID 根据id获取一条记录
func (d *userDao) GetByID(ctx context.Context, id uint64) (*User, error) {
	table := &User{}
	err := d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
	return table, err
}

func (d *userDao) GetByColumns(ctx context.Context, params *query.Params) ([]*User, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if params.Sort != "ignore count" {
		err = d.db.WithContext(ctx).Model(&User{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*User{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// ---------------------------------------测试示例---------------------------------------------------

func newTestUserDao() *Dao {
	now := time.Now()
	testData := &User{
		ID:        1,
		Name:      "foo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 初始化mock cache
	c := NewCache(map[string]interface{}{"no cache": testData}) // 为了测试mysql，禁止缓存

	// 初始化mock dao
	d := NewDao(c, testData)
	d.IDao = newUserDao(d.DB)

	return d
}

func TestUserDao_Create(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	d.SqlMock.ExpectBegin()
	d.SqlMock.ExpectExec("INSERT INTO .*").
		WithArgs(testData.Name, testData.CreatedAt, testData.UpdatedAt, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SqlMock.ExpectCommit()

	err := d.IDao.(*userDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_Update(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	d.SqlMock.ExpectBegin()
	d.SqlMock.ExpectExec("UPDATE .*").
		WithArgs(testData.Name, d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SqlMock.ExpectCommit()

	err := d.IDao.(*userDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_Delete(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	d.SqlMock.ExpectBegin()
	d.SqlMock.ExpectExec("DELETE  FROM .*").
		WithArgs(testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SqlMock.ExpectCommit()

	err := d.IDao.(*userDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_GetByID(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(testData.ID, testData.Name)

	d.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(*userDao).GetByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = d.SqlMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_GetByColumns(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(testData.ID, testData.Name)
	d.SqlMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(*userDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Sort: "ignore count",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SqlMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}
