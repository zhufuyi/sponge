package gotest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
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

// ----------- Example of adding, deleting and checking --------------------

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

func (d *userDao) Create(ctx context.Context, table *User) error {
	return d.db.WithContext(ctx).Create(table).Error
}

func (d *userDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		return nil
	}

	return nil
}

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

// --------------------------------------- Test example --------------------------------------------

func newTestUserDao() *Dao {
	now := time.Now()
	testData := &User{
		ID:        1,
		Name:      "foo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// init mock cache, to test mysql, disable caching
	c := NewCache(map[string]interface{}{"no cache": testData})

	// init mock dao
	d := NewDao(c, testData)
	d.IDao = newUserDao(d.DB)

	return d
}

func TestUserDao_Create(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(testData.Name, testData.CreatedAt, testData.UpdatedAt, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(*userDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_Update(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(testData.Name, d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(*userDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_Delete(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("DELETE  FROM .*").
		WithArgs(testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

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

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(*userDao).GetByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDao_GetByColumns(t *testing.T) {
	d := newTestUserDao()
	defer d.Close()
	testData := d.TestData.(*User)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(testData.ID, testData.Name)
	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(*userDao).GetByColumns(d.Ctx, &query.Params{
		Page:  0,
		Limit: 10,
		Sort:  "ignore count",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}
