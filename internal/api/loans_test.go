package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"kaleido-project/db/sqlc"
	"kaleido-project/internal/contracts"
	"kaleido-project/internal/loans"
)

func TestCreateLoan(t *testing.T) {
	service := &fakeLoansService{
		result: loans.OriginateResult{
			Loan: db.Loan{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				PrincipalMinor:   100_00,
				AprBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
				OutstandingMinor: 108_00,
				Status:           "active",
			},
			OperationID: 11,
			TxHash:      "0xabc",
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, int64(100_00), service.request.PrincipalMinor)
	require.Equal(t, uint16(800), service.request.APRBps)
	require.Nil(t, service.request.ContractID)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(7), body.ID)
	require.Equal(t, "0", body.TokenID)
	require.Equal(t, int64(11), body.OperationID)
	require.Equal(t, "0xabc", body.TxHash)
}

func TestCreateLoanWithContractID(t *testing.T) {
	service := &fakeLoansService{
		result: loans.OriginateResult{
			Loan: db.Loan{
				ID:            8,
				TokenID:       db.Ptr("0"),
				ContractID:    db.Ptr(int64(2)),
				BorrowerRef:   "borrower-1",
				LenderAddress: "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				Status:        "active",
			},
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365,
		"contract_id":2
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.NotNil(t, service.request.ContractID)
	require.Equal(t, int64(2), *service.request.ContractID)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(2), body.ContractID)
}

func TestCreateLoanContractNotFound(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{err: contracts.ErrContractNotFound},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365,
		"contract_id":99
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateLoanValidation(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{"borrower_ref":""}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateLoanWithLenderSubject(t *testing.T) {
	service := &fakeLoansService{
		result: loans.OriginateResult{
			Loan: db.Loan{
				ID:               9,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				LenderIdentityID: db.Ptr(int64(1)),
				Status:           "active",
			},
			LenderSubject: "alice",
		},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_subject":"alice",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, "alice", service.request.LenderSubject)
	require.Empty(t, service.request.LenderAddress)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "alice", body.LenderSubject)
}

func TestCreateLoanLenderExclusivity(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{err: loans.ErrInvalidLender},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"lender_subject":"alice",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestTransferLoanToSubject(t *testing.T) {
	service := &fakeLoansService{
		transferResult: loans.TransferResult{
			Loan: db.Loan{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0x9999999999999999999999999999999999999999",
				LenderIdentityID: db.Ptr(int64(2)),
				Status:           "active",
			},
			LenderSubject: "bob",
			OperationID:   13,
			TxHash:        "0xtransfer",
		},
	}
	handler := newTestHandler(Options{Loans: service})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/transfer", strings.NewReader(`{"to_subject":"bob"}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "bob", service.transferRequest.ToSubject)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "bob", body.LenderSubject)
}

func TestCreateLoanDomainError(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{err: loans.ErrInvalidAmount},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":0,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateLoanNoActiveContract(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{err: loans.ErrNoActiveContract},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusConflict, recorder.Code)

	var body struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, loans.ErrNoActiveContract.Error(), body.Error)
}

func TestCreateLoanLockBusy(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{err: fmt.Errorf("originate: %w", db.ErrLockBusy)},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{
		"borrower_ref":"borrower-1",
		"lender_address":"0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
		"principal_minor":10000,
		"apr_bps":800,
		"term_days":365
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestGetLoan(t *testing.T) {
	service := &fakeLoansService{
		read: loans.ReadResult{
			Loan: db.Loan{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				PrincipalMinor:   100_00,
				AprBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
				OutstandingMinor: 108_00,
				Status:           "active",
				MintOperationID:  db.Ptr(int64(11)),
			},
			OwnerAddress: "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
			MintTxHash:   "0xabc",
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans/7", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, int64(7), service.getID)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(7), body.ID)
	require.NotEmpty(t, body.OwnerAddress)
	require.Equal(t, "0xabc", body.TxHash)
	require.Equal(t, int64(11), body.OperationID)
}

func TestGetLoanNotFound(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{getErr: pgx.ErrNoRows},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans/99", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestListLoans(t *testing.T) {
	service := &fakeLoansService{
		list: []db.Loan{
			{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				PrincipalMinor:   100_00,
				AprBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
				OutstandingMinor: 108_00,
				Status:           "active",
			},
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans?lender=0xabc&status=active&limit=10&offset=5", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "0xabc", service.listRequest.Lender)
	require.Equal(t, "active", service.listRequest.Status)
	require.Equal(t, int32(10), service.listRequest.Limit)
	require.Equal(t, int32(5), service.listRequest.Offset)

	var body struct {
		Loans []loanResponse `json:"loans"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Len(t, body.Loans, 1)
	require.Equal(t, int64(7), body.Loans[0].ID)
}

func TestListLoansAbsentLimitPassesZero(t *testing.T) {
	service := &fakeLoansService{}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, int32(0), service.listRequest.Limit)
}

func TestListLoansInvalidPagination(t *testing.T) {
	tests := map[string]string{
		"limit zero":         "/loans?limit=0",
		"limit negative":     "/loans?limit=-1",
		"limit over max":     "/loans?limit=101",
		"limit not numeric":  "/loans?limit=abc",
		"offset negative":    "/loans?offset=-1",
		"offset not numeric": "/loans?offset=abc",
	}

	for name, target := range tests {
		t.Run(name, func(t *testing.T) {
			handler := newTestHandler(Options{
				Loans: &fakeLoansService{},
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, target, nil)
			handler.ServeHTTP(recorder, asServicer(request))

			require.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	}
}

func TestLoanTerms(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{
			terms: loans.LoanTerms{
				PrincipalMinor:   100_00,
				APRBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
			},
		},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans/7/terms", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)

	var body termsResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(108_00), body.TotalDueMinor)
}

func TestTransferLoan(t *testing.T) {
	service := &fakeLoansService{
		transferResult: loans.TransferResult{
			Loan: db.Loan{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0x1111111111111111111111111111111111111111",
				PrincipalMinor:   100_00,
				AprBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
				OutstandingMinor: 108_00,
				Status:           "active",
			},
			OperationID: 12,
			TxHash:      "0xtransfer",
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/transfer", strings.NewReader(`{
		"to_address":"0x1111111111111111111111111111111111111111"
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, int64(7), service.transferLoanID)
	require.Equal(t, "0x1111111111111111111111111111111111111111", service.transferRequest.ToAddress)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "0x1111111111111111111111111111111111111111", body.LenderAddress)
	require.Equal(t, int64(12), body.OperationID)
	require.Equal(t, "0xtransfer", body.TxHash)
}

func TestTransferLoanConflict(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{transferErr: loans.ErrLoanNotTransferable},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/transfer", strings.NewReader(`{
		"to_address":"0x1111111111111111111111111111111111111111"
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestTransferLoanNotNoteOwner(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{
			transferErr: fmt.Errorf("%w: note is held by 0x2222222222222222222222222222222222222222", loans.ErrNotNoteOwner),
		},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/transfer", strings.NewReader(`{
		"to_address":"0x1111111111111111111111111111111111111111"
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusForbidden, recorder.Code)

	var body errorResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Contains(t, body.Error, "note is held by")
}

func TestDefaultLoan(t *testing.T) {
	service := &fakeLoansService{
		defaultResult: loans.DefaultResult{
			Loan: db.Loan{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				PrincipalMinor:   100_00,
				AprBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
				OutstandingMinor: 108_00,
				Status:           "defaulted",
			},
			OperationID: 13,
			TxHash:      "0xdefault",
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/default", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, int64(7), service.defaultLoanID)

	var body loanResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "defaulted", body.Status)
	require.Equal(t, int64(13), body.OperationID)
	require.Equal(t, "0xdefault", body.TxHash)
}

func TestDefaultLoanConflict(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{defaultErr: loans.ErrLoanNotActive},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/default", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateRepayment(t *testing.T) {
	service := &fakeLoansService{
		repaymentResult: loans.RepaymentResult{
			Loan: db.Loan{
				ID:               7,
				TokenID:          db.Ptr("0"),
				BorrowerRef:      "borrower-1",
				LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
				PrincipalMinor:   100_00,
				AprBps:           800,
				TermDays:         365,
				InterestDueMinor: 800,
				TotalDueMinor:    108_00,
				OutstandingMinor: 58_00,
				Status:           "active",
			},
			Repayment: db.Repayment{
				ID:          3,
				LoanID:      7,
				AmountMinor: 50_00,
				ExternalRef: db.Ptr("payment-1"),
				CreatedAt:   time.Unix(1_700_000_000, 0).UTC(),
			},
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/repayments", strings.NewReader(`{
		"amount_minor":5000,
		"external_ref":"payment-1"
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, int64(7), service.repaymentLoanID)
	require.Equal(t, int64(50_00), service.repaymentRequest.AmountMinor)

	var body repaymentResultResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, int64(3), body.Repayment.ID)
	require.Equal(t, "payment-1", body.Repayment.ExternalRef)
	require.Equal(t, int64(58_00), body.Loan.OutstandingMinor)
}

func TestCreateFinalRepayment(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{
			repaymentResult: loans.RepaymentResult{
				Loan: db.Loan{
					ID:               7,
					TokenID:          db.Ptr("0"),
					BorrowerRef:      "borrower-1",
					LenderAddress:    "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
					PrincipalMinor:   100_00,
					AprBps:           800,
					TermDays:         365,
					InterestDueMinor: 800,
					TotalDueMinor:    108_00,
					OutstandingMinor: 0,
					Status:           "repaid",
				},
				Repayment:             db.Repayment{ID: 4, LoanID: 7, AmountMinor: 108_00},
				SettlementOperationID: 21,
				SettlementTxHash:      "0xsettle",
			},
		},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/repayments", strings.NewReader(`{"amount_minor":10800}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusCreated, recorder.Code)

	var body repaymentResultResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "repaid", body.Loan.Status)
	require.Equal(t, int64(21), body.SettlementOperationID)
	require.Equal(t, "0xsettle", body.SettlementTxHash)
}

func TestCreateRepaymentConflict(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{repaymentErr: loans.ErrLoanNotActive},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/repayments", strings.NewReader(`{"amount_minor":5000}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateRepaymentDuplicateExternalRef(t *testing.T) {
	handler := newTestHandler(Options{
		Loans: &fakeLoansService{repaymentErr: loans.ErrDuplicateExternalRef},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/loans/7/repayments", strings.NewReader(`{
		"amount_minor":5000,
		"external_ref":"payment-1"
	}`))
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestListRepayments(t *testing.T) {
	service := &fakeLoansService{
		repayments: []db.Repayment{
			{ID: 3, LoanID: 7, AmountMinor: 50_00, CreatedAt: time.Unix(1_700_000_000, 0).UTC()},
		},
	}
	handler := newTestHandler(Options{
		Loans: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/loans/7/repayments", nil)
	handler.ServeHTTP(recorder, asServicer(request))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, int64(7), service.listRepaymentsLoanID)

	var body struct {
		Repayments []repaymentResponse `json:"repayments"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Len(t, body.Repayments, 1)
	require.Equal(t, int64(50_00), body.Repayments[0].AmountMinor)
}

type fakeLoansService struct {
	request              loans.OriginateRequest
	result               loans.OriginateResult
	err                  error
	getID                int64
	read                 loans.ReadResult
	getErr               error
	listRequest          loans.ListRequest
	list                 []db.Loan
	listErr              error
	termsID              int64
	terms                loans.LoanTerms
	termsErr             error
	repaymentLoanID      int64
	repaymentRequest     loans.RepaymentRequest
	repaymentResult      loans.RepaymentResult
	repaymentErr         error
	listRepaymentsLoanID int64
	repayments           []db.Repayment
	listRepaymentsErr    error
	transferLoanID       int64
	transferRequest      loans.TransferRequest
	transferCaller       loans.Caller
	transferResult       loans.TransferResult
	transferErr          error
	loan                 db.Loan
	loanErr              error
	operation            db.ChainOperation
	operationErr         error
	defaultLoanID        int64
	defaultResult        loans.DefaultResult
	defaultErr           error
}

func (f *fakeLoansService) Originate(_ context.Context, req loans.OriginateRequest) (loans.OriginateResult, error) {
	f.request = req
	// Pending semantics return the partial result alongside the error, so the fake does too.
	return f.result, f.err
}

func (f *fakeLoansService) Get(_ context.Context, id int64) (loans.ReadResult, error) {
	f.getID = id
	if f.getErr != nil {
		return loans.ReadResult{}, f.getErr
	}
	return f.read, nil
}

func (f *fakeLoansService) List(_ context.Context, req loans.ListRequest) ([]db.Loan, error) {
	f.listRequest = req
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.list, nil
}

func (f *fakeLoansService) Terms(_ context.Context, id int64) (loans.LoanTerms, error) {
	f.termsID = id
	if f.termsErr != nil {
		return loans.LoanTerms{}, f.termsErr
	}
	return f.terms, nil
}

func (f *fakeLoansService) RecordRepayment(_ context.Context, loanID int64, req loans.RepaymentRequest) (loans.RepaymentResult, error) {
	f.repaymentLoanID = loanID
	f.repaymentRequest = req
	return f.repaymentResult, f.repaymentErr
}

func (f *fakeLoansService) ListRepayments(_ context.Context, loanID int64) ([]db.Repayment, error) {
	f.listRepaymentsLoanID = loanID
	if f.listRepaymentsErr != nil {
		return nil, f.listRepaymentsErr
	}
	return f.repayments, nil
}

func (f *fakeLoansService) Transfer(_ context.Context, loanID int64, req loans.TransferRequest, caller loans.Caller) (loans.TransferResult, error) {
	f.transferLoanID = loanID
	f.transferRequest = req
	f.transferCaller = caller
	if f.transferErr != nil {
		return loans.TransferResult{}, f.transferErr
	}
	return f.transferResult, nil
}

func (f *fakeLoansService) Loan(_ context.Context, _ int64) (db.Loan, error) {
	if f.loanErr != nil {
		return db.Loan{}, f.loanErr
	}
	return f.loan, nil
}

func (f *fakeLoansService) Operation(_ context.Context, _ int64) (db.ChainOperation, error) {
	if f.operationErr != nil {
		return db.ChainOperation{}, f.operationErr
	}
	return f.operation, nil
}

func (f *fakeLoansService) Default(_ context.Context, loanID int64) (loans.DefaultResult, error) {
	f.defaultLoanID = loanID
	return f.defaultResult, f.defaultErr
}
