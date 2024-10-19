## mgo

`mgo` is a library wrapped on the official library [mongo-go-driver](https://github.com/mongodb/mongo-go-driver), with added features paging queries.

<br>

### Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/mgo"

    // dsn document: https://www.mongodb.com/docs/manual/reference/connection-string/

    // case 1: specify options in dsn
    db, err := mgo.Init("mongodb://root:123456@192.168.3.37:27017/account?socketTimeoutMS=30000&maxPoolSize=100&minPoolSize=1&maxConnIdleTimeMS=300000")
    // case 2: specify options in code
    db, err := mgo.Init("mongodb://root:123456@192.168.3.37:27017/account",
        mgo.WithOption().SetMaxPoolSize(100),
        mgo.WithOption().SetMinPoolSize(1),
        mgo.WithOption().SetMaxConnIdleTime(5*time.Minute),
        mgo.WithOption().SetSocketTimeout(30*time.Second),
    )

    // close mongodb
    defer mgo.Close(db)
```
