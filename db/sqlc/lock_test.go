package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLockManagerAcquireAndRelease(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	fake := &fakeLockQueries{}
	manager := NewLockManager(fake)
	manager.now = func() time.Time { return now }

	release, err := manager.Acquire(ctx, "ethereum-writer:1337:0xabc", "api-1", 30*time.Second)
	require.NoError(t, err)
	require.Equal(t, "ethereum-writer:1337:0xabc", fake.acquire.Name)
	require.Equal(t, "api-1", fake.acquire.Holder)
	require.True(t, fake.acquire.ExpiresAt.Equal(now.Add(30*time.Second)))

	require.NoError(t, release(ctx))
	require.Equal(t, fake.acquire.Name, fake.release.Name)
	require.Equal(t, fake.acquire.Holder, fake.release.Holder)
}

func TestLockManagerReleaseIsIdempotent(t *testing.T) {
	ctx := context.Background()
	fake := &fakeLockQueries{}
	manager := NewLockManager(fake)

	release, err := manager.Acquire(ctx, "lock", "holder", time.Second)
	require.NoError(t, err)

	require.NoError(t, release(ctx))
	require.NoError(t, release(ctx))
	require.Equal(t, 1, fake.releaseCalls)
}

func TestLockManagerAcquireBusy(t *testing.T) {
	manager := NewLockManager(&fakeLockQueries{acquireErr: pgx.ErrNoRows})

	_, err := manager.Acquire(context.Background(), "lock", "holder", time.Second)
	require.ErrorIs(t, err, ErrLockBusy)
}

func TestLockManagerAcquireWrapsUnexpectedError(t *testing.T) {
	dbErr := errors.New("db down")
	manager := NewLockManager(&fakeLockQueries{acquireErr: dbErr})

	_, err := manager.Acquire(context.Background(), "lock", "holder", time.Second)
	require.ErrorIs(t, err, dbErr)
	require.Contains(t, err.Error(), `acquire lock "lock"`)
}

func TestLockManagerAcquireInvalid(t *testing.T) {
	tests := map[string]struct {
		manager *LockManager
		name    string
		holder  string
		ttl     time.Duration
	}{
		"nil manager": {
			manager: nil,
			name:    "lock",
			holder:  "holder",
			ttl:     time.Second,
		},
		"nil queries": {
			manager: &LockManager{},
			name:    "lock",
			holder:  "holder",
			ttl:     time.Second,
		},
		"empty name": {
			manager: NewLockManager(&fakeLockQueries{}),
			holder:  "holder",
			ttl:     time.Second,
		},
		"empty holder": {
			manager: NewLockManager(&fakeLockQueries{}),
			name:    "lock",
			ttl:     time.Second,
		},
		"zero ttl": {
			manager: NewLockManager(&fakeLockQueries{}),
			name:    "lock",
			holder:  "holder",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := tc.manager.Acquire(context.Background(), tc.name, tc.holder, tc.ttl)
			require.ErrorIs(t, err, ErrInvalidLock)
		})
	}
}

type fakeLockQueries struct {
	acquire      AcquireAppLockParams
	acquireErr   error
	release      ReleaseAppLockParams
	releaseErr   error
	releaseCalls int
}

func (f *fakeLockQueries) AcquireAppLock(_ context.Context, arg AcquireAppLockParams) (AppLock, error) {
	f.acquire = arg
	if f.acquireErr != nil {
		return AppLock{}, f.acquireErr
	}
	return AppLock{
		Name:      arg.Name,
		Holder:    arg.Holder,
		ExpiresAt: arg.ExpiresAt,
	}, nil
}

func (f *fakeLockQueries) ReleaseAppLock(_ context.Context, arg ReleaseAppLockParams) error {
	f.release = arg
	f.releaseCalls++
	return f.releaseErr
}
