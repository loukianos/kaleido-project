package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
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

// handleCreateLoan originates a loan and mints its note on chain.
//
//	@Summary	Originate a loan note
//	@Tags		loans
//	@Accept		json
//	@Produce	json
//	@Param		request	body		createLoanRequest	true	"Loan to originate"
//	@Success	201		{object}	loanResponse
//	@Failure	400		{object}	errorResponse
//	@Failure	409		{object}	errorResponse	"No active contract deployed"
//	@Failure	503		{object}	errorResponse	"Another chain operation is in progress, retry shortly"
//	@Router		/loans [post]
func handleCreateLoan(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createLoanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorBody("invalid json body"))
			return
		}
		if req.BorrowerRef == "" || req.LenderAddress == "" {
			c.JSON(http.StatusBadRequest, errorBody("borrower_ref and lender_address are required"))
			return
		}

		result, err := service.Originate(c.Request.Context(), loans.OriginateRequest{
			BorrowerRef:    req.BorrowerRef,
			LenderAddress:  req.LenderAddress,
			PrincipalMinor: req.PrincipalMinor,
			APRBps:         req.APRBps,
			TermDays:       req.TermDays,
		})
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "create loan failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("create loan failed"))
			return
		}

		c.JSON(http.StatusCreated, loanResponseWithTx(result.Loan, result.OperationID, result.TxHash))
	}
}

// handleGetLoan reads one loan, decorated with its owner and mint tx hash.
//
//	@Summary	Get loan
//	@Tags		loans
//	@Produce	json
//	@Param		id	path		int	true	"Loan ID"
//	@Success	200	{object}	loanResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/loans/{id} [get]
func handleGetLoan(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		result, err := service.Get(c.Request.Context(), id)
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "get loan failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("get loan failed"))
			return
		}
		c.JSON(http.StatusOK, loanResponseFromRead(result))
	}
}

type loansListResponse struct {
	Loans []loanResponse `json:"loans"`
}

// handleListLoans lists loans with optional filters and offset paging.
//
//	@Summary	List loans
//	@Tags		loans
//	@Produce	json
//	@Param		lender	query		string	false	"Filter by lender address (case-insensitive)"
//	@Param		status	query		string	false	"Filter by loan status"	Enums(originating, active, settling, repaid, defaulted)
//	@Param		limit	query		int		false	"Page size"				minimum(1)	maximum(100)	default(50)
//	@Param		offset	query		int		false	"Result offset"			minimum(0)
//	@Success	200		{object}	loansListResponse
//	@Failure	400		{object}	errorResponse
//	@Router		/loans [get]
func handleListLoans(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Zero means "not specified". The service applies its default.
		limit := int32(0)
		if raw := c.Query("limit"); raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 32)
			if err != nil || parsed < 1 || parsed > 100 {
				c.JSON(http.StatusBadRequest, errorBody("limit must be between 1 and 100"))
				return
			}
			limit = int32(parsed)
		}
		offset := int32(0)
		if raw := c.Query("offset"); raw != "" {
			parsed, err := strconv.ParseInt(raw, 10, 32)
			if err != nil || parsed < 0 {
				c.JSON(http.StatusBadRequest, errorBody("offset must be non-negative"))
				return
			}
			offset = int32(parsed)
		}

		items, err := service.List(c.Request.Context(), loans.ListRequest{
			Lender: c.Query("lender"),
			Status: c.Query("status"),
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "list loans failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("list loans failed"))
			return
		}

		responses := make([]loanResponse, 0, len(items))
		for _, item := range items {
			responses = append(responses, loanResponseFromLoan(item))
		}
		c.JSON(http.StatusOK, loansListResponse{Loans: responses})
	}
}

// handleLoanTerms serves the terms JSON targeted by the note's tokenURI.
//
//	@Summary	Loan terms for tokenURI
//	@Tags		loans
//	@Produce	json
//	@Param		id	path		int	true	"Loan ID"
//	@Success	200	{object}	termsResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/loans/{id}/terms [get]
func handleLoanTerms(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		terms, err := service.Terms(c.Request.Context(), id)
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "get loan terms failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("get loan terms failed"))
			return
		}
		c.JSON(http.StatusOK, termsResponse{
			PrincipalMinor:   terms.PrincipalMinor,
			APRBps:           terms.APRBps,
			TermDays:         terms.TermDays,
			InterestDueMinor: terms.InterestDueMinor,
			TotalDueMinor:    terms.TotalDueMinor,
		})
	}
}

type transferLoanRequest struct {
	ToAddress string `json:"to_address"`
}

