package googleauth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testClientID = "test-client.apps.googleusercontent.com"

// testServer serves a JWKS containing priv's public key under kid, so a
// verifier pointed at it can validate tokens signed by priv.
func testServer(t *testing.T, kid string, priv *rsa.PrivateKey) *httptest.Server {
	t.Helper()
	n := base64.RawURLEncoding.EncodeToString(priv.PublicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(priv.PublicKey.E)).Bytes())

	body := fmt.Sprintf(`{"keys":[{"kid":%q,"kty":"RSA","n":%q,"e":%q}]}`, kid, n, e)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func signToken(t *testing.T, priv *rsa.PrivateKey, kid string, claims tokenClaims) string {
	t.Helper()
	header := map[string]string{"alg": "RS256", "kid": kid, "typ": "JWT"}

	headerJSON, err := json.Marshal(header)
	require.NoError(t, err)
	claimsJSON, err := json.Marshal(claims)
	require.NoError(t, err)

	headerSeg := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsSeg := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := headerSeg + "." + claimsSeg
	hashed := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
	require.NoError(t, err)

	return signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func validClaims() tokenClaims {
	return tokenClaims{
		Iss:           issuerHTTPS,
		Aud:           testClientID,
		Exp:           time.Now().Add(time.Hour).Unix(),
		Sub:           "1234567890",
		Email:         "sibling@example.com",
		EmailVerified: true,
		Name:          "Sibling Person",
	}
}

func TestVerifySuccess(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv)
	v := newVerifier(testClientID, srv.URL)

	token := signToken(t, priv, "key-1", validClaims())

	claims, err := v.Verify(context.Background(), token)

	require.NoError(t, err)
	assert.Equal(t, "1234567890", claims.Subject)
	assert.Equal(t, "sibling@example.com", claims.Email)
	assert.True(t, claims.EmailVerified)
	assert.Equal(t, "Sibling Person", claims.Name)
}

func TestVerifyBareIssuerAccepted(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv)
	v := newVerifier(testClientID, srv.URL)

	claims := validClaims()
	claims.Iss = issuerBare
	token := signToken(t, priv, "key-1", claims)

	_, err = v.Verify(context.Background(), token)

	require.NoError(t, err)
}

func TestVerifyWrongAudienceRejected(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv)
	v := newVerifier(testClientID, srv.URL)

	claims := validClaims()
	claims.Aud = "someone-elses-client-id"
	token := signToken(t, priv, "key-1", claims)

	_, err = v.Verify(context.Background(), token)

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyExpiredRejected(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv)
	v := newVerifier(testClientID, srv.URL)

	claims := validClaims()
	claims.Exp = time.Now().Add(-time.Hour).Unix()
	token := signToken(t, priv, "key-1", claims)

	_, err = v.Verify(context.Background(), token)

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyWrongIssuerRejected(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv)
	v := newVerifier(testClientID, srv.URL)

	claims := validClaims()
	claims.Iss = "https://not-google.example.com"
	token := signToken(t, priv, "key-1", claims)

	_, err = v.Verify(context.Background(), token)

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifySignedByWrongKeyRejected(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv) // JWKS only knows priv's key

	impostor, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	v := newVerifier(testClientID, srv.URL)

	// Token claims to be "key-1" but is actually signed by a different key.
	token := signToken(t, impostor, "key-1", validClaims())

	_, err = v.Verify(context.Background(), token)

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyUnknownKidRejected(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	srv := testServer(t, "key-1", priv)
	v := newVerifier(testClientID, srv.URL)

	token := signToken(t, priv, "key-does-not-exist", validClaims())

	_, err = v.Verify(context.Background(), token)

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyMalformedTokenRejected(t *testing.T) {
	v := newVerifier(testClientID, "http://unused.invalid")

	_, err := v.Verify(context.Background(), "not-a-jwt")

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyCachesKeysAcrossCalls(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	fetches := 0
	n := base64.RawURLEncoding.EncodeToString(priv.PublicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(priv.PublicKey.E)).Bytes())
	body := fmt.Sprintf(`{"keys":[{"kid":"key-1","kty":"RSA","n":%q,"e":%q}]}`, n, e)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetches++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	v := newVerifier(testClientID, srv.URL)
	token := signToken(t, priv, "key-1", validClaims())

	_, err = v.Verify(context.Background(), token)
	require.NoError(t, err)
	_, err = v.Verify(context.Background(), token)
	require.NoError(t, err)

	assert.Equal(t, 1, fetches, "second verify should reuse the cached key, not refetch")
}
