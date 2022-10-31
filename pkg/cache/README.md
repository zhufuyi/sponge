## cache

支持memory缓存和redis缓存。

## 使用示例

```go

// 实例化缓存类型，根据CType字段选择使用memory或redis
cache := cache.NewUserExampleCache(&model.CacheType{
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

	if errors.Is(err, model.ErrCacheNotFound) {
		// 从mysql获取
		table := &model.UserExample{}
		err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
		if err != nil {
			// if data is empty, set not found cache to prevent cache penetration(防止缓存穿透)
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
```