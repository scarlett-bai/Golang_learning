package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// 既然是个实现，就要写清楚它的内容
type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify verifiers a token and returns acount id.
func (v *JWTTokenVerifier) Verify(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %v", err)
	}
	if !t.Valid {
		return "", fmt.Errorf("token not valid")
	}
	clm, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not standard claims")
	}
	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valied: %v", err)
	}

	return clm.Subject, nil
}
