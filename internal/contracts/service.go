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

var ErrContractAlreadyDeployed = errors.New("contract already deployed for this chain")

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

func (s *Service) Deploy(ctx context.Context, baseURI string) (db.Contract, error) {
	if baseURI == "" {
		baseURI = s.baseURI
	}

	if _, err := s.repo.ActiveContract(ctx, s.chain.ChainID().Int64()); err == nil {
		return db.Contract{}, ErrContractAlreadyDeployed
	} else if !errors.Is(err, pgx.ErrNoRows) {
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
		Active:       true,
	})
	if err != nil {
		return db.Contract{}, err
	}

	return contract, nil
}
