package contracts

import (
	"context"
	"fmt"

	"kaleido-project/db/sqlc"
)

type Repository struct {
	db      db.TxBeginner
	queries *db.Queries
}

func NewRepository(queries *db.Queries, conn db.TxBeginner) *Repository {
	return &Repository{queries: queries, db: conn}
}

func (r *Repository) ActiveContract(ctx context.Context, chainID int64) (db.Contract, error) {
	return r.queries.GetActiveContractByChainID(ctx, chainID)
}

func (r *Repository) Contract(ctx context.Context, id int64) (db.Contract, error) {
	return r.queries.GetContractByID(ctx, id)
}

func (r *Repository) ListContracts(ctx context.Context, chainID int64) ([]db.Contract, error) {
	return r.queries.ListContractsByChainID(ctx, chainID)
}

// Activate makes the contract the chain's default for new originations, demoting any current default in the same transaction.
func (r *Repository) Activate(ctx context.Context, id int64, chainID int64) (db.Contract, error) {
	var contract db.Contract
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		if err := q.DeactivateActiveContract(ctx, chainID); err != nil {
			return fmt.Errorf("deactivate current contract: %w", err)
		}
		var err error
		contract, err = q.ActivateContract(ctx, id)
		if err != nil {
			return fmt.Errorf("activate contract: %w", err)
		}
		return nil
	})
	if err != nil {
		return db.Contract{}, err
	}
	return contract, nil
}

func (r *Repository) CreateDeployOperation(ctx context.Context) (db.ChainOperation, error) {
	return r.queries.CreateChainOperation(ctx, db.CreateChainOperationParams{
		Kind:   deployOperationKind,
		Status: "created",
	})
}

func (r *Repository) CreateContractOperation(ctx context.Context, kind string, contractID int64) (db.ChainOperation, error) {
	return r.queries.CreateChainOperation(ctx, db.CreateChainOperationParams{
		Kind:       kind,
		Status:     "created",
		ContractID: db.Ptr(contractID),
	})
}

func (r *Repository) SetOperationApplied(ctx context.Context, id int64) (db.ChainOperation, error) {
	return r.queries.SetOperationApplied(ctx, id)
}

func (r *Repository) ActivateDeployedContract(ctx context.Context, opID int64, params db.CreateContractParams) (db.Contract, error) {
	var contract db.Contract
	err := db.WithTx(ctx, r.db, r.queries, func(q *db.Queries) error {
		// An activating deploy demotes the current default in the same transaction, so the one-active-per-chain index holds throughout.
		if params.Active {
			if err := q.DeactivateActiveContract(ctx, params.ChainID); err != nil {
				return fmt.Errorf("deactivate current contract: %w", err)
			}
		}
		var err error
		contract, err = q.CreateContract(ctx, params)
		if err != nil {
			return fmt.Errorf("record deployed contract: %w", err)
		}

		if _, err := q.AttachOperationContract(ctx, db.AttachOperationContractParams{
			ID:         opID,
			ContractID: db.Ptr(contract.ID),
		}); err != nil {
			return fmt.Errorf("attach deploy operation contract: %w", err)
		}

		if _, err := q.SetOperationApplied(ctx, opID); err != nil {
			return fmt.Errorf("mark deploy applied: %w", err)
		}
		return nil
	})
	if err != nil {
		return db.Contract{}, err
	}
	return contract, nil
}

func (r *Repository) SetOperationSubmitted(ctx context.Context, params db.SetOperationSubmittedParams) (db.ChainOperation, error) {
	return r.queries.SetOperationSubmitted(ctx, params)
}

func (r *Repository) SetOperationMined(ctx context.Context, id int64) (db.ChainOperation, error) {
	return r.queries.SetOperationMined(ctx, id)
}

func (r *Repository) SetOperationRetryable(ctx context.Context, params db.SetOperationRetryableParams) (db.ChainOperation, error) {
	return r.queries.SetOperationRetryable(ctx, params)
}

func (r *Repository) SetOperationFailed(ctx context.Context, params db.SetOperationFailedParams) (db.ChainOperation, error) {
	return r.queries.SetOperationFailed(ctx, params)
}
