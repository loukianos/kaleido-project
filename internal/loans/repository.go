package loans

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	db "kaleido-project/db/sqlc"
)

const operationStatusCreated = "created"

type Repository struct {
	db      db.TxBeginner
	queries *db.Queries
}

type CreateOriginationParams struct {
	ContractID       int64
	BorrowerRef      string
	LenderAddress    string
	LenderIdentityID *int64
	PrincipalMinor   int64
	APRBps           int32
	TermDays         int64
	InterestDueMinor int64
	TotalDueMinor    int64
	ExternalRef      *string
}

type RepaymentTxResult struct {
	Loan                db.Loan
	Repayment           db.Repayment
	SettlementOperation db.ChainOperation
}

func NewRepository(queries *db.Queries, conn db.TxBeginner) *Repository {
	return &Repository{queries: queries, db: conn}
}

func (r *Repository) ActiveContract(ctx context.Context, chainID int64) (db.Contract, error) {
	return r.queries.GetActiveContractByChainID(ctx, chainID)
}

func (r *Repository) Contract(ctx context.Context, id int64) (db.Contract, error) {
	return r.queries.GetContractByID(ctx, id)
}

func (r *Repository) Loan(ctx context.Context, id int64) (db.Loan, error) {
	return r.queries.GetLoanByID(ctx, id)
}

func (r *Repository) ListLoans(ctx context.Context, params db.ListLoansParams) ([]db.Loan, error) {
	return r.queries.ListLoans(ctx, params)
}

func (r *Repository) Repayments(ctx context.Context, loanID int64) ([]db.Repayment, error) {
	return r.queries.ListRepaymentsByLoan(ctx, loanID)
}

func (r *Repository) Operation(ctx context.Context, id int64) (db.ChainOperation, error) {
	return r.queries.GetChainOperationByID(ctx, id)
}

func (r *Repository) CreateLoanOperation(ctx context.Context, kind string, contractID int64, loanID int64) (db.ChainOperation, error) {
	return r.queries.CreateChainOperation(ctx, db.CreateChainOperationParams{
		Kind:       kind,
		Status:     operationStatusCreated,
		ContractID: db.Ptr(contractID),
		LoanID:     db.Ptr(loanID),
	})
}

func (r *Repository) CreateOrigination(ctx context.Context, params CreateOriginationParams) (db.Loan, db.ChainOperation, error) {
	var (
		loan db.Loan
		op   db.ChainOperation
	)
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		var err error
		loan, err = q.CreateLoan(ctx, db.CreateLoanParams{
			BorrowerRef:      params.BorrowerRef,
			LenderAddress:    params.LenderAddress,
			LenderIdentityID: params.LenderIdentityID,
			PrincipalMinor:   params.PrincipalMinor,
			AprBps:           params.APRBps,
			TermDays:         params.TermDays,
			InterestDueMinor: params.InterestDueMinor,
			TotalDueMinor:    params.TotalDueMinor,
			OutstandingMinor: params.TotalDueMinor,
			Status:           LoanStatusOriginating,
			ContractID:       db.Ptr(params.ContractID),
			ExternalRef:      params.ExternalRef,
		})
		if err != nil {
			return fmt.Errorf("create loan: %w", err)
		}

		op, err = q.CreateChainOperation(ctx, db.CreateChainOperationParams{
			Kind:       originateOperationKind,
			Status:     operationStatusCreated,
			ContractID: db.Ptr(params.ContractID),
			LoanID:     db.Ptr(loan.ID),
		})
		if err != nil {
			return fmt.Errorf("create originate operation: %w", err)
		}
		return nil
	})
	if err != nil {
		return db.Loan{}, db.ChainOperation{}, err
	}
	return loan, op, nil
}

