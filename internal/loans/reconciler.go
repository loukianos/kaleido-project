package loans

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/chainop"
	contractpkg "kaleido-project/internal/contracts"
)

const (
	reconcileMaxAttempts = 5
	reconcileStaleAfter  = 2 * time.Minute
	reconcileBatchSize   = 20
	reconcileLockTTL     = 60 * time.Second
)

// Reconciler drains the chain-operation journal in the background: it re-drives retryable platform operations, resolves stale submitted transactions by hash, and terminally fails what can't converge.
// One instance runs at a time across the fleet, elected per tick through the shared lock manager.
type Reconciler struct {
	service *Service
	locks   *db.LockManager
	holder  string
	logger  *slog.Logger
}

func NewReconciler(service *Service, locks *db.LockManager, holder string, logger *slog.Logger) *Reconciler {
	return &Reconciler{service: service, locks: locks, holder: holder, logger: logger}
}

// Run drains the journal every interval until ctx is done.
func (r *Reconciler) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.tick(ctx)
		}
	}
}

func (r *Reconciler) tick(ctx context.Context) {
	lockName := fmt.Sprintf("reconciler:%s", r.service.chain.ChainID())
	release, err := r.locks.Acquire(ctx, lockName, r.holder, reconcileLockTTL)
	if errors.Is(err, db.ErrLockBusy) {
		// Another instance is reconciling this tick.
		return
	}
	if err != nil {
		r.logger.WarnContext(ctx, "reconciler leader election failed", "error", err)
		return
	}
	defer func() { _ = release(ctx) }()

	if err := r.DrainOnce(ctx); err != nil {
		r.logger.WarnContext(ctx, "reconcile pass failed", "error", err)
	}
}

// DrainOnce runs one reconcile pass: exhausted operations fail terminally, stale submitted ones resolve by receipt, and retryable ones re-drive.
func (r *Reconciler) DrainOnce(ctx context.Context) error {
	if err := r.failExhausted(ctx); err != nil {
		return err
	}
	if err := r.resolveStaleSubmitted(ctx); err != nil {
		return err
	}
	return r.redriveRetryable(ctx)
}

