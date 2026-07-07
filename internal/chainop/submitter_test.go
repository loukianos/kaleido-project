package chainop

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/eth"
)

func TestSubmitSuccess(t *testing.T) {
	tx := newTestTx()
	journal := &fakeJournal{}
	locks := &fakeLockQueries{}
	backend := fakeBackend{
		receipt: &types.Receipt{
			Status: types.ReceiptStatusSuccessful,
			TxHash: tx.Hash(),
		},
	}

	releasedBeforeReceipt := false
	backend.onReceipt = func() { releasedBeforeReceipt = locks.releaseCalls == 1 }
	writer := &fakeWriter{backend: backend}
	submitter := NewSubmitter(writer, db.NewLockManager(locks), "test-holder", journal)

	var gotNonce *big.Int
	txHash, receipt, err := submitter.Submit(context.Background(), 42, "test",
		func(auth *bind.TransactOpts, _ eth.ContractBackend) (*types.Transaction, error) {
			gotNonce = auth.Nonce
			return tx, nil
		})

	require.NoError(t, err)
	require.Equal(t, tx.Hash().Hex(), txHash)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	require.Equal(t, int64(7), gotNonce.Int64())

	require.NotNil(t, journal.submitted)
	require.Equal(t, int64(42), journal.submitted.ID)
	require.Equal(t, tx.Hash().Hex(), *journal.submitted.TxHash)
	require.Equal(t, int64(7), *journal.submitted.Nonce)
	require.Equal(t, int64(42), journal.minedID)
	require.Nil(t, journal.retryable)
	require.Equal(t, 1, locks.releaseCalls)
	require.True(t, releasedBeforeReceipt)
}

func TestSubmitLockBusy(t *testing.T) {
	journal := &fakeJournal{}
	locks := &fakeLockQueries{acquireErr: pgx.ErrNoRows}
	submitter := newTestSubmitter(journal, locks, nil)

	_, _, err := submitter.Submit(context.Background(), 42, "test", neverSend(t))

	require.ErrorIs(t, err, db.ErrLockBusy)
	require.Nil(t, journal.submitted)
	require.NotNil(t, journal.retryable)
	require.Equal(t, int64(42), journal.retryable.ID)
}

func TestSubmitSendFailure(t *testing.T) {
	journal := &fakeJournal{}
	locks := &fakeLockQueries{}
	submitter := newTestSubmitter(journal, locks, nil)

	sendErr := errors.New("execution reverted")
	_, _, err := submitter.Submit(context.Background(), 42, "test",
		func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) {
			return nil, sendErr
		})

	require.ErrorIs(t, err, sendErr)
	require.Contains(t, err.Error(), "test on chain")
	require.Nil(t, journal.submitted)
	require.NotNil(t, journal.retryable)
	require.Contains(t, *journal.retryable.Error, "execution reverted")
	require.Equal(t, 1, locks.releaseCalls)
}

func TestSubmitRevertedTransaction(t *testing.T) {
	tx := newTestTx()
	journal := &fakeJournal{}
	locks := &fakeLockQueries{}
	submitter := newTestSubmitter(journal, locks, &types.Receipt{
		Status: types.ReceiptStatusFailed,
		TxHash: tx.Hash(),
	})

	_, _, err := submitter.Submit(context.Background(), 42, "test",
		func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) {
			return tx, nil
		})

	// A revert is permanent: the operation fails terminally instead of queueing for retry.
	require.ErrorIs(t, err, ErrReverted)
	require.NotErrorIs(t, err, ErrPending)
	require.NotNil(t, journal.submitted)
	require.Zero(t, journal.minedID)
	require.Nil(t, journal.retryable)
	require.NotNil(t, journal.failed)
	require.Equal(t, int64(42), journal.failed.ID)
}

