// Package identity maps subjects to custodial signing keys.
// A lender's key is provisioned lazily the first time an operation needs their address; the private key is envelope-encrypted at rest and decrypted only for the duration of a request.
package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/eth"
	"kaleido-project/internal/keys"
)

const (
	RoleLender   = "lender"
	RoleServicer = "servicer"
	// PlatformIssuer marks platform-internal identities like servicer pool keys, which no external token ever names.
	PlatformIssuer = "platform"
)

var (
	ErrNoCustodialKey   = errors.New("no custodial key for this address")
	ErrUnknownEncryptor = errors.New("signing key sealed by an unknown encryptor")
	ErrInvalidSubject   = errors.New("subject must not be empty")
	ErrNotOnboarded     = errors.New("lender has not onboarded")
)

const uniqueViolationSQLState = "23505"

// Store is the subset of generated queries the service needs; *db.Queries satisfies it.
type Store interface {
	GetOrCreateIdentity(context.Context, db.GetOrCreateIdentityParams) (db.Identity, error)
	GetIdentityByID(context.Context, int64) (db.Identity, error)
	GetIdentityByIssuerSubject(context.Context, db.GetIdentityByIssuerSubjectParams) (db.Identity, error)
	GetSigningKeyByIdentityID(context.Context, int64) (db.SigningKey, error)
	GetSigningKeyByAddress(context.Context, string) (db.SigningKey, error)
	CreateSigningKey(context.Context, db.CreateSigningKeyParams) (db.SigningKey, error)
}

type Service struct {
	store     Store
	encryptor keys.Encryptor
}

func NewService(store Store, encryptor keys.Encryptor) *Service {
	return &Service{store: store, encryptor: encryptor}
}

// ResolveIdentity returns the lender identity for (issuer, subject), creating it on first sight without provisioning a key.
func (s *Service) ResolveIdentity(ctx context.Context, issuer, subject string) (db.Identity, error) {
	if issuer == "" || subject == "" {
		return db.Identity{}, ErrInvalidSubject
	}
	ident, err := s.store.GetOrCreateIdentity(ctx, db.GetOrCreateIdentityParams{
		Issuer:  issuer,
		Subject: subject,
		Role:    RoleLender,
	})
	if err != nil {
		return db.Identity{}, fmt.Errorf("get or create identity: %w", err)
	}
	return ident, nil
}

// OnboardLender is the explicit onboarding step: it creates the lender identity for (issuer, subject) if new and eagerly provisions its custodial wallet.
// On a network that charges gas, this is where a funding transfer would happen; the local Besu network is gas-free, so none is needed.
// Idempotent: an already-onboarded lender gets their existing identity and key back.
func (s *Service) OnboardLender(ctx context.Context, issuer, subject string) (db.Identity, *eth.Signer, error) {
	ident, err := s.ResolveIdentity(ctx, issuer, subject)
	if err != nil {
		return db.Identity{}, nil, err
	}

	signer, err := s.signerForIdentity(ctx, ident.ID, true)
	if err != nil {
		return db.Identity{}, nil, err
	}
	return ident, signer, nil
}

