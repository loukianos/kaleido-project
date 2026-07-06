package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrMissingPrivateKey = errors.New("deployer private key is required")
	ErrChainIDMismatch   = errors.New("ethereum chain id mismatch")
)

type RPCClient interface {
	ChainID(context.Context) (*big.Int, error)
	PendingNonceAt(context.Context, common.Address) (uint64, error)
	SendTransaction(context.Context, *types.Transaction) error
}

type ContractBackend interface {
	bind.ContractBackend
	bind.DeployBackend
}

// Writer is the shared chain connection.
// Signing is per-request: callers pass a Signer (the platform's DefaultSigner or a custodial identity's) to the submitter.
type Writer interface {
	SignerAddress() common.Address
	DefaultSigner() *Signer
	ChainID() *big.Int
	PendingNonceOf(context.Context, common.Address) (uint64, error)
	Backend() (ContractBackend, error)
}

type Client struct {
	rpc     RPCClient
	backend ContractBackend
	signer  *Signer
	chainID *big.Int
}

func Dial(ctx context.Context, rpcURL string, chainID int64, privateKeyHex string) (*Client, error) {
	rpc, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum rpc: %w", err)
	}

	client, err := New(rpc, chainID, privateKeyHex)
	if err != nil {
		rpc.Close()
		return nil, err
	}
	if err := client.CheckChainID(ctx); err != nil {
		rpc.Close()
		return nil, err
	}
	return client, nil
}

func New(rpc RPCClient, chainID int64, privateKeyHex string) (*Client, error) {
	if rpc == nil {
		return nil, errors.New("ethereum rpc client is required")
	}
	if chainID <= 0 {
		return nil, fmt.Errorf("chain id must be positive, got %d", chainID)
	}

	signer, err := NewSigner(privateKeyHex)
	if err != nil {
		return nil, err
	}

	client := &Client{
		rpc:     rpc,
		signer:  signer,
		chainID: big.NewInt(chainID),
	}
	if backend, ok := rpc.(ContractBackend); ok {
		client.backend = backend
	}
	return client, nil
}

// SignerAddress is the platform default signer's address.
func (c *Client) SignerAddress() common.Address {
	return c.signer.Address()
}

// DefaultSigner is the platform key, used for servicer-side writes (deploy, originate, settle, default).
func (c *Client) DefaultSigner() *Signer {
	return c.signer
}

func (c *Client) ChainID() *big.Int {
	return new(big.Int).Set(c.chainID)
}

func (c *Client) Backend() (ContractBackend, error) {
	if c.backend == nil {
		return nil, errors.New("ethereum contract backend is unavailable")
	}
	return c.backend, nil
}

func (c *Client) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.CheckChainID(ctx)
}

func (c *Client) CheckChainID(ctx context.Context) error {
	got, err := c.rpc.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get ethereum chain id: %w", err)
	}
	if got.Cmp(c.chainID) != 0 {
		return fmt.Errorf("%w: configured=%s actual=%s", ErrChainIDMismatch, c.chainID, got)
	}
	return nil
}

func (c *Client) PendingNonceOf(ctx context.Context, address common.Address) (uint64, error) {
	nonce, err := c.rpc.PendingNonceAt(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("get pending nonce: %w", err)
	}
	return nonce, nil
}

func (c *Client) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	return c.signer.SignTx(c.chainID, tx)
}

func (c *Client) SendSignedTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := c.rpc.SendTransaction(ctx, tx); err != nil {
		return fmt.Errorf("send raw transaction: %w", err)
	}
	return nil
}

func (c *Client) SignAndSendTransaction(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	signed, err := c.SignTransaction(tx)
	if err != nil {
		return nil, err
	}
	if err := c.SendSignedTransaction(ctx, signed); err != nil {
		return nil, err
	}
	return signed, nil
}

func ParsePrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	hexKey = strings.TrimSpace(hexKey)
	if hexKey == "" {
		return nil, ErrMissingPrivateKey
	}
	if len(hexKey) >= 2 && hexKey[0] == '0' && (hexKey[1] == 'x' || hexKey[1] == 'X') {
		hexKey = hexKey[2:]
	}
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, fmt.Errorf("parse deployer private key: %w", err)
	}
	return key, nil
}

// LockNameFor names the writer lock serializing nonce use for one signing address.
// Locks are per key, so different identities' transactions never contend.
func LockNameFor(chainID *big.Int, address common.Address) string {
	return fmt.Sprintf("ethereum-writer:%s:%s", chainID, address.Hex())
}
