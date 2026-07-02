package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
	contractspkg "kaleido-project/internal/contracts"
)

type ContractsService interface {
	Deploy(context.Context, string) (db.Contract, error)
	ActiveContract(context.Context) (db.Contract, error)
}

type deployContractRequest struct {
	BaseURI string `json:"base_uri"`
}

func handleDeployContract(logger *slog.Logger, contracts ContractsService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decode[deployContractRequest](r)
		// An empty body falls back to the configured base URI, so EOF error is ok.
		if err != nil && !errors.Is(err, io.EOF) {
			writeJSON(w, http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		contract, err := contracts.Deploy(r.Context(), req.BaseURI)
		if err != nil {
			if errors.Is(err, contractspkg.ErrContractAlreadyDeployed) {
				writeJSON(w, http.StatusConflict, errorBody(err.Error()))
				return
			}
			if errors.Is(err, db.ErrLockBusy) {
				writeJSON(w, http.StatusServiceUnavailable, errorBody("another chain operation is in progress, retry shortly"))
				return
			}
			logger.ErrorContext(r.Context(), "deploy contract failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("deploy contract failed"))
			return
		}
		writeJSON(w, http.StatusCreated, contractResponseFromStore(contract))
	})
}

func handleActiveContract(logger *slog.Logger, contracts ContractsService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contract, err := contracts.ActiveContract(r.Context())
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, errorBody("active contract not found"))
			return
		}
		if err != nil {
			logger.ErrorContext(r.Context(), "get active contract failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, errorBody("get active contract failed"))
			return
		}
		writeJSON(w, http.StatusOK, contractResponseFromStore(contract))
	})
}

type contractResponse struct {
	ID           int64  `json:"id"`
	ChainID      int64  `json:"chain_id"`
	Address      string `json:"address"`
	DeployTxHash string `json:"deploy_tx_hash,omitempty"`
	BaseURI      string `json:"base_uri"`
	Active       bool   `json:"active"`
}

func contractResponseFromStore(contract db.Contract) contractResponse {
	return contractResponse{
		ID:           contract.ID,
		ChainID:      contract.ChainID,
		Address:      contract.Address,
		DeployTxHash: stringFromNull(contract.DeployTxHash),
		BaseURI:      contract.BaseUri,
		Active:       contract.Active,
	}
}
