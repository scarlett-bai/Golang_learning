package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSlpNaGk3oxOzjM07yw5
ysHHSVhTTuf6fDZiohdMdlPyFIpxZPMEBXBpjAMOkvFEVF3YjTfAN57CXxUcBNfn
cnoXf4GGOsPpDWBsUgeYpJpFQpu5XhH+z2GT7dqvrSey0EqRHTKCTWINMPvlg3r7
xuEDzWK2KFO1MB3j/S56gDEHZbOxUE4xieLULXYh0erVH4TAHCzQqwG6rJ7ty0gn
zN0VxiWTvyklyYSHm5QU0jQo1xSMm4Awb1JiJ3iG1bga7NdEMOvXHtK/cmisBDYi
u3/tn8jjY8o720M9K364PTI11LETPx4AyKwlA8Z0nk6Obnl8ZFY3hPzdmXBhYUmR
cwIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}
	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name: "valid_token",
			tkn:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY5MDI1NDAsImlhdCI6MTY3Njg5NTM0MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjNlNDlkOTY2ZDVjYmFmMDJkNWMwZDIxIn0.l2avoJexiFgjpOwG0kaxtgnN73SvfNLSa9ZgwzNxclh82sepnQrX5paLodoPoQYKdYYF4N36Q4WPfsM_h9fUwOntI9fwvcmrLgpoY3C98PghE7MDX3rOCVJLX4QMYWCalx-P8IgGPOxKxMfdZhjbSWtMEoxWxgnLIQfOdP2z-o52UcNOnvi5xDajuU23oSoROAJbsnxBrfeUcC38krlcEUJv7X6jLyQEojDyTiVS_SgGl1ja8eah3G-0_EZcoZAwG-SqiAo1NNIgbNlHfrPX3jyQ6MgCeA3hs-ZAt0V9Wr4VbrT_zTwDALBRrul_-aOSWG4Zru8mCNs6Nfrf7Hh1Dw",
			now:  time.Unix(1676895340, 0),
			want: "63e49d966d5cbaf02d5c0d21",
			// wantErr: false,
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY5MDI1NDAsImlhdCI6MTY3Njg5NTM0MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjNlNDlkOTY2ZDVjYmFmMDJkNWMwZDIxIn0.l2avoJexiFgjpOwG0kaxtgnN73SvfNLSa9ZgwzNxclh82sepnQrX5paLodoPoQYKdYYF4N36Q4WPfsM_h9fUwOntI9fwvcmrLgpoY3C98PghE7MDX3rOCVJLX4QMYWCalx-P8IgGPOxKxMfdZhjbSWtMEoxWxgnLIQfOdP2z-o52UcNOnvi5xDajuU23oSoROAJbsnxBrfeUcC38krlcEUJv7X6jLyQEojDyTiVS_SgGl1ja8eah3G-0_EZcoZAwG-SqiAo1NNIgbNlHfrPX3jyQ6MgCeA3hs-ZAt0V9Wr4VbrT_zTwDALBRrul_-aOSWG4Zru8mCNs6Nfrf7Hh1Dw",
			now:     time.Unix(1686895340, 0),
			wantErr: true,
			// want: "",
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1676895340, 0),
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.tkn)
			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v", err)
			}
			if c.wantErr && err == nil {
				t.Errorf("want error; got no error")
			}
			if accountID != c.want {
				t.Errorf("wrong account id. want: %q, got %q", c.want, accountID)
			}
		})
	}

	// tkn := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzY5MDI1NDAsImlhdCI6MTY3Njg5NTM0MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjNlNDlkOTY2ZDVjYmFmMDJkNWMwZDIxIn0.l2avoJexiFgjpOwG0kaxtgnN73SvfNLSa9ZgwzNxclh82sepnQrX5paLodoPoQYKdYYF4N36Q4WPfsM_h9fUwOntI9fwvcmrLgpoY3C98PghE7MDX3rOCVJLX4QMYWCalx-P8IgGPOxKxMfdZhjbSWtMEoxWxgnLIQfOdP2z-o52UcNOnvi5xDajuU23oSoROAJbsnxBrfeUcC38krlcEUJv7X6jLyQEojDyTiVS_SgGl1ja8eah3G-0_EZcoZAwG-SqiAo1NNIgbNlHfrPX3jyQ6MgCeA3hs-ZAt0V9Wr4VbrT_zTwDALBRrul_-aOSWG4Zru8mCNs6Nfrf7Hh1Dw"
	// jwt.TimeFunc = func() time.Time {
	// 	return time.Unix(1676895340, 0)
	// }
	// accountId, err := v.Verify(tkn)
	// if err != nil {
	// 	t.Errorf("verification failed: %v", err)
	// }
	// want := "63e49d966d5cbaf02d5c0d21"
	// if accountId != want {
	// 	t.Errorf("wrong account id. want: %q, got %q", want, accountId)
	// }
}
