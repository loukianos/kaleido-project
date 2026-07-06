// Package chainop tracks single chain writes through the chain_operations
// journal while holding the shared writer lock. Both the contracts and loans
// services submit transactions through it.
package chainop

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/eth"
)

const (
	lockTTL        = 30 * time.Second
	receiptTimeout = 2 * time.Minute
)

var (
	// ErrPending marks transient failures whose operation was journaled retryable: the platform owns the retry, so callers report 202 rather than an error.
	ErrPending = errors.New("operation pending retry")
	// ErrReverted marks a transaction the chain mined and rejected; retrying is futile, so the operation is terminally failed.
	ErrReverted = errors.New("transaction reverted on chain")
)

// pendingError matches ErrPending while keeping the cause reachable through errors.Is and errors.Unwrap.
type pendingError struct{ cause error }

func (e *pendingError) Error() string        { return ErrPending.Error() + ": " + e.cause.Error() }
func (e *pendingError) Unwrap() error        { return e.cause }
func (e *pendingError) Is(target error) bool { return target == ErrPending }

// holderSeq distinguishes concurrent acquisitions within one process. The
// lock manager treats an acquisition by the current holder as a lease
// renewal, so if every request shared the process-wide holder name, two
// concurrent submissions would both "hold" the lock and race on the nonce.
// The counter is package-level so all submitters in the process draw from
// the same sequence.
var holderSeq atomic.Uint64

// Journal records chain-operation status transitions. Both service
// repositories satisfy it.
type Journal interface {
	SetOperationSubmitted(context.Context, db.SetOperationSubmittedParams) (db.ChainOperation, error)
	SetOperationMined(context.Context, int64) (db.ChainOperation, error)
	SetOperationRetryable(context.Context, db.SetOperationRetryableParams) (db.ChainOperation, error)
	SetOperationFailed(context.Context, db.SetOperationFailedParams) (db.ChainOperation, error)
}

type Submitter struct {
	chain   eth.Writer
	locks   *db.LockManager
	holder  string
	journal Journal
}

func NewSubmitter(chain eth.Writer, locks *db.LockManager, holder string, journal Journal) *Submitter {
	return &Submitter{
		chain:   chain,
		locks:   locks,
		holder:  holder,
		journal: journal,
	}
}

// Submit performs a single chain write signed by the platform's default signer, tracking the operation through submitted and mined.
func (s *Submitter) Submit(
	ctx context.Context,
	opID int64,
	action string,
	send func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	return s.SubmitAs(ctx, s.chain.DefaultSigner(), opID, action, send)
}

// SubmitAs performs a single chain write signed by the given signer, tracking the operation through submitted and mined.
// The writer lock is per signing address, so writes by different identities never contend; it is released after the transaction is submitted but before the receipt is awaited.
func (s *Submitter) SubmitAs(
	ctx context.Context,
	signer *eth.Signer,
	opID int64,
	action string,
	send func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	holder := fmt.Sprintf("%s#%d", s.holder, holderSeq.Add(1))
	release, err := s.locks.Acquire(ctx, eth.LockNameFor(s.chain.ChainID(), signer.Address()), holder, lockTTL)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}
	return s.submitLocked(ctx, signer, release, opID, action, send)
}

// poolSeq rotates the starting pool index across submissions so load spreads over the pool even without contention.
var poolSeq atomic.Uint64

// SubmitAsAny performs a single chain write signed by whichever pool signer's lock is free, falling through busy ones.
// Every pool member must be equally authorized for the action; when all locks are busy the operation is retryable with ErrLockBusy, matching single-signer behavior.
func (s *Submitter) SubmitAsAny(
	ctx context.Context,
	pool []*eth.Signer,
	opID int64,
	action string,
	send func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	if len(pool) == 0 {
		return "", nil, s.Retryable(ctx, opID, fmt.Errorf("%s: signer pool is empty", action))
	}

	holder := fmt.Sprintf("%s#%d", s.holder, holderSeq.Add(1))
	start := int(poolSeq.Add(1)) % len(pool)
	for i := range pool {
		signer := pool[(start+i)%len(pool)]
		release, err := s.locks.Acquire(ctx, eth.LockNameFor(s.chain.ChainID(), signer.Address()), holder, lockTTL)
		if errors.Is(err, db.ErrLockBusy) {
			continue
		}
		if err != nil {
			return "", nil, s.Retryable(ctx, opID, err)
		}
		return s.submitLocked(ctx, signer, release, opID, action, send)
	}
	return "", nil, s.Retryable(ctx, opID, db.ErrLockBusy)
}

// submitLocked runs the write while holding the signer's lock, releasing it once the transaction is in the pending pool.
func (s *Submitter) submitLocked(
	ctx context.Context,
	signer *eth.Signer,
	release func(context.Context) error,
	opID int64,
	action string,
	send func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	defer func() { _ = release(ctx) }()

	backend, err := s.chain.Backend()
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}

	nonce, err := s.chain.PendingNonceOf(ctx, signer.Address())
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}
	auth, err := signer.TransactOpts(ctx, s.chain.ChainID(), nonce)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}

	tx, err := send(auth, backend)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, fmt.Errorf("%s on chain: %w", action, err))
	}
	if _, err := s.journal.SetOperationSubmitted(ctx, db.SetOperationSubmittedParams{
		ID:            opID,
		TxHash:        db.Ptr(tx.Hash().Hex()),
		Nonce:         db.Ptr(int64(nonce)),
		SignerAddress: db.Ptr(signer.Address().Hex()),
	}); err != nil {
		return "", nil, fmt.Errorf("mark %s submitted: %w", action, err)
	}

	// The transaction is now in the pending pool; other writers can safely take the lock.
	// A failed release is harmless: the lock expires on its own after the TTL.
	_ = release(ctx)

	receiptCtx, cancel := context.WithTimeout(ctx, receiptTimeout)
	defer cancel()
	receipt, err := bind.WaitMined(receiptCtx, backend, tx)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, fmt.Errorf("wait for %s receipt: %w", action, err))
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		// A revert is permanent: the chain evaluated the transaction and rejected it, so retrying the same call is futile.
		return "", nil, s.Failed(ctx, opID, fmt.Errorf("%s %w: %s", action, ErrReverted, tx.Hash()))
	}
	if _, err := s.journal.SetOperationMined(ctx, opID); err != nil {
		return "", nil, fmt.Errorf("mark %s mined: %w", action, err)
	}
	return tx.Hash().Hex(), receipt, nil
}

// Retryable records err against the operation so the reconciler can retry it, and wraps it in ErrPending so callers report the operation as in flight.
func (s *Submitter) Retryable(ctx context.Context, opID int64, err error) error {
	_, _ = s.journal.SetOperationRetryable(ctx, db.SetOperationRetryableParams{
		ID:    opID,
		Error: db.Ptr(err.Error()),
	})
	return &pendingError{cause: err}
}

// Failed terminally fails the operation, then returns err unchanged.
func (s *Submitter) Failed(ctx context.Context, opID int64, err error) error {
	_, _ = s.journal.SetOperationFailed(ctx, db.SetOperationFailedParams{
		ID:    opID,
		Error: db.Ptr(err.Error()),
	})
	return err
}