func (r *Reconciler) failExhausted(ctx context.Context) error {
	ops, err := r.service.repo.ExhaustedOperations(ctx, reconcileMaxAttempts, reconcileBatchSize)
	if err != nil {
		return fmt.Errorf("list exhausted operations: %w", err)
	}
	for _, op := range ops {
		r.logger.WarnContext(ctx, "operation exhausted retries", "operation", op.ID, "kind", op.Kind)
		if err := r.failOperation(ctx, op, "retry attempts exhausted"); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reconciler) resolveStaleSubmitted(ctx context.Context) error {
	ops, err := r.service.repo.StaleSubmittedOperations(ctx, time.Now().UTC().Add(-reconcileStaleAfter), reconcileBatchSize)
	if err != nil {
		return fmt.Errorf("list stale submitted operations: %w", err)
	}
	for _, op := range ops {
		if err := r.resolveByReceipt(ctx, op); err != nil {
			return err
		}
	}
	return nil
}

// resolveByReceipt settles the fate of a transaction whose receipt was never observed: mined and successful means apply, mined and reverted means fail, dropped means retry.
func (r *Reconciler) resolveByReceipt(ctx context.Context, op db.ChainOperation) error {
	if op.TxHash == nil {
		return r.failOperation(ctx, op, "submitted without a transaction hash")
	}
	backend, err := r.service.chain.Backend()
	if err != nil {
		return fmt.Errorf("get backend: %w", err)
	}

	receipt, err := backend.TransactionReceipt(ctx, common.HexToHash(*op.TxHash))
	if errors.Is(err, ethereum.NotFound) {
		// The pending pool dropped it; queue a fresh submission for re-drivable kinds.
		r.logger.InfoContext(ctx, "submitted transaction dropped", "operation", op.ID, "kind", op.Kind, "tx", *op.TxHash)
		if !redrivable(op.Kind) {
			return r.failOperation(ctx, op, "transaction dropped and operation is not re-drivable")
		}
		_, journalErr := r.service.repo.SetOperationRetryable(ctx, db.SetOperationRetryableParams{
			ID:    op.ID,
			Error: db.Ptr("transaction dropped from the pending pool; will resubmit"),
		})
		return journalErr
	}
	if err != nil {
		// The chain is unreachable; leave the operation for the next pass.
		r.logger.WarnContext(ctx, "receipt lookup failed", "operation", op.ID, "error", err)
		return nil
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return r.failOperation(ctx, op, "transaction reverted on chain")
	}

	r.logger.InfoContext(ctx, "recovered mined transaction", "operation", op.ID, "kind", op.Kind, "tx", *op.TxHash)
	return r.applyMined(ctx, op, receipt)
}

// applyMined applies a mined-but-unapplied operation's effects to the loan ledger.
func (r *Reconciler) applyMined(ctx context.Context, op db.ChainOperation, receipt *types.Receipt) error {
	if _, err := r.service.repo.SetOperationMined(ctx, op.ID); err != nil {
		return fmt.Errorf("mark operation %d mined: %w", op.ID, err)
	}

	switch op.Kind {
	case originateOperationKind:
		note, err := r.bindContract(ctx, op)
		if err != nil {
			return err
		}
		tokenID, err := parseOriginatedTokenID(note, receipt)
		if err != nil {
			return r.failOperation(ctx, op, "mined origination missing LoanOriginated event")
		}
		_, err = r.service.repo.ApplyOrigination(ctx, *op.LoanID, tokenID.String(), op.ID, *op.ContractID)
		return err
	case settleOperationKind:
		_, err := r.service.repo.ApplySettlement(ctx, *op.LoanID, op.ID)
		return err
	case markDefaultedOperationKind:
		_, err := r.service.repo.ApplyDefault(ctx, *op.LoanID, op.ID)
		return err
	case transferOperationKind:
		note, err := r.bindContract(ctx, op)
		if err != nil {
			return err
		}
		to, err := parseTransferDestination(note, receipt)
		if err != nil {
			return r.failOperation(ctx, op, "mined transfer missing Transfer event")
		}
		identityID, err := r.service.repo.IdentityIDByAddress(ctx, to.Hex())
		if err != nil {
			return fmt.Errorf("resolve transfer destination identity: %w", err)
		}
		_, err = r.service.repo.ApplyTransfer(ctx, *op.LoanID, to.Hex(), identityID, op.ID)
		return err
	case grantRoleOperationKind:
		_, err := r.service.repo.SetOperationApplied(ctx, op.ID)
		return err
	default:
		// Deploys carry state (address, base URI) the journal doesn't; an admin redeploys instead.
		return r.failOperation(ctx, op, "kind cannot be reconciled automatically")
	}
}

func (r *Reconciler) redriveRetryable(ctx context.Context) error {
	ops, err := r.service.repo.RetryableOperations(ctx, reconcileMaxAttempts, reconcileBatchSize)
	if err != nil {
		return fmt.Errorf("list retryable operations: %w", err)
	}
	for _, op := range ops {
		if err := r.redrive(ctx, op); err != nil {
			return err
		}
	}
	return nil
}

// redrive re-runs a retryable platform operation from durable state.
// Transient failures leave it retryable with attempts incremented; the exhausted pass eventually gives up.
func (r *Reconciler) redrive(ctx context.Context, op db.ChainOperation) error {
	if !redrivable(op.Kind) {
		// Transfers aren't re-signed without a fresh request, and deploys carry unjournaled state.
		return r.failOperation(ctx, op, "kind is not re-drivable")
	}
	if op.LoanID == nil || op.ContractID == nil {
		return r.failOperation(ctx, op, "operation is missing loan or contract linkage")
	}

	loan, err := r.service.repo.Loan(ctx, *op.LoanID)
	if err != nil {
		return fmt.Errorf("load loan %d: %w", *op.LoanID, err)
	}
	contract, err := r.service.repo.Contract(ctx, *op.ContractID)
	if err != nil {
		return fmt.Errorf("load contract %d: %w", *op.ContractID, err)
	}

	r.logger.InfoContext(ctx, "re-driving operation", "operation", op.ID, "kind", op.Kind, "loan", loan.ID, "attempts", op.Attempts)

	switch {
	case op.Kind == originateOperationKind && loan.Status == LoanStatusOriginating:
		_, _, err = r.service.driveOrigination(ctx, loan, contract, op.ID)
	case op.Kind == settleOperationKind && loan.Status == LoanStatusSettling:
		if _, err = r.service.settle(ctx, loan, op); err == nil {
			_, err = r.service.repo.ApplySettlement(ctx, loan.ID, op.ID)
		}
	case op.Kind == markDefaultedOperationKind && loan.Status == LoanStatusActive:
		tokenID, parseErr := parseTokenID(loan)
		if parseErr != nil {
			return r.failOperation(ctx, op, parseErr.Error())
		}
		if _, _, err = r.service.submitOperation(ctx, op.ID, contract.Address, "default",
			func(auth *bind.TransactOpts, note *contractpkg.LoanNote) (*types.Transaction, error) {
				return note.MarkDefaulted(auth, tokenID)
			}); err == nil {
			_, err = r.service.repo.ApplyDefault(ctx, loan.ID, op.ID)
		}
	default:
		// The loan moved on (or terminally failed) while this operation waited; retrying would fight the current state.
		return r.failOperation(ctx, op, fmt.Sprintf("loan is %s; operation no longer applies", loan.Status))
	}

	if errors.Is(err, chainop.ErrPending) {
		// Still transient; attempts were incremented and the next pass tries again.
		return nil
	}
	if errors.Is(err, chainop.ErrReverted) {
		return r.failLoanForOperation(ctx, op)
	}
	return err
}

// failOperation terminally fails an operation and, where the operation carries the loan's fate, the loan.
func (r *Reconciler) failOperation(ctx context.Context, op db.ChainOperation, reason string) error {
	if _, err := r.service.repo.SetOperationFailed(ctx, db.SetOperationFailedParams{ID: op.ID, Error: db.Ptr(reason)}); err != nil {
		return fmt.Errorf("fail operation %d: %w", op.ID, err)
	}
	return r.failLoanForOperation(ctx, op)
}

// failLoanForOperation fails the loan when the dead operation was load-bearing for its state: a loan can't become active without its mint, or repaid without its burn.
// A failed default leaves the loan active, which remains true.
func (r *Reconciler) failLoanForOperation(ctx context.Context, op db.ChainOperation) error {
	if op.LoanID == nil {
		return nil
	}
	if op.Kind != originateOperationKind && op.Kind != settleOperationKind {
		return nil
	}
	if _, err := r.service.repo.FailLoan(ctx, *op.LoanID); err != nil {
		return fmt.Errorf("fail loan %d: %w", *op.LoanID, err)
	}
	return nil
}

func (r *Reconciler) bindContract(ctx context.Context, op db.ChainOperation) (*contractpkg.LoanNote, error) {
	if op.ContractID == nil {
		return nil, fmt.Errorf("operation %d has no contract", op.ID)
	}
	contract, err := r.service.repo.Contract(ctx, *op.ContractID)
	if err != nil {
		return nil, fmt.Errorf("load contract %d: %w", *op.ContractID, err)
	}
	backend, err := r.service.chain.Backend()
	if err != nil {
		return nil, err
	}
	return contractpkg.NewLoanNote(common.HexToAddress(contract.Address), backend)
}

func redrivable(kind string) bool {
	switch kind {
	case originateOperationKind, settleOperationKind, markDefaultedOperationKind:
		return true
	}
	return false
}

func parseTransferDestination(note *contractpkg.LoanNote, receipt *types.Receipt) (common.Address, error) {
	for _, log := range receipt.Logs {
		if event, err := note.ParseTransfer(*log); err == nil {
			return event.To, nil
		}
	}
	return common.Address{}, errors.New("transfer event missing from receipt")
}