// LenderAddress resolves a subject named in a request body to an onboarded lender's identity and custodial address.
// Unlike OnboardLender it never creates anything: naming a lender who hasn't onboarded is ErrNotOnboarded, so provisioning strictly precedes participation.
func (s *Service) LenderAddress(ctx context.Context, issuer, subject string) (db.Identity, common.Address, error) {
	if issuer == "" || subject == "" {
		return db.Identity{}, common.Address{}, ErrInvalidSubject
	}
	ident, err := s.store.GetIdentityByIssuerSubject(ctx, db.GetIdentityByIssuerSubjectParams{
		Issuer:  issuer,
		Subject: subject,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Identity{}, common.Address{}, fmt.Errorf("%w: %s", ErrNotOnboarded, subject)
	}
	if err != nil {
		return db.Identity{}, common.Address{}, fmt.Errorf("get identity: %w", err)
	}

	row, err := s.store.GetSigningKeyByIdentityID(ctx, ident.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		// An identity row without a wallet has been seen but never onboarded.
		return db.Identity{}, common.Address{}, fmt.Errorf("%w: %s", ErrNotOnboarded, subject)
	}
	if err != nil {
		return db.Identity{}, common.Address{}, fmt.Errorf("get signing key: %w", err)
	}
	return ident, common.HexToAddress(row.Address), nil
}

// SignerForIdentity returns the signer for an identity's existing custodial key, or ErrNoCustodialKey when none was ever provisioned.
func (s *Service) SignerForIdentity(ctx context.Context, identityID int64) (*eth.Signer, error) {
	return s.signerForIdentity(ctx, identityID, false)
}

// EnsureServicerPool provisions size platform-internal servicer identities with signing keys and returns their signers.
// Pool keys sign servicer chain writes so they don't serialize on one nonce sequence; they never hold assets.
// Idempotent: existing pool identities get their existing keys back, so the pool is stable across restarts.
func (s *Service) EnsureServicerPool(ctx context.Context, size int) ([]*eth.Signer, error) {
	signers := make([]*eth.Signer, 0, size)
	for i := 1; i <= size; i++ {
		ident, err := s.store.GetOrCreateIdentity(ctx, db.GetOrCreateIdentityParams{
			Issuer:  PlatformIssuer,
			Subject: fmt.Sprintf("servicer-pool-%d", i),
			Role:    RoleServicer,
		})
		if err != nil {
			return nil, fmt.Errorf("get or create pool identity %d: %w", i, err)
		}
		signer, err := s.signerForIdentity(ctx, ident.ID, true)
		if err != nil {
			return nil, fmt.Errorf("provision pool key %d: %w", i, err)
		}
		signers = append(signers, signer)
	}
	return signers, nil
}

// SignerForAddress returns the signer for the custodial key holding address, or ErrNoCustodialKey when we don't hold it.
func (s *Service) SignerForAddress(ctx context.Context, address common.Address) (*eth.Signer, error) {
	row, err := s.store.GetSigningKeyByAddress(ctx, address.Hex())
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoCustodialKey
	}
	if err != nil {
		return nil, fmt.Errorf("get signing key by address: %w", err)
	}
	return s.decryptSigner(ctx, row)
}

// Identity returns the identity row by id, for read paths that surface the subject.
func (s *Service) Identity(ctx context.Context, id int64) (db.Identity, error) {
	return s.store.GetIdentityByID(ctx, id)
}

func (s *Service) signerForIdentity(ctx context.Context, identityID int64, provision bool) (*eth.Signer, error) {
	row, err := s.store.GetSigningKeyByIdentityID(ctx, identityID)
	if errors.Is(err, pgx.ErrNoRows) {
		if !provision {
			return nil, ErrNoCustodialKey
		}
		return s.provisionKey(ctx, identityID)
	}
	if err != nil {
		return nil, fmt.Errorf("get signing key: %w", err)
	}
	return s.decryptSigner(ctx, row)
}

func (s *Service) provisionKey(ctx context.Context, identityID int64) (*eth.Signer, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate signing key: %w", err)
	}
	ciphertext, version, err := s.encryptor.Encrypt(ctx, crypto.FromECDSA(key))
	if err != nil {
		return nil, fmt.Errorf("encrypt signing key: %w", err)
	}

	signer := eth.NewSignerFromKey(key)
	if _, err := s.store.CreateSigningKey(ctx, db.CreateSigningKeyParams{
		IdentityID: identityID,
		Address:    signer.Address().Hex(),
		Ciphertext: ciphertext,
		Encryptor:  s.encryptor.Scheme(),
		KeyVersion: int32(version),
	}); err != nil {
		// A concurrent request provisioned this identity's key first; use theirs and discard ours.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationSQLState {
			row, getErr := s.store.GetSigningKeyByIdentityID(ctx, identityID)
			if getErr != nil {
				return nil, fmt.Errorf("get signing key after conflict: %w", getErr)
			}
			return s.decryptSigner(ctx, row)
		}
		return nil, fmt.Errorf("store signing key: %w", err)
	}
	return signer, nil
}

func (s *Service) decryptSigner(ctx context.Context, row db.SigningKey) (*eth.Signer, error) {
	if row.Encryptor != s.encryptor.Scheme() {
		return nil, fmt.Errorf("%w: row sealed by %s, service configured for %s", ErrUnknownEncryptor, row.Encryptor, s.encryptor.Scheme())
	}
	plaintext, err := s.encryptor.Decrypt(ctx, row.Ciphertext, int(row.KeyVersion))
	if err != nil {
		return nil, err
	}
	key, err := crypto.ToECDSA(plaintext)
	if err != nil {
		return nil, fmt.Errorf("parse decrypted signing key: %w", err)
	}
	return eth.NewSignerFromKey(key), nil
}
