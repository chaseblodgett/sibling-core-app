// Package googleauth verifies "Sign in with Google" ID tokens. It checks
// the token's RS256 signature against Google's published public keys and
// validates the standard claims (issuer, audience, expiry) itself, rather
// than calling Google's tokeninfo debug endpoint on every sign-in — Google
// recommends local signature verification for production use.
// https://developers.google.com/identity/sign-in/web/backend-auth
package googleauth

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultCertsURL = "https://www.googleapis.com/oauth2/v3/certs"

	// Google documents both forms as valid issuer values.
	issuerHTTPS = "https://accounts.google.com"
	issuerBare  = "accounts.google.com"

	// How long a fetched key is trusted before we'll bother re-fetching.
	// Google rotates signing keys infrequently (on the order of weeks), so
	// this only governs how quickly we'd notice a rotation, not security.
	keysCacheTTL = time.Hour
)

// ErrInvalidToken covers every way an ID token can fail verification —
// malformed, bad signature, wrong audience/issuer, or expired. Callers
// don't need to distinguish the reason; it's always "reject this sign-in".
var ErrInvalidToken = errors.New("invalid google id token")

// Claims is the subset of a verified Google ID token's claims this app
// uses to identify the signed-in user.
type Claims struct {
	Subject       string // Google's stable per-user id ("sub")
	Email         string
	EmailVerified bool
	Name          string
}

// Verifier checks a Google ID token and returns its claims if valid.
type Verifier interface {
	Verify(ctx context.Context, idToken string) (*Claims, error)
}

// NewVerifier constructs a Verifier that checks tokens were issued for
// clientID (your OAuth Client ID — not secret, safe to read from env).
func NewVerifier(clientID string) Verifier {
	return newVerifier(clientID, defaultCertsURL)
}

func newVerifier(clientID, certsURL string) *verifier {
	return &verifier{
		clientID:   clientID,
		certsURL:   certsURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		keys:       map[string]*rsa.PublicKey{},
	}
}

type verifier struct {
	clientID   string
	certsURL   string
	httpClient *http.Client

	mu        sync.Mutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
}

type tokenClaims struct {
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Exp           int64  `json:"exp"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

func (v *verifier) Verify(ctx context.Context, idToken string) (*Claims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	if err := decodeSegment(parts[0], &header); err != nil {
		return nil, ErrInvalidToken
	}
	if header.Alg != "RS256" {
		return nil, ErrInvalidToken
	}

	key, err := v.keyFor(ctx, header.Kid)
	if err != nil {
		return nil, err
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}
	hashed := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature); err != nil {
		return nil, ErrInvalidToken
	}

	var claims tokenClaims
	if err := decodeSegment(parts[1], &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Iss != issuerHTTPS && claims.Iss != issuerBare {
		return nil, ErrInvalidToken
	}
	if claims.Aud != v.clientID {
		return nil, ErrInvalidToken
	}
	if time.Now().Unix() >= claims.Exp {
		return nil, ErrInvalidToken
	}

	return &Claims{
		Subject:       claims.Sub,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Name:          claims.Name,
	}, nil
}

// keyFor returns the public key for kid, fetching/refreshing the key set
// if needed. A kid that's already cached and fresh costs no network call.
func (v *verifier) keyFor(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.Lock()
	key, ok := v.keys[kid]
	fresh := ok && time.Since(v.fetchedAt) < keysCacheTTL
	v.mu.Unlock()
	if fresh {
		return key, nil
	}

	if err := v.refreshKeys(ctx); err != nil {
		if ok {
			// Refresh failed (network blip) but we already know this kid
			// from a previous fetch — serve it rather than fail outright.
			return key, nil
		}
		return nil, err
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	key, ok = v.keys[kid]
	if !ok {
		return nil, ErrInvalidToken
	}
	return key, nil
}

type jwkSet struct {
	Keys []struct {
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

func (v *verifier) refreshKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.certsURL, nil)
	if err != nil {
		return err
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetching google certs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching google certs: status %d", resp.StatusCode)
	}

	var set jwkSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return fmt.Errorf("decoding google certs: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(set.Keys))
	for _, k := range set.Keys {
		if k.Kty != "RSA" {
			continue
		}
		pub, err := rsaPublicKeyFromJWK(k.N, k.E)
		if err != nil {
			continue
		}
		keys[k.Kid] = pub
	}

	v.mu.Lock()
	v.keys = keys
	v.fetchedAt = time.Now()
	v.mu.Unlock()
	return nil
}

func rsaPublicKeyFromJWK(nParam, eParam string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nParam)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eParam)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(new(big.Int).SetBytes(eBytes).Int64()),
	}, nil
}

func decodeSegment(segment string, v any) error {
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(decoded, v)
}
