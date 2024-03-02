package mgo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model embedded structs, add `bson: ",inline"` when defining table structs
type Model struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

// SetModelValue set model fields
func (p *Model) SetModelValue() {
	now := time.Now()
	if !p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}

	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
		p.UpdatedAt = now
	}
}

// ExcludeDeleted exclude soft deleted records
func ExcludeDeleted(filter bson.M) bson.M {
	if filter == nil {
		filter = bson.M{}
	}
	filter["deleted_at"] = bson.M{"$exists": false}
	return filter
}

// EmbedUpdatedAt embed updated_at datetime column
func EmbedUpdatedAt(update bson.M) bson.M {
	updateM := bson.M{}
	if v, ok := update["$set"]; ok {
		if m, ok2 := v.(bson.M); ok2 {
			m["updated_at"] = time.Now()
			updateM["$set"] = m
		}
	} else {
		update["updated_at"] = time.Now()
		updateM["$set"] = update
	}
	return updateM
}

// EmbedDeletedAt embed deleted_at datetime column
func EmbedDeletedAt(update bson.M) bson.M {
	updateM := bson.M{}
	if v, ok := update["$set"]; ok {
		if m, ok2 := v.(bson.M); ok2 {
			m["deleted_at"] = time.Now()
			updateM["$set"] = m
		}
	} else {
		updateM["$set"] = bson.M{"deleted_at": time.Now()}
	}
	return updateM
}

// ConvertToObjectIDs convert ids to objectIDs
func ConvertToObjectIDs(ids []string) []primitive.ObjectID {
	oids := []primitive.ObjectID{}
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		oids = append(oids, oid)
	}
	return oids
}
