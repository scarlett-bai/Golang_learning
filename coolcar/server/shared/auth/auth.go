package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"coolcar/shared/id"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor creates a grpc auth interceptor
func Interceptor(publikKeyFile string) (grpc.UnaryServerInterceptor, error) {
	f, err := os.Open(publikKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file: %v", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key:%v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	i := &interceptor{
		publicKey: pubKey,
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
	}
	return i.HandleReq, nil
}

type tokenVerifier interface {
	Verify(token string) (string, error)
}

type interceptor struct {
	publicKey *rsa.PublicKey
	verifier  tokenVerifier
}

func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	aid := impersonationFromContext(ctx)
	if aid != "" {
		fmt.Printf("impersonating %q\n", aid)
		return handler(ContextWithAccoundID(ctx, id.AccountID(aid)), req)
	}
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	aid, err = i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
	}

	return handler(ContextWithAccoundID(ctx, id.AccountID(aid)), req)
}

func impersonationFromContext(c context.Context) string {
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return ""
	}
	imp := m[ImpersonateAccountHeader]
	if len(imp) == 0 {
		return ""
	}
	return imp[0]
}

const (
	// ImpersonateAccountHeader defines the header for account id impersonation
	ImpersonateAccountHeader = "impresonate-account-id"
	authorizationHeader      = "authorization"
	bearerPrefix             = "Bearer "
)

func tokenFromContext(c context.Context) (string, error) {
	unauthenticated := status.Error(codes.Unauthenticated, "")
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", unauthenticated
	}
	// fmt.Println(m)
	tkn := ""
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if tkn == "" {
		return "", unauthenticated
	}

	return tkn, nil
}

type accountIDKey struct{}

// ContextWithAccountID creates a context with given account
func ContextWithAccoundID(c context.Context, aid id.AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}

// AccountIDFromContext gets account id from context
//
//	Returns unauthenticated error if no account id is available.
func AccountIDFromContext(c context.Context) (id.AccountID, error) {
	v := c.Value(accountIDKey{})
	aid, ok := v.(id.AccountID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return aid, nil
}