// handleTransferLoan reassigns the note to a new lender on chain.
//
//	@Summary	Transfer loan note
//	@Tags		loans
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int					true	"Loan ID"
//	@Param		request	body		transferLoanRequest	true	"Transfer target"
//	@Success	200		{object}	loanResponse
//	@Failure	400		{object}	errorResponse
//	@Failure	404		{object}	errorResponse
//	@Failure	409		{object}	errorResponse	"Loan is not transferable"
//	@Failure	503		{object}	errorResponse	"Another chain operation is in progress, retry shortly"
//	@Router		/loans/{id}/transfer [post]
func handleTransferLoan(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		var req transferLoanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		result, err := service.Transfer(c.Request.Context(), id, loans.TransferRequest{ToAddress: req.ToAddress})
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "transfer loan failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("transfer loan failed"))
			return
		}

		c.JSON(http.StatusOK, loanResponseWithTx(result.Loan, result.OperationID, result.TxHash))
	}
}

// handleDefaultLoan marks an active loan defaulted on chain and in the API.
//
//	@Summary	Mark loan defaulted
//	@Tags		loans
//	@Produce	json
//	@Param		id	path		int	true	"Loan ID"
//	@Success	200	{object}	loanResponse
//	@Failure	404	{object}	errorResponse
//	@Failure	409	{object}	errorResponse	"Loan is not active"
//	@Failure	503	{object}	errorResponse	"Another chain operation is in progress, retry shortly"
//	@Router		/loans/{id}/default [post]
func handleDefaultLoan(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		result, err := service.Default(c.Request.Context(), id)
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "default loan failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("default loan failed"))
			return
		}

		c.JSON(http.StatusOK, loanResponseWithTx(result.Loan, result.OperationID, result.TxHash))
	}
}

type createRepaymentRequest struct {
	AmountMinor int64  `json:"amount_minor"`
	ExternalRef string `json:"external_ref"`
}

// handleCreateRepayment records a repayment; paying the balance to zero
// settles the loan and burns the note.
//
//	@Summary	Record repayment
//	@Tags		loans
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"Loan ID"
//	@Param		request	body		createRepaymentRequest	true	"Repayment to record"
//	@Success	201		{object}	repaymentResultResponse
//	@Failure	400		{object}	errorResponse	"Invalid amount or overpayment"
//	@Failure	404		{object}	errorResponse
//	@Failure	409		{object}	errorResponse	"Loan not active or duplicate external_ref"
//	@Failure	503		{object}	errorResponse	"Another chain operation is in progress, retry shortly"
//	@Router		/loans/{id}/repayments [post]
func handleCreateRepayment(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		var req createRepaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		result, err := service.RecordRepayment(c.Request.Context(), id, loans.RepaymentRequest{
			AmountMinor: req.AmountMinor,
			ExternalRef: req.ExternalRef,
		})
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "create repayment failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("create repayment failed"))
			return
		}

		c.JSON(http.StatusCreated, repaymentResultResponse{
			Repayment:             repaymentResponseFromRepayment(result.Repayment),
			Loan:                  loanResponseFromLoan(result.Loan),
			SettlementOperationID: result.SettlementOperationID,
			SettlementTxHash:      result.SettlementTxHash,
		})
	}
}

type repaymentsListResponse struct {
	Repayments []repaymentResponse `json:"repayments"`
}

// handleListRepayments lists a loan's repayments oldest-first.
//
//	@Summary	List repayments
//	@Tags		loans
//	@Produce	json
//	@Param		id	path		int	true	"Loan ID"
//	@Success	200	{object}	repaymentsListResponse
//	@Failure	404	{object}	errorResponse
//	@Router		/loans/{id}/repayments [get]
func handleListRepayments(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		items, err := service.ListRepayments(c.Request.Context(), id)
		if err != nil {
			if writeLoanError(c, err) {
				return
			}
			logger.ErrorContext(c.Request.Context(), "list repayments failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("list repayments failed"))
			return
		}

		repayments := make([]repaymentResponse, 0, len(items))
		for _, item := range items {
			repayments = append(repayments, repaymentResponseFromRepayment(item))
		}
		c.JSON(http.StatusOK, repaymentsListResponse{Repayments: repayments})
	}
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

func pathID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, errorBody("invalid id"))
		return 0, false
	}
	return id, true
}

func writeLoanError(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		c.JSON(http.StatusNotFound, errorBody("loan not found"))
	case errors.Is(err, loans.ErrInvalidAmount),
		errors.Is(err, loans.ErrInvalidTerm),
		errors.Is(err, loans.ErrInvalidAddress),
		errors.Is(err, loans.ErrOverpayment):
		c.JSON(http.StatusBadRequest, errorBody(err.Error()))
	case errors.Is(err, loans.ErrLoanNotActive),
		errors.Is(err, loans.ErrLoanNotTransferable),
		errors.Is(err, loans.ErrNotNoteOwner),
		errors.Is(err, loans.ErrLoanMissingToken),
		errors.Is(err, loans.ErrLoanMissingContract),
		errors.Is(err, loans.ErrNoActiveContract),
		errors.Is(err, loans.ErrDuplicateExternalRef):
		c.JSON(http.StatusConflict, errorBody(err.Error()))
	case errors.Is(err, db.ErrLockBusy):
		c.JSON(http.StatusServiceUnavailable, errorBody("another chain operation is in progress, retry shortly"))
	default:
		return false
	}
	return true
}
