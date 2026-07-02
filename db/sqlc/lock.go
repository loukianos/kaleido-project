package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrLockBusy    = errors.New("lock is held by another process")
	ErrInvalidLock = errors.New("lock name, holder, and ttl are required")
)

type LockQueries interface {
	AcquireAppLock(context.Context, AcquireAppLockParams) (AppLock, error)
	ReleaseAppLock(context.Context, ReleaseAppLockParams) error
}

type LockManager struct {
	queries LockQueries
	now     func() time.Time
}

func NewLockManager(queries LockQueries) *LockManager {
	return &LockManager{
		queries: queries,
		now:     time.Now,
	}
}

func (m *LockManager) Acquire(ctx context.Context, name string, holder string, ttl time.Duration) (func(context.Context) error, error) {
	if m == nil || m.queries == nil || name == "" || holder == "" || ttl <= 0 {
		return nil, ErrInvalidLock
	}

	expiresAt := m.now().UTC().Add(ttl)
	lock, err := m.queries.AcquireAppLock(ctx, AcquireAppLockParams{
		Name:      name,
		Holder:    holder,
		ExpiresAt: expiresAt,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrLockBusy
	}
	if err != nil {
		return nil, fmt.Errorf("acquire lock %q: %w", name, err)
	}
	if lock.Holder != holder {
		return nil, ErrLockBusy
	}

	released := false
	return func(ctx context.Context) error {
		if released {
			return nil
		}
		released = true
		if err := m.queries.ReleaseAppLock(ctx, ReleaseAppLockParams{
			Name:   name,
			Holder: holder,
		}); err != nil {
			return fmt.Errorf("release lock %q: %w", name, err)
		}
		return nil
	}, nil
}
