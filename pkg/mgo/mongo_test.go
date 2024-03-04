package mgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mgoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

func TestInit(t *testing.T) {
	dsns := []string{
		"mongodb://root:123456@192.168.3.37:27017/account",
		"mongodb://root:123456@192.168.3.37:27017/account?connectTimeoutMS=2000",
		// error
		"mongodb-dsn",
		"mongodb://root:123456@192.168.3.37",
	}
	for _, dsn := range dsns {
		timeout := time.Second * 2
		opts := &mgoOptions.ClientOptions{Timeout: &timeout}
		db, err := Init(dsn, opts)
		if err != nil {
			t.Log(err)
			continue
		}
		time.Sleep(time.Millisecond * 100)
		defer Close(db)
	}

	defer func() { recover() }()
	db := &mongo.Database{}
	_ = Close(db)
}

func TestInit2(t *testing.T) {
	uri := "mongodb://root:123456@192.168.3.37:27017"
	dbName := "account"
	timeout := time.Second * 3
	opts := &mgoOptions.ClientOptions{Timeout: &timeout}
	db, err := Init2(uri, dbName, opts)
	if err != nil {
		t.Log(err)
		return
	}
	time.Sleep(time.Millisecond * 100)
	defer Close(db)
}

func TestModel_SetModelValue(t *testing.T) {
	m := new(Model)
	m.SetModelValue()

	assert.NotNil(t, m.ID)
	assert.NotNil(t, m.CreatedAt)
	assert.NotNil(t, m.UpdatedAt)
}

func TestExcludeDeleted(t *testing.T) {
	filter := bson.M{"foo": "bar"}
	filter = ExcludeDeleted(filter)
	assert.NotNil(t, filter["deleted_at"])

	filter = ExcludeDeleted(nil)
	assert.NotNil(t, filter["deleted_at"])
}

func TestEmbedUpdatedAt(t *testing.T) {
	update := bson.M{"$set": bson.M{"foo": "bar"}}
	update = EmbedUpdatedAt(update)
	m := update["$set"].(bson.M)
	assert.NotNil(t, m["updated_at"])

	update = bson.M{"foo": "bar"}
	update = EmbedUpdatedAt(update)
	m = update["$set"].(bson.M)
	assert.NotNil(t, m["updated_at"])
}

func TestEmbedDeletedAt(t *testing.T) {
	update := bson.M{"$set": bson.M{"foo": "bar"}}
	update = EmbedDeletedAt(update)
	m := update["$set"].(bson.M)
	assert.NotNil(t, m["deleted_at"])

	update = bson.M{"foo": "bar"}
	update = EmbedDeletedAt(update)
	m = update["$set"].(bson.M)
	assert.NotNil(t, m["deleted_at"])
}

func TestConvertToObjectIDs(t *testing.T) {
	ids := []string{"65c9ae1b1378ae7f0787a039", "invalid_id"}
	oids := ConvertToObjectIDs(ids)
	assert.Equal(t, len(oids), 1)
}
