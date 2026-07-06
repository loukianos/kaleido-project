// Package auth verifies bearer tokens and describes who is calling.
// Identity comes from OIDC: the token's (issuer, subject) names the caller and realm roles carry servicer/admin powers.
package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/coreos/go-oidc/v3/oidc"
)

const (
	RoleServicer = "servicer"
	RoleAdmin    = "admin"
)

var ErrInvalidToken = errors.New("invalid bearer token")

// Principal is the authenticated caller.
// A principal with no roles is a lender; servicer and admin powers come only from token claims.
type Principal struct {
	Issuer  string
	Subject string
	Roles   []string
}

func (p Principal) HasRole(role string) bool {
	return slices.Contains(p.Roles, role)
}

// Verifier turns a raw bearer token into a Principal.
type Verifier interface {
	Verify(ctx context.Context, rawToken string) (Principal, error)
}

// OIDCConfig points at the issuer.
// JWKSURL is separate from IssuerURL because the network path to the issuer can differ from the issuer name in tokens (e.g. a compose service reaching Keycloak by service name while tokens carry the host-published issuer).
type OIDCConfig struct {
	IssuerURL string
	JWKSURL   string
	Audience  string
}

// OIDCVerifier validates JWTs against the issuer's JWKS, checking issuer and audience.
type OIDCVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewOIDCVerifier(ctx context.Context, cfg OIDCConfig) (*OIDCVerifier, error) {
	if cfg.IssuerURL == "" || cfg.JWKSURL == "" || cfg.Audience == "" {
		return nil, errors.New("oidc issuer url, jwks url, and audience are required")
	}
	keySet := oidc.NewRemoteKeySet(ctx, cfg.JWKSURL)
	return &OIDCVerifier{
		verifier: oidc.NewVerifier(cfg.IssuerURL, keySet, &oidc.Config{ClientID: cfg.Audience}),
	}, nil
}

func (v *OIDCVerifier) Verify(ctx context.Context, rawToken string) (Principal, error) {
	token, err := v.verifier.Verify(ctx, rawToken)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	// Keycloak carries realm roles in the realm_access claim.
	var claims struct {
		RealmAccess struct {
			Roles []string `json:"roles"`
		} `json:"realm_access"`
	}
	if err := token.Claims(&claims); err != nil {
		return Principal{}, fmt.Errorf("%w: parse claims: %s", ErrInvalidToken, err)
	}

	return Principal{
		Issuer:  token.Issuer,
		Subject: token.Subject,
		Roles:   claims.RealmAccess.Roles,
	}, nil
}
