package identity

import (
	"context"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/keys"
)

const (
	testMasterKey = "56659eba0b040dfd7503c987c0c26428a476b3ee49747a181f6f7361e253e4a3"
	testIssuerURL = "http://issuer.test/realms/loan-notes"
)

func newTestService(t *testing.T) (*Service, *fakeStore) {
	t.Helper()
	encryptor, err := keys.NewAESGCM(testMasterKey)
	require.NoError(t, err)
	store := &fakeStore{identities: map[string]db.Identity{}, keysByIdentity: map[int64]db.SigningKey{}}
	return NewService(store, encryptor), store
}

func TestResolveLenderProvisionsOnFirstSight(t *testing.T) {
	service, store := newTestService(t)

	ident, signer, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)
	require.Equal(t, "alice", ident.Subject)
	require.Equal(t, testIssuerURL, ident.Issuer)
	require.Equal(t, RoleLender, ident.Role)
	require.NotNil(t, signer)

	row := store.keysByIdentity[ident.ID]
	require.Equal(t, signer.Address().Hex(), row.Address)
	require.Equal(t, keys.AESGCMScheme, row.Encryptor)
	// The stored ciphertext must not be the raw key.
	require.NotContains(t, string(row.Ciphertext), signer.Address().Hex())
}

func TestResolveLenderReusesExistingKey(t *testing.T) {
	service, _ := newTestService(t)

	_, first, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)
	_, second, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)
	require.Equal(t, first.Address(), second.Address())

	_, other, err := service.ResolveLender(context.Background(), testIssuerURL, "bob")
	require.NoError(t, err)
	require.NotEqual(t, first.Address(), other.Address())
}

func TestResolveLenderRejectsEmptySubject(t *testing.T) {
	service, _ := newTestService(t)
	_, _, err := service.ResolveLender(context.Background(), testIssuerURL, "")
	require.ErrorIs(t, err, ErrInvalidSubject)
	_, _, err = service.ResolveLender(context.Background(), "", "alice")
	require.ErrorIs(t, err, ErrInvalidSubject)
}

func TestResolveIdentityDoesNotProvisionKey(t *testing.T) {
	service, store := newTestService(t)

	ident, err := service.ResolveIdentity(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)
	require.Empty(t, store.keysByIdentity)

	_, err = service.SignerForIdentity(context.Background(), ident.ID)
	require.ErrorIs(t, err, ErrNoCustodialKey)
}

func TestSignerForIdentityReturnsProvisionedKey(t *testing.T) {
	service, _ := newTestService(t)

	ident, signer, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)

	found, err := service.SignerForIdentity(context.Background(), ident.ID)
	require.NoError(t, err)
	require.Equal(t, signer.Address(), found.Address())
}

func TestResolveLenderSurvivesProvisionRace(t *testing.T) {
	service, store := newTestService(t)
	store.conflictOnce = true

	_, signer, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)
	// The winner's key (installed by the fake during the conflict) is used, not the loser's.
	require.Equal(t, store.keysByIdentity[store.identities[testIssuerURL+"/alice"].ID].Address, signer.Address().Hex())
}

func TestSignerForAddress(t *testing.T) {
	service, _ := newTestService(t)

	_, signer, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)

	found, err := service.SignerForAddress(context.Background(), signer.Address())
	require.NoError(t, err)
	require.Equal(t, signer.Address(), found.Address())

	_, err = service.SignerForAddress(context.Background(), common.HexToAddress("0x2222222222222222222222222222222222222222"))
	require.ErrorIs(t, err, ErrNoCustodialKey)
}

func TestDecryptRejectsUnknownEncryptor(t *testing.T) {
	service, store := newTestService(t)

	ident, _, err := service.ResolveLender(context.Background(), testIssuerURL, "alice")
	require.NoError(t, err)

	row := store.keysByIdentity[ident.ID]
	row.Encryptor = "aws-kms"
	store.keysByIdentity[ident.ID] = row

	_, err = service.SignerForAddress(context.Background(), common.HexToAddress(row.Address))
	require.ErrorIs(t, err, ErrUnknownEncryptor)
}

type fakeStore struct {
	identities     map[string]db.Identity
	keysByIdentity map[int64]db.SigningKey
	nextID         int64
	conflictOnce   bool
}

func (f *fakeStore) GetOrCreateIdentity(_ context.Context, arg db.GetOrCreateIdentityParams) (db.Identity, error) {
	key := arg.Issuer + "/" + arg.Subject
	if ident, ok := f.identities[key]; ok {
		return ident, nil
	}
	f.nextID++
	ident := db.Identity{ID: f.nextID, Issuer: arg.Issuer, Subject: arg.Subject, Role: arg.Role}
	f.identities[key] = ident
	return ident, nil
}

func (f *fakeStore) GetIdentityByID(_ context.Context, id int64) (db.Identity, error) {
	for _, ident := range f.identities {
		if ident.ID == id {
			return ident, nil
		}
	}
	return db.Identity{}, pgx.ErrNoRows
}

func (f *fakeStore) GetSigningKeyByIdentityID(_ context.Context, identityID int64) (db.SigningKey, error) {
	row, ok := f.keysByIdentity[identityID]
	if !ok {
		return db.SigningKey{}, pgx.ErrNoRows
	}
	return row, nil
}

func (f *fakeStore) GetSigningKeyByAddress(_ context.Context, address string) (db.SigningKey, error) {
	for _, row := range f.keysByIdentity {
		if strings.EqualFold(row.Address, address) {
			return row, nil
		}
	}
	return db.SigningKey{}, pgx.ErrNoRows
}

func (f *fakeStore) CreateSigningKey(ctx context.Context, arg db.CreateSigningKeyParams) (db.SigningKey, error) {
	if f.conflictOnce {
		// Simulate a concurrent provision winning the unique race: install a different key, then report the conflict.
		f.conflictOnce = false
		encryptor, err := keys.NewAESGCM(testMasterKey)
		if err != nil {
			return db.SigningKey{}, err
		}
		winnerKey, err := crypto.GenerateKey()
		if err != nil {
			return db.SigningKey{}, err
		}
		ciphertext, version, err := encryptor.Encrypt(ctx, crypto.FromECDSA(winnerKey))
		if err != nil {
			return db.SigningKey{}, err
		}
		f.keysByIdentity[arg.IdentityID] = db.SigningKey{
			ID:         arg.IdentityID,
			IdentityID: arg.IdentityID,
			Address:    crypto.PubkeyToAddress(winnerKey.PublicKey).Hex(),
			Ciphertext: ciphertext,
			Encryptor:  keys.AESGCMScheme,
			KeyVersion: int32(version),
		}
		return db.SigningKey{}, &pgconn.PgError{Code: uniqueViolationSQLState}
	}
	if _, ok := f.keysByIdentity[arg.IdentityID]; ok {
		return db.SigningKey{}, &pgconn.PgError{Code: uniqueViolationSQLState}
	}
	row := db.SigningKey{
		ID:         arg.IdentityID,
		IdentityID: arg.IdentityID,
		Address:    arg.Address,
		Ciphertext: arg.Ciphertext,
		Encryptor:  arg.Encryptor,
		KeyVersion: arg.KeyVersion,
	}
	f.keysByIdentity[arg.IdentityID] = row
	return row, nil
}
