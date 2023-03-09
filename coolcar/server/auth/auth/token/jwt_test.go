package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAqSlpNaGk3oxOzjM07yw5ysHHSVhTTuf6fDZiohdMdlPyFIpx
ZPMEBXBpjAMOkvFEVF3YjTfAN57CXxUcBNfncnoXf4GGOsPpDWBsUgeYpJpFQpu5
XhH+z2GT7dqvrSey0EqRHTKCTWINMPvlg3r7xuEDzWK2KFO1MB3j/S56gDEHZbOx
UE4xieLULXYh0erVH4TAHCzQqwG6rJ7ty0gnzN0VxiWTvyklyYSHm5QU0jQo1xSM
m4Awb1JiJ3iG1bga7NdEMOvXHtK/cmisBDYiu3/tn8jjY8o720M9K364PTI11LET
Px4AyKwlA8Z0nk6Obnl8ZFY3hPzdmXBhYUmRcwIDAQABAoIBAAxdIHlxBRHXA3OC
vp20h8zP+PbdY8linyYw+2iJd2c2n3zs2XjdYi/blXtMReZrh+j6qvc3We82xVIZ
wuB/v0TYs5r4Jo1pAEGgCIq+T2PIesNxikzb19nkceFymGB5hFJBAPY3WNq7DefE
oXPIq1pP2+1JB7NO2vFXXfCFUyLHVASr8FY2o3mX9oMxh8qy74UDPcz/E+k36D1i
1WmGzY6LF8GG6UOEi4DCE4XpIkM7Z+4CuNFkusZPjuk4FVQxzuR+lWBqLuaU2oaD
Bj5LT2VpTulnc8tbHBk6TtWHU7fJdq4XIctpiGrEtnD6QjtDZkLsF1ilf3kgNeco
XB9j1jECgYEA5o8JgcCAX8c843HBFZQkwqaRsb87QSuPZjWqwbfXRsRY3FcxAW0Q
R6GXj0ZWxkVh63AS6TkuF+LxCS0ZP+5tc6Rsq/g+FdYNbe/1cmIawzWInlFK4Oou
eKwdUuh4syZeGQZIMbpKURswpc3r/mDt5W1UBmh54UYp8etEN2N8iu8CgYEAu9P+
rm5nAis3TeECYg+qYEppci9cnXwJLgWzWsRwLgdOVgvwrmjf8rdJahBeHNO049w6
rD/tLvCEjkY14aoXEx45HJqc5zBpww9atQUcPqtyUnmn9AB0WXPbMjUDLJQliYle
hNETMAkiEqCZMdsqfh/pX+zZ3z6FKDan+vk08b0CgYB9ntnTNIu9o9TtKAHIPBt7
Yz5m1ob2j0Fmsz8CpaRKDplMFMXCvSXtoYHusqh9Bzi/CyWCpYETyrcCBOyJBOPl
6mS7nlpVk3dluyTE2eczDWwOtsRRn8cKQN0JW1jIY9NJVz7muVXcsy/iZzx6MV3t
b5AknbAqqgwYn9NfSnmFSQKBgCWRPTs+MbQpWKJnAscCQx2HRJfmCSwmht+BnGHn
MFjEdVKYiMcZitFM/44LQAecAG4iukmBb7sXuCuMt3IvRvY38UxbUE6dTEoLZCUY
pJGUUQVV99XB0YOivJDKMZxU9T0REKqX9rKA4SPAo2NpZpJbZ54cDWetZYypgeec
uI4xAoGALi2kFnj2d3HRJtL7+b+a3tGECRJ9fUYFcJIA30kYXy0wjYUTWJcqaYk+
f9QaX9gKKecHGuPCpBu5SxV3Vake136zVxh1JjdwdoOu4zGzCPN5gVpWTTDrpwUg
SMgG/jLJATX8Z+SLEL+bpsUEafhqVzZBSiM/ZzWYlWR06k37rAo=
-----END RSA PRIVATE KEY-----`

	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSlpNaGk3oxOzjM07yw5
ysHHSVhTTuf6fDZiohdMdlPyFIpxZPMEBXBpjAMOkvFEVF3YjTfAN57CXxUcBNfn
cnoXf4GGOsPpDWBsUgeYpJpFQpu5XhH+z2GT7dqvrSey0EqRHTKCTWINMPvlg3r7
xuEDzWK2KFO1MB3j/S56gDEHZbOxUE4xieLULXYh0erVH4TAHCzQqwG6rJ7ty0gn
zN0VxiWTvyklyYSHm5QU0jQo1xSMm4Awb1JiJ3iG1bga7NdEMOvXHtK/cmisBDYi
u3/tn8jjY8o720M9K364PTI11LETPx4AyKwlA8Z0nk6Obnl8ZFY3hPzdmXBhYUmR
cwIDAQAB
-----END PUBLIC KEY-----`
)

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key: %v", err)
	}
	g := NewJWTTokenGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1676895340, 0)
	}
	tkn, err := g.GenerateToken("63e49d966d5cbaf02d5c0d21", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate togken: %v", err)
	}

	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY5MDI1NDAsImlhdCI6MTY3Njg5NTM0MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjNlNDlkOTY2ZDVjYmFmMDJkNWMwZDIxIn0.l2avoJexiFgjpOwG0kaxtgnN73SvfNLSa9ZgwzNxclh82sepnQrX5paLodoPoQYKdYYF4N36Q4WPfsM_h9fUwOntI9fwvcmrLgpoY3C98PghE7MDX3rOCVJLX4QMYWCalx-P8IgGPOxKxMfdZhjbSWtMEoxWxgnLIQfOdP2z-o52UcNOnvi5xDajuU23oSoROAJbsnxBrfeUcC38krlcEUJv7X6jLyQEojDyTiVS_SgGl1ja8eah3G-0_EZcoZAwG-SqiAo1NNIgbNlHfrPX3jyQ6MgCeA3hs-ZAt0V9Wr4VbrT_zTwDALBRrul_-aOSWG4Zru8mCNs6Nfrf7Hh1Dw"
	if tkn != want {
		t.Errorf("wrong token generated. want: %q; got %q", want, tkn)
	}
}
