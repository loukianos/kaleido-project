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

// LocalIssuer marks identities created before OIDC lands; the OIDC slice replaces it with the token issuer.
const LocalIssuer = "local"

const RoleLender = "lender"

var (
	ErrNoCustodialKey   = errors.New("no custodial key for this address")
	ErrUnknownEncryptor = errors.New("signing key sealed by an unknown encryptor")
	ErrInvalidSubject   = errors.New("subject must not be empty")
)

const uniqueViolationSQLState = "23505"

// Store is the subset of generated queries the service needs; *db.Queries satisfies it.
type Store interface {
	GetOrCreateIdentity(context.Context, db.GetOrCreateIdentityParams) (db.Identity, error)
	GetIdentityByID(context.Context, int64) (db.Identity, error)
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

// ResolveLender returns the lender identity for subject and the signer for its custodial key, provisioning both on first sight.
// On a network that charges gas, provisioning is where a funding transfer would happen; the local Besu network is gas-free, so none is needed.
func (s *Service) ResolveLender(ctx context.Context, subject string) (db.Identity, *eth.Signer, error) {
	if subject == "" {
		return db.Identity{}, nil, ErrInvalidSubject
	}
	ident, err := s.store.GetOrCreateIdentity(ctx, db.GetOrCreateIdentityParams{
		Issuer:  LocalIssuer,
		Subject: subject,
		Role:    RoleLender,
	})
	if err != nil {
		return db.Identity{}, nil, fmt.Errorf("get or create identity: %w", err)
	}

	signer, err := s.signerForIdentity(ctx, ident.ID)
	if err != nil {
		return db.Identity{}, nil, err
	}
	return ident, signer, nil
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

func (s *Service) signerForIdentity(ctx context.Context, identityID int64) (*eth.Signer, error) {
	row, err := s.store.GetSigningKeyByIdentityID(ctx, identityID)
	if errors.Is(err, pgx.ErrNoRows) {
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
		Encryptor:  keys.AESGCMScheme,
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
	if row.Encryptor != keys.AESGCMScheme {
		return nil, fmt.Errorf("%w: %s", ErrUnknownEncryptor, row.Encryptor)
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