func TestSubmitUsesUniqueHolderPerAcquisition(t *testing.T) {
	tx := newTestTx()
	journal := &fakeJournal{}
	locks := &fakeLockQueries{}
	submitter := newTestSubmitter(journal, locks, &types.Receipt{
		Status: types.ReceiptStatusSuccessful,
		TxHash: tx.Hash(),
	})

	send := func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) {
		return tx, nil
	}
	_, _, err := submitter.Submit(context.Background(), 1, "test", send)
	require.NoError(t, err)
	_, _, err = submitter.Submit(context.Background(), 2, "test", send)
	require.NoError(t, err)

	require.Len(t, locks.acquireHolders, 2)
	require.NotEqual(t, locks.acquireHolders[0], locks.acquireHolders[1])
	for _, holder := range locks.acquireHolders {
		require.Contains(t, holder, "test-holder#")
	}
}

func TestRetryableRecordsAndReturnsError(t *testing.T) {
	journal := &fakeJournal{}
	submitter := newTestSubmitter(journal, &fakeLockQueries{}, nil)

	cause := errors.New("boom")
	err := submitter.Retryable(context.Background(), 42, cause)

	// The returned error matches ErrPending for 202 mapping while keeping the cause reachable.
	require.ErrorIs(t, err, ErrPending)
	require.ErrorIs(t, err, cause)
	require.NotNil(t, journal.retryable)
	require.Equal(t, int64(42), journal.retryable.ID)
	require.Equal(t, "boom", *journal.retryable.Error)
}

func newTestSubmitter(journal *fakeJournal, locks *fakeLockQueries, receipt *types.Receipt) *Submitter {
	writer := &fakeWriter{backend: fakeBackend{receipt: receipt}}
	return NewSubmitter(writer, db.NewLockManager(locks), "test-holder", journal)
}

func newTestTx() *types.Transaction {
	to := common.HexToAddress("0x2222222222222222222222222222222222222222")
	return types.NewTx(&types.LegacyTx{Nonce: 7, To: &to, Gas: 21000, GasPrice: big.NewInt(0), Value: big.NewInt(0)})
}

func neverSend(t *testing.T) func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) {
	return func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) {
		t.Fatal("send must not be called")
		return nil, nil
	}
}

// testSignerKey is deterministic test-only material so the fake writer can produce real transact opts.
var testSignerKey = strings.Repeat("1", 64)

// testPool builds n distinct signers from deterministic throwaway keys.
func testPool(t *testing.T, n int) []*eth.Signer {
	t.Helper()
	pool := make([]*eth.Signer, 0, n)
	for i := 1; i <= n; i++ {
		key := strings.Repeat(fmt.Sprintf("%02d", i), 32)
		signer, err := eth.NewSigner(key)
		require.NoError(t, err)
		pool = append(pool, signer)
	}
	return pool
}

func TestSubmitAsAnyFallsThroughBusyLocks(t *testing.T) {
	tx := newTestTx()
	journal := &fakeJournal{}
	pool := testPool(t, 3)
	writer := &fakeWriter{backend: fakeBackend{receipt: &types.Receipt{Status: types.ReceiptStatusSuccessful, TxHash: tx.Hash()}}}
	// Every pool lock except the second signer's is held elsewhere.
	locks := &fakeLockQueries{busyNames: map[string]bool{
		eth.LockNameFor(writer.ChainID(), pool[0].Address()): true,
		eth.LockNameFor(writer.ChainID(), pool[2].Address()): true,
	}}
	submitter := NewSubmitter(writer, db.NewLockManager(locks), "test-holder", journal)

	txHash, _, err := submitter.SubmitAsAny(context.Background(), pool, 42, "test",
		func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) { return tx, nil })

	require.NoError(t, err)
	require.Equal(t, tx.Hash().Hex(), txHash)
	require.NotNil(t, journal.submitted)
	require.Equal(t, pool[1].Address().Hex(), *journal.submitted.SignerAddress)
}

