package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/require"
)

const testAudience = "loan-notes-api"

type tokenClaims struct {
	Issuer      string   `json:"iss"`
	Subject     string   `json:"sub"`
	Audience    []string `json:"aud"`
	Expiry      int64    `json:"exp"`
	IssuedAt    int64    `json:"iat"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

// testIssuer serves a JWKS over httptest and mints signed tokens for it.
type testIssuer struct {
	key    *rsa.PrivateKey
	server *httptest.Server
}

func newTestIssuer(t *testing.T) *testIssuer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key:       key.Public(),
		KeyID:     "test-key",
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	t.Cleanup(server.Close)
	return &testIssuer{key: key, server: server}
}

func (ti *testIssuer) mint(t *testing.T, claims tokenClaims) string {
	t.Helper()
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: ti.key},
		(&jose.SignerOptions{}).WithHeader("kid", "test-key"),
	)
	require.NoError(t, err)

	payload, err := json.Marshal(claims)
	require.NoError(t, err)
	jws, err := signer.Sign(payload)
	require.NoError(t, err)
	raw, err := jws.CompactSerialize()
	require.NoError(t, err)
	return raw
}

func (ti *testIssuer) verifier(t *testing.T, issuerURL string) *OIDCVerifier {
	t.Helper()
	verifier, err := NewOIDCVerifier(context.Background(), OIDCConfig{
		IssuerURL: issuerURL,
		JWKSURL:   ti.server.URL,
		Audience:  testAudience,
	})
	require.NoError(t, err)
	return verifier
}

func validClaims(issuer string) tokenClaims {
	claims := tokenClaims{
		Issuer:   issuer,
		Subject:  "alice",
		Audience: []string{testAudience},
		Expiry:   time.Now().Add(time.Hour).Unix(),
		IssuedAt: time.Now().Unix(),
	}
	return claims
}

func TestVerifyReturnsPrincipal(t *testing.T) {
	issuer := newTestIssuer(t)
	verifier := issuer.verifier(t, "http://issuer.test/realms/loan-notes")

	claims := validClaims("http://issuer.test/realms/loan-notes")
	claims.RealmAccess.Roles = []string{RoleServicer, RoleAdmin}

	principal, err := verifier.Verify(context.Background(), issuer.mint(t, claims))
	require.NoError(t, err)
	require.Equal(t, "alice", principal.Subject)
	require.Equal(t, "http://issuer.test/realms/loan-notes", principal.Issuer)
	require.True(t, principal.HasRole(RoleServicer))
	require.True(t, principal.HasRole(RoleAdmin))
}

func TestVerifyLenderHasNoRoles(t *testing.T) {
	issuer := newTestIssuer(t)
	verifier := issuer.verifier(t, "http://issuer.test/realms/loan-notes")

	principal, err := verifier.Verify(context.Background(), issuer.mint(t, validClaims("http://issuer.test/realms/loan-notes")))
	require.NoError(t, err)
	require.False(t, principal.HasRole(RoleServicer))
	require.False(t, principal.HasRole(RoleAdmin))
}

func TestVerifyRejectsWrongIssuer(t *testing.T) {
	issuer := newTestIssuer(t)
	verifier := issuer.verifier(t, "http://issuer.test/realms/loan-notes")

	_, err := verifier.Verify(context.Background(), issuer.mint(t, validClaims("http://evil.test/realms/loan-notes")))
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyRejectsWrongAudience(t *testing.T) {
	issuer := newTestIssuer(t)
	verifier := issuer.verifier(t, "http://issuer.test/realms/loan-notes")

	claims := validClaims("http://issuer.test/realms/loan-notes")
	claims.Audience = []string{"someone-else"}

	_, err := verifier.Verify(context.Background(), issuer.mint(t, claims))
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	issuer := newTestIssuer(t)
	verifier := issuer.verifier(t, "http://issuer.test/realms/loan-notes")

	claims := validClaims("http://issuer.test/realms/loan-notes")
	claims.Expiry = time.Now().Add(-time.Hour).Unix()

	_, err := verifier.Verify(context.Background(), issuer.mint(t, claims))
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerifyRejectsGarbage(t *testing.T) {
	issuer := newTestIssuer(t)
	verifier := issuer.verifier(t, "http://issuer.test/realms/loan-notes")

	_, err := verifier.Verify(context.Background(), "not-a-jwt")
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestNewOIDCVerifierRequiresConfig(t *testing.T) {
	_, err := NewOIDCVerifier(context.Background(), OIDCConfig{})
	require.Error(t, err)
}