func (r *Repository) ApplyOrigination(ctx context.Context, loanID int64, tokenID string, opID int64, contractID int64) (db.Loan, error) {
	var loan db.Loan
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		var err error
		loan, err = q.SetLoanActive(ctx, db.SetLoanActiveParams{
			ID:              loanID,
			TokenID:         db.Ptr(tokenID),
			MintOperationID: db.Ptr(opID),
			ContractID:      db.Ptr(contractID),
		})
		if err != nil {
			return fmt.Errorf("mark loan active: %w", err)
		}
		if _, err := q.SetOperationApplied(ctx, opID); err != nil {
			return fmt.Errorf("mark originate applied: %w", err)
		}
		return nil
	})
	if err != nil {
		return db.Loan{}, err
	}
	return loan, nil
}

func (r *Repository) RecordRepayment(ctx context.Context, loanID int64, req RepaymentRequest) (RepaymentTxResult, error) {
	var result RepaymentTxResult
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		loan, err := q.GetLoanByIDForUpdate(ctx, loanID)
		if err != nil {
			return err
		}
		if loan.Status != LoanStatusActive {
			return ErrLoanNotActive
		}

		outstanding, err := ApplyRepayment(loan.OutstandingMinor, req.AmountMinor)
		if err != nil {
			return err
		}

		repayment, err := q.CreateRepayment(ctx, db.CreateRepaymentParams{
			LoanID:      loan.ID,
			AmountMinor: req.AmountMinor,
			ExternalRef: nullableString(req.ExternalRef),
		})
		if err != nil {
			if isUniqueViolation(err, "repayments_external_ref_unique") {
				return ErrDuplicateExternalRef
			}
			return fmt.Errorf("create repayment: %w", err)
		}

		nextStatus := LoanStatusActive
		var settleOp db.ChainOperation
		if outstanding == 0 {
			nextStatus = LoanStatusSettling
			settleOp, err = createSettleOperation(ctx, q, loan)
			if err != nil {
				return err
			}
		}

		loan, err = q.UpdateLoanOutstandingAndStatus(ctx, db.UpdateLoanOutstandingAndStatusParams{
			ID:               loan.ID,
			OutstandingMinor: outstanding,
			Status:           nextStatus,
		})
		if err != nil {
			return fmt.Errorf("update loan outstanding: %w", err)
		}

		result = RepaymentTxResult{Loan: loan, Repayment: repayment, SettlementOperation: settleOp}
		return nil
	})
	if err != nil {
		return RepaymentTxResult{}, err
	}
	return result, nil
}

func createSettleOperation(ctx context.Context, q *db.Queries, loan db.Loan) (db.ChainOperation, error) {
	if loan.TokenID == nil {
		return db.ChainOperation{}, ErrLoanMissingToken
	}
	if loan.ContractID == nil {
		return db.ChainOperation{}, ErrLoanMissingContract
	}
	return q.CreateChainOperation(ctx, db.CreateChainOperationParams{
		Kind:       settleOperationKind,
		Status:     operationStatusCreated,
		ContractID: loan.ContractID,
		LoanID:     db.Ptr(loan.ID),
	})
}

func (r *Repository) ApplySettlement(ctx context.Context, loanID int64, opID int64) (db.Loan, error) {
	return r.setLoanStatusAndApplyOperation(ctx, loanID, LoanStatusRepaid, opID, "settle")
}

func (r *Repository) ApplyDefault(ctx context.Context, loanID int64, opID int64) (db.Loan, error) {
	return r.setLoanStatusAndApplyOperation(ctx, loanID, LoanStatusDefaulted, opID, "default")
}

func (r *Repository) ApplyTransfer(ctx context.Context, loanID int64, lenderAddress string, lenderIdentityID *int64, opID int64) (db.Loan, error) {
	var loan db.Loan
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		var err error
		loan, err = q.UpdateLoanLender(ctx, db.UpdateLoanLenderParams{
			ID:               loanID,
			LenderAddress:    lenderAddress,
			LenderIdentityID: lenderIdentityID,
		})
		if err != nil {
			return fmt.Errorf("update loan lender: %w", err)
		}
		if _, err := q.SetOperationApplied(ctx, opID); err != nil {
			return fmt.Errorf("mark transfer applied: %w", err)
		}
		return nil
	})
	if err != nil {
		return db.Loan{}, err
	}
	return loan, nil
}

