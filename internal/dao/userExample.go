package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/model"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var _ UserExampleDao = (*userExampleDao)(nil)

// UserExampleDao defining the dao interface
type UserExampleDao interface {
	Create(ctx context.Context, table *model.UserExample) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.UserExample) error
	GetByID(ctx context.Context, id uint64) (*model.UserExample, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]*model.UserExample, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.UserExample, int64, error)
}

type userExampleDao struct {
	db    *gorm.DB
	cache cache.UserExampleCache
}

// NewUserExampleDao creating the dao interface
func NewUserExampleDao(db *gorm.DB, cache cache.UserExampleCache) UserExampleDao {
	return &userExampleDao{db: db, cache: cache}
}

// Create a record, insert the record and the id value is written back to the table
func (d *userExampleDao) Create(ctx context.Context, table *model.UserExample) error {
	err := d.db.WithContext(ctx).Create(table).Error
	_ = d.cache.Del(ctx, table.ID)
	return err
}

// DeleteByID delete a record based on id
func (d *userExampleDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserExample{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.cache.Del(ctx, id)

	return nil
}

// UpdateByID update records by id
func (d *userExampleDao) UpdateByID(ctx context.Context, table *model.UserExample) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}
	// todo generate the update fields code to here
	// delete the templates code start
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
	// delete the templates code end

	err := d.db.WithContext(ctx).Model(table).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.cache.Del(ctx, table.ID)

	return nil
}

// GetByID get a record based on id
func (d *userExampleDao) GetByID(ctx context.Context, id uint64) (*model.UserExample, error) {
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// get from mysql
		table := &model.UserExample{}
		err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
		if err != nil {
			// if data is empty, set not found cache to prevent cache penetration(preventing Cache Penetration)
			if errors.Is(err, model.ErrRecordNotFound) {
				err = d.cache.SetCacheWithNotFound(ctx, id)
				if err != nil {
					return nil, err
				}
				return nil, model.ErrRecordNotFound
			}
			return nil, err
		}

		// set cache
		err = d.cache.Set(ctx, id, table, cacheBase.DefaultExpireTime)
		if err != nil {
			return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByIDs get multiple rows by ids
func (d *userExampleDao) GetByIDs(ctx context.Context, ids []uint64) ([]*model.UserExample, error) {
	records := []*model.UserExample{}

	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []uint64
	for _, id := range ids {
		item, ok := itemMap[cast.ToString(id)]
		if !ok {
			missedIDs = append(missedIDs, id)
			continue
		}
		records = append(records, item)
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in mysql
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if errors.Is(err, cacheBase.ErrPlaceholder) {
				continue
			} else {
				realMissedIDs = append(realMissedIDs, id)
			}
		}

		if len(realMissedIDs) > 0 {
			var missedData []*model.UserExample
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				records = append(records, missedData...)
				err = d.cache.MultiSet(ctx, missedData, 10*time.Minute)
				if err != nil {
					return nil, err
				}
			} else {
				for _, id := range realMissedIDs {
					_ = d.cache.SetCacheWithNotFound(ctx, id)
				}
			}
		}
	}
	return records, nil
}

// GetByColumns filter multiple rows based on paging and column information
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	size: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default to = when the value is null, have =, ! =, >, >=, <, <=, like
//	value: column name
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
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
//			Value: "male",
//		},
//	}
func (d *userExampleDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.UserExample, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.UserExample{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.UserExample{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}
