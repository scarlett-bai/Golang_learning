package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".status"
)

// Mongo defines a mongo dao
type Mongo struct {
	col *mongo.Collection // 表名？
}

// NewMongo creates a new mongo dao.
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
	}
}

// TripRecord defines a trip record in mongo db
type TripRecord struct {
	mgutil.IDField       `bson:"inline"` // Go的一个语法糖
	mgutil.UpdateAtField `bson:"inline"`
	Trip                 *rentalpb.Trip `bson:"trip"`
}

// TODO 表格驱动测试

// CreateTrip creates a trip
func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgutil.NewObjID()
	r.UpdateAt = mgutil.UpdatedAt()

	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTrip gets a trip.
func (m *Mongo) GetTrip(c context.Context, id id.TripID, accountID id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}
	res := m.col.FindOne(c, bson.M{
		mgutil.IDFieldName: objID,
		accountIDField:     accountID,
	})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var tr TripRecord // 这样写就会给TripRecord分配一个空间
	err = res.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode: %v", err)
	}
	return &tr, nil
}

// GetTrips gets trips for the account by status.
// if status is not specified, gets all trips for the account
func (m *Mongo) GetTrips(c context.Context, accountID id.AccountID, status rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}
	if status != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = status
	}
	res, err := m.col.Find(c, filter, options.Find().SetSort(bson.M{
		mgutil.IDFieldName: -1,
	}))
	if err != nil {
		return nil, err
	}

	var trips []*TripRecord
	for res.Next(c) {
		var trip TripRecord
		err := res.Decode(&trip)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

//	用updateAt 来实现乐观锁  同时更新的问题
//
// UpdateTrip updates a trip
func (m *Mongo) UpdateTrip(c context.Context, tid id.TripID, aid id.AccountID, updateAt int64, trip *rentalpb.Trip) error {
	objID, err := objid.FromID(tid)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}
	newUpdatedAt := mgutil.UpdatedAt()
	res, err := m.col.UpdateOne(c, bson.M{
		mgutil.IDFieldName:        objID,
		accountIDField:            aid.String(),
		mgutil.UpdatedAtFieldName: updateAt,
	}, mgutil.Set(bson.M{
		tripField:                 trip,
		mgutil.UpdatedAtFieldName: newUpdatedAt,
	}))

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