func (r *Repository) setLoanStatusAndApplyOperation(ctx context.Context, loanID int64, status string, opID int64, action string) (db.Loan, error) {
	var loan db.Loan
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		var err error
		loan, err = q.SetLoanStatus(ctx, db.SetLoanStatusParams{
			ID:     loanID,
			Status: status,
		})
		if err != nil {
			return fmt.Errorf("mark loan %s: %w", status, err)
		}
		if _, err := q.SetOperationApplied(ctx, opID); err != nil {
			return fmt.Errorf("mark %s applied: %w", action, err)
		}
		return nil
	})
	if err != nil {
		return db.Loan{}, err
	}
	return loan, nil
}

func (r *Repository) SetOperationSubmitted(ctx context.Context, params db.SetOperationSubmittedParams) (db.ChainOperation, error) {
	return r.queries.SetOperationSubmitted(ctx, params)
}

func (r *Repository) SetOperationMined(ctx context.Context, id int64) (db.ChainOperation, error) {
	return r.queries.SetOperationMined(ctx, id)
}

func (r *Repository) SetOperationRetryable(ctx context.Context, params db.SetOperationRetryableParams) (db.ChainOperation, error) {
	return r.queries.SetOperationRetryable(ctx, params)
}

func (r *Repository) SetOperationFailed(ctx context.Context, params db.SetOperationFailedParams) (db.ChainOperation, error) {
	return r.queries.SetOperationFailed(ctx, params)
}

func (r *Repository) SetOperationApplied(ctx context.Context, id int64) (db.ChainOperation, error) {
	return r.queries.SetOperationApplied(ctx, id)
}

// FailLoan terminally fails a loan whose chain operation cannot complete.
func (r *Repository) FailLoan(ctx context.Context, loanID int64) (db.Loan, error) {
	return r.queries.SetLoanStatus(ctx, db.SetLoanStatusParams{ID: loanID, Status: LoanStatusFailed})
}

func (r *Repository) LoanByExternalRef(ctx context.Context, externalRef string) (db.Loan, error) {
	return r.queries.GetLoanByExternalRef(ctx, db.Ptr(externalRef))
}

func (r *Repository) RetryableOperations(ctx context.Context, maxAttempts int32, limit int32) ([]db.ChainOperation, error) {
	return r.queries.ListRetryableOperations(ctx, db.ListRetryableOperationsParams{MaxAttempts: maxAttempts, LimitCount: limit})
}

func (r *Repository) ExhaustedOperations(ctx context.Context, maxAttempts int32, limit int32) ([]db.ChainOperation, error) {
	return r.queries.ListExhaustedOperations(ctx, db.ListExhaustedOperationsParams{MaxAttempts: maxAttempts, LimitCount: limit})
}

func (r *Repository) StaleSubmittedOperations(ctx context.Context, staleBefore time.Time, limit int32) ([]db.ChainOperation, error) {
	return r.queries.ListStaleSubmittedOperations(ctx, db.ListStaleSubmittedOperationsParams{StaleBefore: staleBefore, LimitCount: limit})
}

// IdentityIDByAddress resolves a custodial key's identity, or nil when the address isn't custodied.
func (r *Repository) IdentityIDByAddress(ctx context.Context, address string) (*int64, error) {
	row, err := r.queries.GetSigningKeyByAddress(ctx, address)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return db.Ptr(row.IdentityID), nil
}

// nullableString maps "" to NULL, so absent refs don't collide in the unique index on (loan_id, external_ref).
func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	// 23505 is the Postgres error code for unique_violation.
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
