## cache

memory and redis cache libraries.

## Example of use

```go

// Choose to create a memory or redis cache depending on CType
cache := cache.NewUserExampleCache(&database.CacheType{
  CType: "redis",
  Rdb:   c.RedisClient,
})

// -----------------------------------------------------------------------------------------

type userExampleDao struct {
	db    *gorm.DB
	cache cache.UserExampleCache
}

// NewUserExampleDao creating the dao interface
func NewUserExampleDao(db *gorm.DB, cache cache.UserExampleCache) UserExampleDao {
	return &userExampleDao{db: db, cache: cache}
}

// GetByID get a record based on id
func (d *userExampleDao) GetByID(ctx context.Context, id uint64) (*model.UserExample, error) {
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, database.ErrCacheNotFound) {
		// get from mysql
		table := &model.UserExample{}
		err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
		if err != nil {
			// if data is empty, set not found cache to prevent cache penetration(preventing Cache Penetration)
			if errors.Is(err, database.ErrRecordNotFound) {
				err = d.cache.SetCacheWithNotFound(ctx, id)
				if err != nil {
					return nil, err
				}
				return nil, database.ErrRecordNotFound
			}
			return nil, err
		}

		// set cache
		err = d.cache.Set(ctx, id, table, 10*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, database.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}
```
