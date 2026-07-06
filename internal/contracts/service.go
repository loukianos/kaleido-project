package contracts

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/chainop"
	"kaleido-project/internal/eth"
)

var ErrContractNotFound = errors.New("contract not found for this chain")

const deployOperationKind = "deploy_contract"

type Service struct {
	repo      *Repository
	chain     eth.Writer
	submitter *chainop.Submitter
	baseURI   string
}

func NewService(repo *Repository, chain eth.Writer, locks *db.LockManager, baseURI string, lockHolder string) *Service {
	return &Service{
		repo:      repo,
		chain:     chain,
		submitter: chainop.NewSubmitter(chain, locks, lockHolder, repo),
		baseURI:   baseURI,
	}
}

func (s *Service) ActiveContract(ctx context.Context) (db.Contract, error) {
	return s.repo.ActiveContract(ctx, s.chain.ChainID().Int64())
}

func (s *Service) ListContracts(ctx context.Context) ([]db.Contract, error) {
	return s.repo.ListContracts(ctx, s.chain.ChainID().Int64())
}

// Contract returns the contract by id, treating ids from other chains as not found.
func (s *Service) Contract(ctx context.Context, id int64) (db.Contract, error) {
	contract, err := s.repo.Contract(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Contract{}, ErrContractNotFound
	}
	if err != nil {
		return db.Contract{}, err
	}
	if contract.ChainID != s.chain.ChainID().Int64() {
		return db.Contract{}, ErrContractNotFound
	}
	return contract, nil
}

// Activate makes the contract the chain's default for new originations.
func (s *Service) Activate(ctx context.Context, id int64) (db.Contract, error) {
	contract, err := s.Contract(ctx, id)
	if err != nil {
		return db.Contract{}, err
	}
	return s.repo.Activate(ctx, contract.ID, contract.ChainID)
}

// Deploy deploys a new LoanNote instance; each instance is its own loan series.
// The chain's first contract becomes the origination default automatically; later deploys only take over the default when activate is set.
func (s *Service) Deploy(ctx context.Context, baseURI string, activate bool) (db.Contract, error) {
	if baseURI == "" {
		baseURI = s.baseURI
	}

	if _, err := s.repo.ActiveContract(ctx, s.chain.ChainID().Int64()); errors.Is(err, pgx.ErrNoRows) {
		activate = true
	} else if err != nil {
		return db.Contract{}, fmt.Errorf("check existing contract: %w", err)
	}

	op, err := s.repo.CreateDeployOperation(ctx)
	if err != nil {
		return db.Contract{}, fmt.Errorf("create deploy operation: %w", err)
	}

	var address common.Address
	txHash, _, err := s.submitter.Submit(ctx, op.ID, "deploy",
		func(auth *bind.TransactOpts, backend eth.ContractBackend) (*types.Transaction, error) {
			addr, tx, _, err := DeployLoanNote(auth, backend, s.chain.SignerAddress())
			if err != nil {
				return nil, err
			}
			address = addr
			return tx, nil
		})
	if err != nil {
		return db.Contract{}, err
	}

	contract, err := s.repo.ActivateDeployedContract(ctx, op.ID, db.CreateContractParams{
		ChainID:      s.chain.ChainID().Int64(),
		Address:      address.Hex(),
		DeployTxHash: db.Ptr(txHash),
		BaseUri:      baseURI,
		Active:       activate,
	})
	if err != nil {
		return db.Contract{}, err
	}

	return contract, nil
}
