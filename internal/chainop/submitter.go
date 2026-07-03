// Package chainop tracks single chain writes through the chain_operations
// journal while holding the shared writer lock. Both the contracts and loans
// services submit transactions through it.
package chainop

import (
	"context"
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

// Submit performs a single chain write, tracking the operation through submitted and mined.
// Shared writer lock is released after the transaction is submitted but before the receipt is awaited.
func (s *Submitter) Submit(
	ctx context.Context,
	opID int64,
	action string,
	send func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	holder := fmt.Sprintf("%s#%d", s.holder, holderSeq.Add(1))
	release, err := s.locks.Acquire(ctx, eth.LockName(s.chain), holder, lockTTL)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}
	defer func() { _ = release(ctx) }()

	backend, err := s.chain.Backend()
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}

	nonce, err := s.chain.PendingNonce(ctx)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}
	auth, err := s.chain.TransactOpts(ctx, nonce)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, err)
	}

	tx, err := send(auth, backend)
	if err != nil {
		return "", nil, s.Retryable(ctx, opID, fmt.Errorf("%s on chain: %w", action, err))
	}
	if _, err := s.journal.SetOperationSubmitted(ctx, db.SetOperationSubmittedParams{
		ID:     opID,
		TxHash: db.Ptr(tx.Hash().Hex()),
		Nonce:  db.Ptr(int64(nonce)),
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
		return "", nil, s.Retryable(ctx, opID, fmt.Errorf("%s transaction failed: %s", action, tx.Hash()))
	}
	if _, err := s.journal.SetOperationMined(ctx, opID); err != nil {
		return "", nil, fmt.Errorf("mark %s mined: %w", action, err)
	}
	return tx.Hash().Hex(), receipt, nil
}

// Retryable records err against the operation so it can be retried later, then returns err unchanged.
func (s *Submitter) Retryable(ctx context.Context, opID int64, err error) error {
	_, _ = s.journal.SetOperationRetryable(ctx, db.SetOperationRetryableParams{
		ID:    opID,
		Error: db.Ptr(err.Error()),
	})
	return err
}
