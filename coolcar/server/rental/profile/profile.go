package profile

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityReslover interface {
	Resolve(c context.Context, photo []byte) (*rentalpb.Identity, error)
}

type Service struct {
	BlobClient        blobpb.BlobServiceClient
	PhotoGetExpire    time.Duration
	PhotoUploadExpire time.Duration
	IdentityReslover  IdentityReslover
	Mongo             *dao.Mongo
	Logger            *zap.Logger
}

// GetProfile gets profile for the current account
func (s *Service) GetProfile(c context.Context, req *rentalpb.GetProfileRequest) (*rentalpb.Profile, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	pr, err := s.Mongo.GetProfile(c, aid)
	if err != nil {
		code := s.logAndConvertProfileErr(err)
		if code == codes.NotFound {
			// 当数据库里没有数据时， 返回一个空的结构体即可
			return &rentalpb.Profile{}, nil
		}
		return nil, status.Error(code, "")
	}
	if pr.Profile == nil {
		return &rentalpb.Profile{}, nil
	}
	return pr.Profile, nil
}

// SubmitProfile submits a profile
func (s *Service) SubmitProfile(c context.Context, i *rentalpb.Identity) (*rentalpb.Profile, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	p := &rentalpb.Profile{
		Identity:       i,
		IdentityStatus: rentalpb.IdentityStatus_PENDING,
	}
	err = s.Mongo.UpdateProfile(c, aid, rentalpb.IdentityStatus_UNSUBMITTED, p)
	if err != nil {
		s.Logger.Error("cannot update profile", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	// s.Logger.Info("here")
	// submit 之后 帮助用户完成审核
	go func() {
		// s.Logger.Info("start check the verify status")
		time.Sleep(3 * time.Second)
		err := s.Mongo.UpdateProfile(context.Background(), aid, rentalpb.IdentityStatus_PENDING, &rentalpb.Profile{
			Identity:       i,
			IdentityStatus: rentalpb.IdentityStatus_VERIFIED,
		})
		if err != nil {
			s.Logger.Error("cannot verify identity", zap.Error(err))
		}
	}()
	return p, nil
}

func (s *Service) ClearProfile(c context.Context, req *rentalpb.ClearProfileRequest) (*rentalpb.Profile, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	p := &rentalpb.Profile{}
	err = s.Mongo.UpdateProfile(c, aid, rentalpb.IdentityStatus_VERIFIED, p)
	if err != nil {
		s.Logger.Error("cannot update profile", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return p, nil
}

func (s *Service) ClearProfilePhoto(c context.Context, req *rentalpb.ClearProfilePhotoRequest) (*rentalpb.ClearProfilePhotoResponse, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	err = s.Mongo.UpdateProfilePhoto(c, aid, id.BlobID(""))
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &rentalpb.ClearProfilePhotoResponse{}, nil
}

// CompleteProfilePhoto completes profile photo, returns AI recoginiton results
func (s *Service) CompleteProfilePhoto(c context.Context, req *rentalpb.CompleteProfilePhotoRequest) (*rentalpb.Identity, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	pr, err := s.Mongo.GetProfile(c, aid)
	if err != nil {
		return nil, status.Error(s.logAndConvertProfileErr(err), "")
	}
	if pr.PhotoBlobID == "" {
		return nil, status.Error(codes.NotFound, "")
	}

	br, err := s.BlobClient.GetBlob(c, &blobpb.GetBlobRequest{
		Id: pr.PhotoBlobID,
	})
	if err != nil {
		s.Logger.Error("cannot get blob", zap.Error(err))
		return nil, status.Error(codes.Aborted, "")
	}

	s.Logger.Info("got profile photo", zap.Int("szie", len(br.Data)))
	return &rentalpb.Identity{
		LicNumber:       "22223445465",
		Name:            "落意",
		Gender:          rentalpb.Gender_FEMALE,
		BirthDateMillis: 955510075,
	}, nil
}

// CreateProfilePhoto creates profile photo
func (s *Service) CreateProfilePhoto(c context.Context, req *rentalpb.CreateProfilePhotoRequest) (*rentalpb.CreateProfilePhotoResponse, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	br, err := s.BlobClient.CreateBlob(c, &blobpb.CreateBlobRequest{
		AccountId:           aid.String(),
		UploadUrlTimeoutSec: int32(s.PhotoUploadExpire.Seconds()),
	})
	if err != nil {
		s.Logger.Error("cannot create blob", zap.Error(err))
		return nil, status.Error(codes.Aborted, "")
	}

	err = s.Mongo.UpdateProfilePhoto(c, aid, id.BlobID(br.Id))
	if err != nil {
		s.Logger.Error("cannot update profile photo", zap.Error(err))
		return nil, status.Error(codes.Aborted, "")
	}

	return &rentalpb.CreateProfilePhotoResponse{
		UploadUrl: br.UploadUrl,
	}, nil

}

// GetProfilePhoto gets profile photo
func (s *Service) GetProfilePhoto(c context.Context, req *rentalpb.GetProfilePhotoRequest) (*rentalpb.GetProfilePhotoResponse, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	pr, err := s.Mongo.GetProfile(c, aid)
	if err != nil {
		return nil, status.Error(s.logAndConvertProfileErr(err), "")
	}
	if pr.PhotoBlobID == "" {
		return nil, status.Error(codes.NotFound, "")
	}

	br, err := s.BlobClient.GetBlobURL(c, &blobpb.GetBlobURLRequest{
		Id:         pr.PhotoBlobID,
		TimeoutSec: int32(s.PhotoGetExpire.Seconds()),
	})
	if err != nil {
		s.Logger.Error("cannot get blob", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &rentalpb.GetProfilePhotoResponse{
		Url: br.Url,
	}, nil
}

func (s *Service) logAndConvertProfileErr(err error) codes.Code {
	if err == mongo.ErrNoDocuments {
		return codes.NotFound
	}
	s.Logger.Error("cannot get profile,", zap.Error(err))
	return codes.Internal
}
