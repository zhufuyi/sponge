## cache

内存类型的有：memory和NoSQL

NoSQL的主要有: redis

各类库只要实现了cache定义的接口(driver)即可。
> 这里的接口driver命名参考了Go官方mysql接口的命名规范

## 多级缓存

### 二级缓存

这里的多级主要是指二级缓存：本地缓存(一级缓存L1)+redis缓存(二级缓存L2)  
使用本地缓存可以减少应用服务器到redis之间的网络I/O开销

> 需要注意的是：在并发量不大的系统内，本地缓存的意义不大，反而增加维护的困难。但在高并发系统中，
> 本地缓存可以大大节约带宽。但是要注意本地缓存不是银弹，它会引起多个副本间数据的
> 不一致，还会占据大量的内存，所以不适合保存特别大的数据，而且需要严格考虑刷新机制。

### 过期时间

本地缓存过期时间比分布式缓存小至少一半，以防止本地缓存太久造成多实例数据不一致。

<br>

## 使用示例

```go
// GetByID 根据id获取一条记录
func (d *userExampleDao) GetByID(ctx context.Context, id uint64) (*model.UserExample, error) {
	record, err := d.cache.Get(ctx, id)

	if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// 从mysql获取
	if errors.Is(err, goredis.ErrRedisNotFound) {
		table := &model.UserExample{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
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

		if table.ID == 0 {
			return nil, model.ErrRecordNotFound
		}

		// set cache
		err = d.cache.Set(ctx, id, table, cacheBase.DefaultExpireTime)
		if err != nil {
			return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
		}

		return table, nil
	}

	if err != nil {
		// fail fast, if cache error return, don't request to db
		return nil, err
	}

	return record, nil
}
```

<br>

## 缓存问题

需要注意以下几个问题：

- 缓存穿透
- 缓存击穿
- 缓存雪崩

可以参考文章：[Redis缓存三大问题](https://mp.weixin.qq.com/s/HjzwefprYSGraU1aJcJ25g)

## Reference
- ristretto：https://github.com/dgraph-io/ristretto (号称最快的本地缓存)
- [Ristretto简介：高性能Go缓存](https://www.yuque.com/kshare/2020/ade1d9b5-5925-426a-9566-3a5587af2181)
- bigcache: https://github.com/allegro/bigcache
- freecache: https://github.com/coocood/freecache
- concurrent_map: https://github.com/easierway/concurrent_map
- gocache: https://github.com/eko/gocache (Built-in stores, eg: bigcache,memcache,redis)