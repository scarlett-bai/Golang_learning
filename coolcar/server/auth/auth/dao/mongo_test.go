package dao

import (
	"context"
	"os"
	"testing"

	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"

	"go.mongodb.org/mongo-driver/bson"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {
	// start container
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("63e49d966d5cbaf02d5c0d10")),
			openIDField:        "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("63e49d966d5cbaf02d5c0d22")),
			openIDField:        "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}
	// mgutil.NewObjID = func() primitive.ObjectID {
	// 	return objid.MustFromID(id.AccountID("63e49d966d5cbaf02d5c0d21"))
	// }

	mgutil.NewObjectIDWithValue(id.AccountID("63e49d966d5cbaf02d5c0d21"))
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "63e49d966d5cbaf02d5c0d10",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "63e49d966d5cbaf02d5c0d22",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "63e49d966d5cbaf02d5c0d21",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("cannot resolve account id for %q: %v", cc.openID, err)
			}
			if id.String() != cc.want {
				t.Errorf("resolve account id: wamt:%q; got %q", cc.want, id)
			}

		})
	}

	// id, err := m.ResolveAccountID(c, "123")
	// if err != nil {
	// 	t.Errorf("failed resolve account id for 123: %v", err)
	// } else {
	// 	want := "63e49d966d5cbaf02d5c0d10"
	// 	if id != want {
	// 		t.Errorf("resolve account id: want : %q, got: %q", want, id)
	// 	}
	// 	// remove container
	// }
}

func MustFromID(s string) {
	panic("unimplemented")
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
