package cos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Service
type Service struct {
	client *cos.Client
	secID  string
	secKey string
}

func NewService(addr, secID, secKey string) (*Service, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse addr: %v", err)
	}
	// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥

	return &Service{
		client: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  secID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
				SecretKey: secKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			},
		}),
		secID:  secID,
		secKey: secKey,
	}, nil

}

func (s *Service) SignURL(c context.Context, method, path string, timeout time.Duration) (string, error) {
	u, err := s.client.Object.GetPresignedURL(c, method, path, s.secID, s.secKey, timeout, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *Service) Get(c context.Context, path string) (io.ReadCloser, error) {
	res, err := s.client.Object.Get(c, path, nil)
	var b io.ReadCloser
	if res != nil {
		b = res.Body
	}
	if err != nil {
		return b, err
	}
	if res.StatusCode >= 400 {
		return b, fmt.Errorf("got err response: %+v", res)
	}
	return b, nil
}
