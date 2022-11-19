## gotest

gotest是一个模拟测试cache、dao、handler的库。

<br>

## 使用示例

### 模拟测试cache

缓存接口示例

```go
// UserExampleCache cache interface
type UserExampleCache interface {
	Set(ctx context.Context, id uint64, data *model.UserExample, duration time.Duration) error
	Get(ctx context.Context, id uint64) (ret *model.UserExample, err error)
}

// userExampleCache define a cache struct
type userExampleCache struct {
	cache cache.Cache
}

// NewUserExampleCache new a cache
func NewUserExampleCache(rdb *redis.Client) UserExampleCache {
	return &userExampleCache{
		// ...
	}
}
```

测试缓存示例

```go
func newUserExampleCache() *gotest.RedisCache {
	testData := map[string]interface{}{
		"1": &model.UserExample{ID:1},
	}

	rc := gotest.NewRedisCache(testData)
	rc.ICache = NewUserExampleCache(rc.RedisClient)
	return rc
}

func Test_userExampleCache_Set(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserExample)
	err := c.ICache.(UserExampleCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_Get(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserExample)
	err := c.ICache.(UserExampleCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserExampleCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, record, got)
}
```

<br>

### 模拟测试dao

点击查看具体[示例](dao_test.go)。

<br>

### 模拟测试handler

```go
func newHandler() *Handler {
	now := time.Now()
	testData := &User{
		ID:        1,
		Name:      "foo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// init mock cache
	c := gotest.NewCache(map[string]interface{}{"no cache": testData})
	c.ICache = struct{}{} // instantiated cache interface

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = struct{}{} // instantiated dao interface

	// init mock handler
	h := gotest.NewHandler(d, testData)
	h.IHandler = struct{}{} // instantiated handler interface

	h := newHandler()
	defer h.Close()

	h.GoRunHttpServer([]gotest.RouterInfo{
		{
			FuncName: "GetByID",
			Method:   http.MethodGet,
			Path:     "/user/:id",
			HandlerFunc: func(c *gin.Context) {
				c.String(http.StatusOK, testData.Name)
			},
		},
	})

	return h
}

func TestGetHello(t *testing.T) {
	h := newHandler()
	defer h.Close()
	testData := h.TestData.(*User)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Get(result, h.GetRequestURL("GetByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}
```