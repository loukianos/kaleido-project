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

type Writer interface {
	SignerAddress() common.Address
	ChainID() *big.Int
	PendingNonce(context.Context) (uint64, error)
	TransactOpts(context.Context, uint64) (*bind.TransactOpts, error)
	Backend() (ContractBackend, error)
}

type Client struct {
	rpc           RPCClient
	backend       ContractBackend
	privateKey    *ecdsa.PrivateKey
	signerAddress common.Address
	chainID       *big.Int
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

	privateKey, err := ParsePrivateKey(privateKeyHex)
	if err != nil {
		return nil, err
	}

	client := &Client{
		rpc:           rpc,
		privateKey:    privateKey,
		signerAddress: crypto.PubkeyToAddress(privateKey.PublicKey),
		chainID:       big.NewInt(chainID),
	}
	if backend, ok := rpc.(ContractBackend); ok {
		client.backend = backend
	}
	return client, nil
}

func (c *Client) SignerAddress() common.Address {
	return c.signerAddress
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

func (c *Client) TransactOpts(ctx context.Context, nonce uint64) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("create transaction signer: %w", err)
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.GasPrice = big.NewInt(0)
	return opts, nil
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

func (c *Client) PendingNonce(ctx context.Context) (uint64, error) {
	nonce, err := c.rpc.PendingNonceAt(ctx, c.signerAddress)
	if err != nil {
		return 0, fmt.Errorf("get pending nonce: %w", err)
	}
	return nonce, nil
}

func (c *Client) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(c.chainID)
	signed, err := types.SignTx(tx, signer, c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}
	return signed, nil
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

func LockName(w Writer) string {
	return fmt.Sprintf("ethereum-writer:%s:%s", w.ChainID(), w.SignerAddress().Hex())
}
