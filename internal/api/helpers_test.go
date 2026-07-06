package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/auth"
)

// Test bearer tokens the default fake verifier accepts.
// Handler tests set these on requests instead of minting real JWTs; the auth package owns real-token coverage.
const (
	servicerToken = "servicer-token"
	adminToken    = "admin-token"
	aliceToken    = "alice-token"
	bobToken      = "bob-token"
	testIssuer    = "http://issuer.test/realms/loan-notes"
)

// Identity ids the default fake identity service assigns.
const (
	aliceIdentityID int64 = 1
	bobIdentityID   int64 = 2
)

func newTestHandler(opts Options) http.Handler {
	if opts.Loans == nil {
		opts.Loans = &fakeLoansService{}
	}
	if opts.Contracts == nil {
		opts.Contracts = &fakeContractsService{}
	}
	if opts.Verifier == nil {
		opts.Verifier = fakeVerifier{
			servicerToken: {Issuer: testIssuer, Subject: "servicer-sa", Roles: []string{auth.RoleServicer, auth.RoleAdmin}},
			adminToken:    {Issuer: testIssuer, Subject: "admin-sa", Roles: []string{auth.RoleAdmin}},
			aliceToken:    {Issuer: testIssuer, Subject: "alice", Roles: nil},
			bobToken:      {Issuer: testIssuer, Subject: "bob", Roles: nil},
		}
	}
	if opts.Identities == nil {
		opts.Identities = fakeIdentityService{
			"alice": aliceIdentityID,
			"bob":   bobIdentityID,
		}
	}
	return New("test", slog.New(slog.NewTextHandler(io.Discard, nil)), opts)
}

// asServicer et al. attach a bearer token the default fake verifier maps to that principal.
func asServicer(r *http.Request) *http.Request { return withToken(r, servicerToken) }
func asAdmin(r *http.Request) *http.Request    { return withToken(r, adminToken) }
func asAlice(r *http.Request) *http.Request    { return withToken(r, aliceToken) }
func asBob(r *http.Request) *http.Request      { return withToken(r, bobToken) }

func withToken(r *http.Request, token string) *http.Request {
	r.Header.Set("Authorization", "Bearer "+token)
	return r
}

type fakeVerifier map[string]auth.Principal

func (f fakeVerifier) Verify(_ context.Context, raw string) (auth.Principal, error) {
	principal, ok := f[raw]
	if !ok {
		return auth.Principal{}, auth.ErrInvalidToken
	}
	return principal, nil
}

type fakeIdentityService map[string]int64

func (f fakeIdentityService) ResolveIdentity(_ context.Context, issuer, subject string) (db.Identity, error) {
	id, ok := f[subject]
	if !ok {
		return db.Identity{}, errors.New("unknown test subject: " + subject)
	}
	return db.Identity{ID: id, Issuer: issuer, Subject: subject, Role: "lender"}, nil
}
