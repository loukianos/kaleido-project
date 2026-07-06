package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/identity"
	"kaleido-project/internal/loans"
)

func TestRequestsWithoutTokenAreUnauthorized(t *testing.T) {
	handler := newTestHandler(Options{})

	for _, route := range []struct{ method, path string }{
		{http.MethodPost, "/loans"},
		{http.MethodGet, "/loans"},
		{http.MethodGet, "/loans/1"},
		{http.MethodPost, "/loans/1/transfer"},
		{http.MethodPost, "/admin/contracts/deploy"},
		{http.MethodGet, "/contracts"},
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(route.method, route.path, nil))
		require.Equal(t, http.StatusUnauthorized, recorder.Code, "%s %s", route.method, route.path)
	}
}

func TestRequestsWithBadTokenAreUnauthorized(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans", nil)
	request.Header.Set("Authorization", "Bearer forged-token")
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestSystemEndpointsArePublic(t *testing.T) {
	handler := newTestHandler(Options{})

	for _, path := range []string{"/", "/healthz"} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		require.Equal(t, http.StatusOK, recorder.Code, path)
	}
}

func TestLenderCannotOriginate(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_subject":"alice",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asAlice(request))

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestLenderCannotDeployContracts(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin/contracts/deploy", strings.NewReader(`{}`))
	handler.ServeHTTP(recorder, asAlice(request))

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAdminWithoutServicerRoleCannotOriginate(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_subject":"alice",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asAdmin(request))

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestLenderListIsScopedToTheirLoans(t *testing.T) {
	service := &fakeLoansService{}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asAlice(httptest.NewRequest(http.MethodGet, "/loans", nil)))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, aliceIdentityID, service.listRequest.LenderIdentityID)
}

func TestServicerListIsUnscoped(t *testing.T) {
	service := &fakeLoansService{}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asServicer(httptest.NewRequest(http.MethodGet, "/loans", nil)))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, int64(0), service.listRequest.LenderIdentityID)
}

func TestLenderCannotReadForeignLoan(t *testing.T) {
	service := &fakeLoansService{
		read: loans.ReadResult{Loan: db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", LenderIdentityID: db.Ptr(bobIdentityID), Status: "active"}},
		loan: db.Loan{ID: 7, LenderIdentityID: db.Ptr(bobIdentityID)},
	}
	handler := newTestHandler(Options{Loans: service})

	for _, path := range []string{"/loans/7", "/loans/7/repayments", "/loans/7/terms"} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, asAlice(httptest.NewRequest(http.MethodGet, path, nil)))
		require.Equal(t, http.StatusNotFound, recorder.Code, path)
	}

	// The loan's holder still sees it.
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asBob(httptest.NewRequest(http.MethodGet, "/loans/7", nil)))
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestLenderReadsTheirOwnLoan(t *testing.T) {
	service := &fakeLoansService{
		read: loans.ReadResult{Loan: db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", LenderIdentityID: db.Ptr(aliceIdentityID), Status: "active"}, LenderSubject: "alice"},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asAlice(httptest.NewRequest(http.MethodGet, "/loans/7", nil)))

	require.Equal(t, http.StatusOK, recorder.Code)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "alice", body.LenderSubject)
}

func TestTransferPassesLenderCaller(t *testing.T) {
	service := &fakeLoansService{
		transferResult: loans.TransferResult{Loan: db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", Status: "active"}},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/transfer", strings.NewReader(`{"to_subject":"bob"}`))
	handler.ServeHTTP(recorder, asAlice(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.False(t, service.transferCaller.Servicer)
	require.Equal(t, aliceIdentityID, service.transferCaller.IdentityID)
}

func TestOnboardLender(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asAlice(httptest.NewRequest(http.MethodPost, "/lenders/onboard", nil)))

	require.Equal(t, http.StatusOK, recorder.Code)

	var body onboardResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "alice", body.Subject)
	require.Equal(t, testIssuer, body.Issuer)
	require.NotEmpty(t, body.Address)
}

func TestOnboardRejectsPlatformAccounts(t *testing.T) {
	handler := newTestHandler(Options{})

	for name, wrap := range map[string]func(*http.Request) *http.Request{
		"servicer": asServicer,
		"admin":    asAdmin,
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, wrap(httptest.NewRequest(http.MethodPost, "/lenders/onboard", nil)))
		require.Equal(t, http.StatusForbidden, recorder.Code, name)
	}
}

func TestOnboardRequiresToken(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/lenders/onboard", nil))

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestCreateLoanLenderNotOnboarded(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{err: fmt.Errorf("resolve lender identity: %w", identity.ErrNotOnboarded)},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_subject":"never-onboarded",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestTransferPassesServicerCaller(t *testing.T) {
	service := &fakeLoansService{
		transferResult: loans.TransferResult{Loan: db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", Status: "active"}},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/transfer", strings.NewReader(`{"to_address":"0x1111111111111111111111111111111111111111"}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.True(t, service.transferCaller.Servicer)
}
