package contracts

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/chainop"
	"kaleido-project/internal/eth"
)

var ErrContractNotFound = errors.New("contract not found for this chain")

const (
	deployOperationKind    = "deploy_contract"
	grantRoleOperationKind = "grant_role"
)

// Role ids mirror the contract's keccak256 constants so grants don't need a chain read first.
var (
	originatorRole = crypto.Keccak256Hash([]byte("ORIGINATOR_ROLE"))
	servicerRole   = crypto.Keccak256Hash([]byte("SERVICER_ROLE"))
)

type Service struct {
	repo      *Repository
	chain     eth.Writer
	submitter *chainop.Submitter
	baseURI   string
	// poolAddresses are the servicer pool keys that must hold the business roles on every contract instance.
	poolAddresses []common.Address
}

func NewService(repo *Repository, chain eth.Writer, locks *db.LockManager, baseURI string, lockHolder string, poolAddresses []common.Address) *Service {
	return &Service{
		repo:          repo,
		chain:         chain,
		submitter:     chainop.NewSubmitter(chain, locks, lockHolder, repo),
		baseURI:       baseURI,
		poolAddresses: poolAddresses,
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

	// The constructor grants roles only to the admin key; pool keys get theirs here.
	// The contract is recorded either way: grants are idempotent, and the startup reconcile heals any that failed.
	if err := s.EnsureRoles(ctx, contract); err != nil {
		return contract, fmt.Errorf("contract %d deployed but pool role grants failed: %w", contract.ID, err)
	}

	return contract, nil
}

// EnsureRoles grants the business roles to every pool key that is missing them on the contract, signed by the admin key.
// Idempotent: keys that already hold a role are skipped, so it doubles as the startup reconcile for pre-existing contracts.
func (s *Service) EnsureRoles(ctx context.Context, contract db.Contract) error {
	if len(s.poolAddresses) == 0 {
		return nil
	}
	backend, err := s.chain.Backend()
	if err != nil {
		return fmt.Errorf("get backend: %w", err)
	}
	note, err := NewLoanNote(common.HexToAddress(contract.Address), backend)
	if err != nil {
		return fmt.Errorf("bind loan note: %w", err)
	}

	for _, address := range s.poolAddresses {
		for _, role := range [][32]byte{originatorRole, servicerRole} {
			has, err := note.HasRole(&bind.CallOpts{Context: ctx}, role, address)
			if err != nil {
				return fmt.Errorf("check role on contract %d: %w", contract.ID, err)
			}
			if has {
				continue
			}

			op, err := s.repo.CreateContractOperation(ctx, grantRoleOperationKind, contract.ID)
			if err != nil {
				return fmt.Errorf("create grant operation: %w", err)
			}
			_, _, err = s.submitter.Submit(ctx, op.ID, "grant role",
				func(auth *bind.TransactOpts, backend eth.ContractBackend) (*types.Transaction, error) {
					note, err := NewLoanNote(common.HexToAddress(contract.Address), backend)
					if err != nil {
						return nil, fmt.Errorf("bind loan note: %w", err)
					}
					return note.GrantRole(auth, role, address)
				})
			if err != nil {
				return fmt.Errorf("grant role to %s on contract %d: %w", address.Hex(), contract.ID, err)
			}
			if _, err := s.repo.SetOperationApplied(ctx, op.ID); err != nil {
				return fmt.Errorf("mark grant applied: %w", err)
			}
		}
	}
	return nil
}

// EnsureAllRoles reconciles pool role grants across every contract on this chain.
func (s *Service) EnsureAllRoles(ctx context.Context) error {
	contracts, err := s.ListContracts(ctx)
	if err != nil {
		return fmt.Errorf("list contracts: %w", err)
	}
	for _, contract := range contracts {
		if err := s.EnsureRoles(ctx, contract); err != nil {
			return err
		}
	}
	return nil
}