func TestSubmitAsAnyAllBusy(t *testing.T) {
	journal := &fakeJournal{}
	pool := testPool(t, 2)
	writer := &fakeWriter{}
	locks := &fakeLockQueries{busyNames: map[string]bool{
		eth.LockNameFor(writer.ChainID(), pool[0].Address()): true,
		eth.LockNameFor(writer.ChainID(), pool[1].Address()): true,
	}}
	submitter := NewSubmitter(writer, db.NewLockManager(locks), "test-holder", journal)

	_, _, err := submitter.SubmitAsAny(context.Background(), pool, 42, "test", neverSend(t))

	require.ErrorIs(t, err, db.ErrLockBusy)
	require.NotNil(t, journal.retryable)
	// Both locks were tried before giving up.
	require.Len(t, locks.acquireNames, 2)
}

func TestSubmitAsAnyEmptyPool(t *testing.T) {
	journal := &fakeJournal{}
	submitter := newTestSubmitter(journal, &fakeLockQueries{}, nil)

	_, _, err := submitter.SubmitAsAny(context.Background(), nil, 42, "test", neverSend(t))

	require.Error(t, err)
	require.NotNil(t, journal.retryable)
}

type fakeWriter struct {
	backend eth.ContractBackend
}

func (f *fakeWriter) DefaultSigner() *eth.Signer {
	signer, err := eth.NewSigner(testSignerKey)
	if err != nil {
		panic(err)
	}
	return signer
}

func (f *fakeWriter) SignerAddress() common.Address {
	return f.DefaultSigner().Address()
}

func (f *fakeWriter) ChainID() *big.Int { return big.NewInt(1337) }

func (f *fakeWriter) PendingNonceOf(context.Context, common.Address) (uint64, error) { return 7, nil }

func (f *fakeWriter) Backend() (eth.ContractBackend, error) { return f.backend, nil }

// fakeBackend embeds the interface so only the methods the submitter actually uses need real implementations.
type fakeBackend struct {
	eth.ContractBackend
	receipt   *types.Receipt
	onReceipt func()
}

func (f fakeBackend) TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error) {
	if f.onReceipt != nil {
		f.onReceipt()
	}
	return f.receipt, nil
}

type fakeJournal struct {
	submitted *db.SetOperationSubmittedParams
	minedID   int64
	retryable *db.SetOperationRetryableParams
	failed    *db.SetOperationFailedParams
}

func (f *fakeJournal) SetOperationFailed(_ context.Context, arg db.SetOperationFailedParams) (db.ChainOperation, error) {
	f.failed = &arg
	return db.ChainOperation{ID: arg.ID}, nil
}

func (f *fakeJournal) SetOperationSubmitted(_ context.Context, arg db.SetOperationSubmittedParams) (db.ChainOperation, error) {
	f.submitted = &arg
	return db.ChainOperation{ID: arg.ID}, nil
}

func (f *fakeJournal) SetOperationMined(_ context.Context, id int64) (db.ChainOperation, error) {
	f.minedID = id
	return db.ChainOperation{ID: id}, nil
}

func (f *fakeJournal) SetOperationRetryable(_ context.Context, arg db.SetOperationRetryableParams) (db.ChainOperation, error) {
	f.retryable = &arg
	return db.ChainOperation{ID: arg.ID}, nil
}

type fakeLockQueries struct {
	acquireErr     error
	busyNames      map[string]bool
	acquireHolders []string
	acquireNames   []string
	releaseCalls   int
}

func (f *fakeLockQueries) AcquireAppLock(_ context.Context, arg db.AcquireAppLockParams) (db.AppLock, error) {
	f.acquireHolders = append(f.acquireHolders, arg.Holder)
	f.acquireNames = append(f.acquireNames, arg.Name)
	if f.acquireErr != nil {
		return db.AppLock{}, f.acquireErr
	}
	if f.busyNames[arg.Name] {
		return db.AppLock{}, pgx.ErrNoRows
	}
	return db.AppLock{Name: arg.Name, Holder: arg.Holder, ExpiresAt: arg.ExpiresAt}, nil
}

func (f *fakeLockQueries) ReleaseAppLock(context.Context, db.ReleaseAppLockParams) error {
	f.releaseCalls++
	return nil
}
