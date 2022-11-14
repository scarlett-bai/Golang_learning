package trip

import (
	"context"
	trippb "coolcar/proto/gen/go"
)

// type TripServiceServer interface {
// 	GetTrip(context.Context, *GetTripRequest) (*GetTripResponse, error)
// }

// Service implements trip service.
type Service struct{}

func (*Service) GetTrip(c context.Context, req *trippb.GetTripRequest) (*trippb.GetTripResponse, error) {
	return &trippb.GetTripResponse{
		Id: req.Id,
		Trip: &trippb.Trip{
			Start:       "Shanghai",
			End:         "Hangzhou",
			DurationSec: 36000,
			FeeCent:     10000,
			StartPos: &trippb.Location{
				Latitude:  30,
				Longitude: 120,
			},
			EndPos: &trippb.Location{
				Latitude:  35,
				Longitude: 115,
			},
			PathLocations: []*trippb.Location{
				{
					Latitude:  31,
					Longitude: 119,
				},
				{
					Latitude:  32,
					Longitude: 118,
				},
			},
			Status: trippb.TripStatus_IN_PROGRESS,
		},
	}, nil

}
