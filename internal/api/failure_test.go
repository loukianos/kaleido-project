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
	"kaleido-project/internal/loans"
)

func TestCreateLoanPendingReturns202(t *testing.T) {
	service := &fakeLoansService{
		result: loans.OriginateResult{
			Loan:        db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", Status: "originating"},
			OperationID: 11,
		},
		err: fmt.Errorf("originate: %w", loans.ErrOperationPending),
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"b",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":0,
		"term_days":30
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusAccepted, recorder.Code)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(7), body.ID)
	require.Equal(t, "originating", body.Status)
	require.Equal(t, int64(11), body.OperationID)
	require.Contains(t, body.Message, "pending")
}

func TestCreateLoanIdempotentReplayReturns200(t *testing.T) {
	service := &fakeLoansService{
		result: loans.OriginateResult{
			Loan:     db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", Status: "active", ExternalRef: db.Ptr("order-1")},
			Existing: true,
		},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"b",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":0,
		"term_days":30,
		"external_ref":"order-1"
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "order-1", service.request.ExternalRef)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "order-1", body.ExternalRef)
}

func TestRepaymentPendingSettlementReturns202(t *testing.T) {
	service := &fakeLoansService{
		repaymentResult: loans.RepaymentResult{
			Loan:                  db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", Status: "settling"},
			Repayment:             db.Repayment{ID: 3, LoanID: 7, AmountMinor: 100_00},
			SettlementOperationID: 12,
		},
		repaymentErr: fmt.Errorf("settle: %w", loans.ErrOperationPending),
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/repayments", strings.NewReader(`{"amount_minor":10000}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusAccepted, recorder.Code)

	var body repaymentResultResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(3), body.Repayment.ID)
	require.Equal(t, "settling", body.Loan.Status)
	require.Equal(t, int64(12), body.SettlementOperationID)
	require.Contains(t, body.Message, "pending")
}

func TestDefaultPendingReturns202(t *testing.T) {
	service := &fakeLoansService{
		defaultResult: loans.DefaultResult{
			Loan:        db.Loan{ID: 7, BorrowerRef: "b", LenderAddress: "0x1", Status: "active"},
			OperationID: 13,
		},
		defaultErr: fmt.Errorf("default: %w", loans.ErrOperationPending),
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/default", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusAccepted, recorder.Code)
}

func TestGetOperation(t *testing.T) {
	service := &fakeLoansService{
		operation: db.ChainOperation{
			ID:            11,
			Kind:          "originate",
			Status:        "retryable",
			Attempts:      2,
			Error:         db.Ptr("dial tcp: connection refused"),
			LoanID:        db.Ptr(int64(7)),
			SignerAddress: db.Ptr("0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73"),
		},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asServicer(httptest.NewRequest(http.MethodGet, "/operations/11", nil)))

	require.Equal(t, http.StatusOK, recorder.Code)

	var body operationResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "retryable", body.Status)
	require.Equal(t, int32(2), body.Attempts)
	require.Contains(t, body.Error, "connection refused")
	require.Equal(t, int64(7), body.LoanID)
}

func TestGetOperationIsServicerOnly(t *testing.T) {
	handler := newTestHandler(Options{})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, asAlice(httptest.NewRequest(http.MethodGet, "/operations/11", nil)))

	require.Equal(t, http.StatusForbidden, recorder.Code)
}
