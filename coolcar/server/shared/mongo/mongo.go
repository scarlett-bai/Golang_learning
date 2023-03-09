package mgo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

const IDField = "_id"

type ObjID struct {
	ID primitive.ObjectID `bson:"_id"`
}

// Set returns a $set update document
func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}
