package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signer holds one secp256k1 key and signs transactions as its address.
// The platform key and every custodial identity key are Signers; chain writes pick the Signer that owns the action.
type Signer struct {
	key     *ecdsa.PrivateKey
	address common.Address
}

func NewSigner(privateKeyHex string) (*Signer, error) {
	key, err := ParsePrivateKey(privateKeyHex)
	if err != nil {
		return nil, err
	}
	return NewSignerFromKey(key), nil
}

func NewSignerFromKey(key *ecdsa.PrivateKey) *Signer {
	return &Signer{key: key, address: crypto.PubkeyToAddress(key.PublicKey)}
}

func (s *Signer) Address() common.Address {
	return s.address
}

// TransactOpts binds the signer to a chain and nonce for a single contract write.
func (s *Signer) TransactOpts(ctx context.Context, chainID *big.Int, nonce uint64) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(s.key, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transaction signer: %w", err)
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.GasPrice = big.NewInt(0)
	return opts, nil
}

func (s *Signer) SignTx(chainID *big.Int, tx *types.Transaction) (*types.Transaction, error) {
	signed, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), s.key)
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}
	return signed, nil
}
