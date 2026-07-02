package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"kaleido-project/db/sqlc"
	"kaleido-project/internal/loans"
)

type LoansService interface {
	Originate(context.Context, loans.OriginateRequest) (loans.OriginateResult, error)
	Get(context.Context, int64) (loans.ReadResult, error)
	List(context.Context, loans.ListRequest) ([]db.Loan, error)
	Terms(context.Context, int64) (loans.LoanTerms, error)
	RecordRepayment(context.Context, int64, loans.RepaymentRequest) (loans.RepaymentResult, error)
	ListRepayments(context.Context, int64) ([]db.Repayment, error)
	Transfer(context.Context, int64, loans.TransferRequest) (loans.TransferResult, error)
	Default(context.Context, int64) (loans.DefaultResult, error)
}

type createLoanRequest struct {
	BorrowerRef    string `json:"borrower_ref"`
	LenderAddress  string `json:"lender_address"`
	PrincipalMinor int64  `json:"principal_minor"`
	APRBps         uint16 `json:"apr_bps"`
	TermDays       int64  `json:"term_days"`
}

func handleCreateLoan(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decode[createLoanRequest](r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody("invalid json body"))
			return
		}
		if req.BorrowerRef == "" || req.LenderAddress == "" {
			writeJSON(w, http.StatusBadRequest, errorBody("borrower_ref and lender_address are required"))
			return
		}

		result, err := service.Originate(r.Context(), loans.OriginateRequest{
			BorrowerRef:    req.BorrowerRef,
			LenderAddress:  req.LenderAddress,
			PrincipalMinor: req.PrincipalMinor,
			APRBps:         req.APRBps,
			TermDays:       req.TermDays,
		})
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "create loan failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("create loan failed"))
			return
		}

		writeJSON(w, http.StatusCreated, loanResponseWithTx(result.Loan, result.OperationID, result.TxHash))
	})
}

func handleGetLoan(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}

		result, err := service.Get(r.Context(), id)
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "get loan failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("get loan failed"))
			return
		}
		writeJSON(w, http.StatusOK, loanResponseFromRead(result))
	})
}

func handleListLoans(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Zero means "not specified"; the service applies its default.
		limit := int32(0)
		if raw := r.URL.Query().Get("limit"); raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 32)
			if err != nil || parsed < 1 || parsed > 100 {
				writeJSON(w, http.StatusBadRequest, errorBody("limit must be between 1 and 100"))
				return
			}
			limit = int32(parsed)
		}
		offset := int32(0)
		if raw := r.URL.Query().Get("offset"); raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 32)
			if err != nil || parsed < 0 {
				writeJSON(w, http.StatusBadRequest, errorBody("offset must be non-negative"))
				return
			}
			offset = int32(parsed)
		}

		items, err := service.List(r.Context(), loans.ListRequest{
			Lender: r.URL.Query().Get("lender"),
			Status: r.URL.Query().Get("status"),
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			logger.ErrorContext(r.Context(), "list loans failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("list loans failed"))
			return
		}

		responses := make([]loanResponse, 0, len(items))
		for _, item := range items {
			responses = append(responses, loanResponseFromLoan(item))
		}
		writeJSON(w, http.StatusOK, map[string]any{"loans": responses})
	})
}

func handleLoanTerms(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}

		terms, err := service.Terms(r.Context(), id)
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "get loan terms failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("get loan terms failed"))
			return
		}
		writeJSON(w, http.StatusOK, termsResponse{
			PrincipalMinor:   terms.PrincipalMinor,
			APRBps:           terms.APRBps,
			TermDays:         terms.TermDays,
			InterestDueMinor: terms.InterestDueMinor,
			TotalDueMinor:    terms.TotalDueMinor,
		})
	})
}

type transferLoanRequest struct {
	ToAddress string `json:"to_address"`
}

func handleTransferLoan(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}

		req, err := decode[transferLoanRequest](r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		result, err := service.Transfer(r.Context(), id, loans.TransferRequest{ToAddress: req.ToAddress})
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "transfer loan failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("transfer loan failed"))
			return
		}

		writeJSON(w, http.StatusOK, loanResponseWithTx(result.Loan, result.OperationID, result.TxHash))
	})
}

func handleDefaultLoan(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}

		result, err := service.Default(r.Context(), id)
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "default loan failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("default loan failed"))
			return
		}

		writeJSON(w, http.StatusOK, loanResponseWithTx(result.Loan, result.OperationID, result.TxHash))
	})
}

type createRepaymentRequest struct {
	AmountMinor int64  `json:"amount_minor"`
	ExternalRef string `json:"external_ref"`
}

func handleCreateRepayment(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}

		req, err := decode[createRepaymentRequest](r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		result, err := service.RecordRepayment(r.Context(), id, loans.RepaymentRequest{
			AmountMinor: req.AmountMinor,
			ExternalRef: req.ExternalRef,
		})
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "create repayment failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("create repayment failed"))
			return
		}

		writeJSON(w, http.StatusCreated, repaymentResultResponse{
			Repayment:             repaymentResponseFromRepayment(result.Repayment),
			Loan:                  loanResponseFromLoan(result.Loan),
			SettlementOperationID: result.SettlementOperationID,
			SettlementTxHash:      result.SettlementTxHash,
		})
	})
}

