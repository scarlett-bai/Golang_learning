package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"coolcar/shared/server"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background() // 建一个新的context 会因为没有用户信息而报错unauthorized 的错误
	// c := auth.ContextWithAccoundID(context.Background(), id.AccountID("account1"))

	pm := &profileManager{}
	cm := &carManager{}
	s := newService(c, t, pm, cm)

	nowFunc = func() int64 {
		return 1680233350
	}

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
	}
	pm.iID = "identity1"
	golden := `{"account_id":%q,"car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1680233350},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1680233350},"status":1,"identity_id":"identity1"}`
	cases := []struct {
		name         string
		accountID    string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:      "normal_create",
			accountID: "account1",
			tripID:    "5f8132eb12714bf629489054",
			want:      fmt.Sprintf(golden, "account1"),
		},
		{
			name:       "profile_err",
			accountID:  "account2",
			tripID:     "5f8132eb12714bf629489055",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:         "car_verify_err",
			accountID:    "account3",
			tripID:       "5f8132eb12714bf629489056",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name:         "car_unlock_err",
			accountID:    "account4",
			tripID:       "5f8132eb12714bf629489057",
			carUnlockErr: fmt.Errorf("unlock"),
			want:         fmt.Sprintf(golden, "account4"),
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjectIDWithValue(id.TripID(cc.tripID))
			pm.err = cc.profileErr
			cm.unlockErr = cc.carUnlockErr
			cm.verifyErr = cc.carVerifyErr
			c := auth.ContextWithAccoundID(context.Background(), id.AccountID(cc.accountID))
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			if res.Id != cc.tripID {
				t.Errorf("incorrect id; want %q, got %q", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshall response: %v", err)
			}
			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response: want %s, got %s", cc.want, got)
			}
		})
	}

}

func TestTripLifecycle(t *testing.T) {
	c := auth.ContextWithAccoundID(context.Background(), id.AccountID("account_for_lifecycle"))
	s := newService(c, t, &profileManager{}, &carManager{})
	tid := id.TripID("5f8132eb22714bf629489056")
	mgutil.NewObjectIDWithValue(tid)
	cases := []struct {
		name    string
		now     int64
		op      func() (*rentalpb.Trip, error)
		want    string
		wantErr bool
	}{
		{
			name: "create_trip",
			now:  10000,
			op: func() (*rentalpb.Trip, error) {
				e, err := s.CreateTrip(c, &rentalpb.CreateTripRequest{
					CarId: "car1",
					Start: &rentalpb.Location{
						Latitude:  32.123,
						Longitude: 114.2525,
					},
				})
				if err != nil {
					return nil, err
				}
				return e.Trip, nil
			},
			want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"status":1}`,
		},
		{
			name: "update_trip",
			now:  20000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
					Current: &rentalpb.Location{
						Longitude: 28.23,
						Latitude:  123.24,
					},
				})
			},
		},
		{
			name: "finish_trip",
			now:  30000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id:      tid.String(),
					EndTrip: true,
				})
			},
		},
		{
			name: "query_trip",
			now:  40000,
			op: func() (*rentalpb.Trip, error) {
				return s.GetTrip(c, &rentalpb.GetTripRequest{
					Id: tid.String(),
				})
			},
		},
		{
			name: "update_after_finished",
			now:  50000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
				})
			},
			wantErr: true,
		},
	}

	rand.Seed(1345)
	for _, cc := range cases {
		nowFunc = func() int64 {
			return cc.now
		}
		trip, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error; got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: operation failed: %v", cc.name, err)
			continue
		}
		b, err := json.Marshal(trip)
		if err != nil {
			t.Errorf("%s: failed marshalling response: %v", cc.name, err)
		}
		got := string(b)
		if cc.want != got {
			t.Errorf("%s: incoorect response;want: %s, got %s", cc.name, cc.want, got)
		}
	}
}

func newService(c context.Context, t *testing.T, pm ProfileManager, cm CarManager) *Service {
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo client: %v", err)
	}

	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("cannot create logger %v", err)
	}

	db := mc.Database("coolcar")
	mongotesting.SetupIndexs(c, db)

	return &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		DistanceCalc:   &distCalc{},
		Mongo:          dao.NewMongo(db),
		Logger:         logger,
	}
}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
	lockErr   error
}

// CarManager defines the ACL for car management
func (m *carManager) Verify(c context.Context, cid id.CarID, loc *rentalpb.Location) error {
	return m.verifyErr
}
func (m *carManager) Unlock(c context.Context, cid id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	return m.unlockErr
}

func (m *carManager) Lock(c context.Context, cid id.CarID) error {
	return m.lockErr
}

type distCalc struct {
}

func (*distCalc) DistanceKm(c context.Context, from *rentalpb.Location, to *rentalpb.Location) (float64, error) {
	if from.Latitude == to.Latitude && from.Longitude == from.Longitude {
		return 0, nil
	}
	return 100, nil
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
