package eth

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const devPrivateKey = "0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"

func TestNewDerivesSignerAddress(t *testing.T) {
	client, err := New(&fakeRPC{chainID: big.NewInt(1337)}, 1337, devPrivateKey)
	require.NoError(t, err)

	want := common.HexToAddress("0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73")
	require.Equal(t, want, client.SignerAddress())
}

func TestNewRequiresPrivateKey(t *testing.T) {
	_, err := New(&fakeRPC{chainID: big.NewInt(1337)}, 1337, "")
	require.ErrorIs(t, err, ErrMissingPrivateKey)
}

func TestCheckChainID(t *testing.T) {
	client, err := New(&fakeRPC{chainID: big.NewInt(1337)}, 1337, devPrivateKey)
	require.NoError(t, err)

	require.NoError(t, client.CheckChainID(context.Background()))
}

func TestCheckChainIDMismatch(t *testing.T) {
	client, err := New(&fakeRPC{chainID: big.NewInt(1)}, 1337, devPrivateKey)
	require.NoError(t, err)

	err = client.CheckChainID(context.Background())
	require.ErrorIs(t, err, ErrChainIDMismatch)
}

func TestSignTransactionRecoversSigner(t *testing.T) {
	client, err := New(&fakeRPC{chainID: big.NewInt(1337)}, 1337, devPrivateKey)
	require.NoError(t, err)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     7,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1),
		Gas:       21_000,
		To:        &common.Address{},
		Value:     big.NewInt(0),
	})

	signed, err := client.SignTransaction(tx)
	require.NoError(t, err)

	signer := types.LatestSignerForChainID(big.NewInt(1337))
	sender, err := types.Sender(signer, signed)
	require.NoError(t, err)
	require.Equal(t, client.SignerAddress(), sender)
}

func TestPendingNonce(t *testing.T) {
	rpc := &fakeRPC{chainID: big.NewInt(1337), nonce: 42}
	client, err := New(rpc, 1337, devPrivateKey)
	require.NoError(t, err)

	nonce, err := client.PendingNonce(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(42), nonce)
	require.Equal(t, client.SignerAddress(), rpc.nonceAddress)
}

func TestSignAndSendTransaction(t *testing.T) {
	rpc := &fakeRPC{chainID: big.NewInt(1337)}
	client, err := New(rpc, 1337, devPrivateKey)
	require.NoError(t, err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		GasPrice: big.NewInt(1),
		Gas:      21_000,
		To:       &common.Address{},
		Value:    big.NewInt(0),
	})

	signed, err := client.SignAndSendTransaction(context.Background(), tx)
	require.NoError(t, err)
	require.NotNil(t, rpc.sentTx)
	require.Equal(t, signed.Hash(), rpc.sentTx.Hash())
}

type fakeRPC struct {
	chainID      *big.Int
	chainIDErr   error
	nonce        uint64
	nonceAddress common.Address
	nonceErr     error
	sentTx       *types.Transaction
	sendErr      error
}

func (f *fakeRPC) ChainID(context.Context) (*big.Int, error) {
	if f.chainIDErr != nil {
		return nil, f.chainIDErr
	}
	return new(big.Int).Set(f.chainID), nil
}

func (f *fakeRPC) PendingNonceAt(_ context.Context, address common.Address) (uint64, error) {
	f.nonceAddress = address
	return f.nonce, f.nonceErr
}

func (f *fakeRPC) SendTransaction(_ context.Context, tx *types.Transaction) error {
	f.sentTx = tx
	return f.sendErr
}