func handleListRepayments(logger *slog.Logger, service LoansService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}

		items, err := service.ListRepayments(r.Context(), id)
		if err != nil {
			if writeLoanError(w, err) {
				return
			}
			logger.ErrorContext(r.Context(), "list repayments failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("list repayments failed"))
			return
		}

		repayments := make([]repaymentResponse, 0, len(items))
		for _, item := range items {
			repayments = append(repayments, repaymentResponseFromRepayment(item))
		}
		writeJSON(w, http.StatusOK, map[string]any{"repayments": repayments})
	})
}

type loanResponse struct {
	ID               int64  `json:"id"`
	TokenID          string `json:"token_id,omitempty"`
	BorrowerRef      string `json:"borrower_ref"`
	LenderAddress    string `json:"lender_address"`
	PrincipalMinor   int64  `json:"principal_minor"`
	APRBps           int32  `json:"apr_bps"`
	TermDays         int64  `json:"term_days"`
	InterestDueMinor int64  `json:"interest_due_minor"`
	TotalDueMinor    int64  `json:"total_due_minor"`
	OutstandingMinor int64  `json:"outstanding_minor"`
	Status           string `json:"status"`
	OperationID      int64  `json:"operation_id,omitempty"`
	TxHash           string `json:"tx_hash,omitempty"`
	OwnerAddress     string `json:"owner_address,omitempty"`
}

type repaymentResultResponse struct {
	Repayment             repaymentResponse `json:"repayment"`
	Loan                  loanResponse      `json:"loan"`
	SettlementOperationID int64             `json:"settlement_operation_id,omitempty"`
	SettlementTxHash      string            `json:"settlement_tx_hash,omitempty"`
}

type repaymentResponse struct {
	ID          int64     `json:"id"`
	LoanID      int64     `json:"loan_id"`
	AmountMinor int64     `json:"amount_minor"`
	ExternalRef string    `json:"external_ref,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func repaymentResponseFromRepayment(repayment db.Repayment) repaymentResponse {
	return repaymentResponse{
		ID:          repayment.ID,
		LoanID:      repayment.LoanID,
		AmountMinor: repayment.AmountMinor,
		ExternalRef: stringFromNull(repayment.ExternalRef),
		CreatedAt:   repayment.CreatedAt,
	}
}

func loanResponseWithTx(loan db.Loan, operationID int64, txHash string) loanResponse {
	response := loanResponseFromLoan(loan)
	response.OperationID = operationID
	response.TxHash = txHash
	return response
}

func loanResponseFromRead(result loans.ReadResult) loanResponse {
	response := loanResponseFromLoan(result.Loan)
	if result.Loan.MintOperationID != nil {
		response.OperationID = *result.Loan.MintOperationID
	}
	response.TxHash = result.MintTxHash
	response.OwnerAddress = result.OwnerAddress
	return response
}

func loanResponseFromLoan(loan db.Loan) loanResponse {
	return loanResponse{
		ID:               loan.ID,
		TokenID:          stringFromNull(loan.TokenID),
		BorrowerRef:      loan.BorrowerRef,
		LenderAddress:    loan.LenderAddress,
		PrincipalMinor:   loan.PrincipalMinor,
		APRBps:           loan.AprBps,
		TermDays:         loan.TermDays,
		InterestDueMinor: loan.InterestDueMinor,
		TotalDueMinor:    loan.TotalDueMinor,
		OutstandingMinor: loan.OutstandingMinor,
		Status:           loan.Status,
	}
}

func stringFromNull(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}

type termsResponse struct {
	PrincipalMinor   int64  `json:"principal_minor"`
	APRBps           uint16 `json:"apr_bps"`
	TermDays         int64  `json:"term_days"`
	InterestDueMinor int64  `json:"interest_due_minor"`
	TotalDueMinor    int64  `json:"total_due_minor"`
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, errorBody("invalid id"))
		return 0, false
	}
	return id, true
}

func writeLoanError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeJSON(w, http.StatusNotFound, errorBody("loan not found"))
	case errors.Is(err, loans.ErrInvalidAmount),
		errors.Is(err, loans.ErrInvalidTerm),
		errors.Is(err, loans.ErrInvalidAddress),
		errors.Is(err, loans.ErrOverpayment):
		writeJSON(w, http.StatusBadRequest, errorBody(err.Error()))
	case errors.Is(err, loans.ErrLoanNotActive),
		errors.Is(err, loans.ErrLoanNotTransferable),
		errors.Is(err, loans.ErrLoanMissingToken),
		errors.Is(err, loans.ErrLoanMissingContract),
		errors.Is(err, loans.ErrNoActiveContract),
		errors.Is(err, loans.ErrDuplicateExternalRef):
		writeJSON(w, http.StatusConflict, errorBody(err.Error()))
	case errors.Is(err, db.ErrLockBusy):
		writeJSON(w, http.StatusServiceUnavailable, errorBody("another chain operation is in progress, retry shortly"))
	default:
		return false
	}
	return true
}
