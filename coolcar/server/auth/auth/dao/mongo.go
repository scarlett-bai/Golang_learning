package dao

import (
	"context"
	mgo "coolcar/shared/mongo"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const openIDField = "open_id" // 自然语言命名

// Mongo defines a mongo dao
type Mongo struct {
	col      *mongo.Collection
	NewObjID func() primitive.ObjectID
}

// NewMongo creates a new mongo dao.
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("account"),
		NewObjID: primitive.NewObjectID,
	}
}

func (m *Mongo) ResolveAccountID(c context.Context, openID string) (string, error) {
	// m.col.InsertOne(c, bson.M{
	// 	mgo.IDField: m.NewObjID(),
	// 	openIDField: openID,
	// })
	insertedID := m.NewObjID()
	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, mgo.SetOnInsert(bson.M{
		mgo.IDField: insertedID,
		openIDField: openID,
	}), options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))
	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v", err)
	}
	var row mgo.ObjID
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}
	return row.ID.Hex(), nil
}
