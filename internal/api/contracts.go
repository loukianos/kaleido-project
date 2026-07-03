package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
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

func handleDeployContract(logger *slog.Logger, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req deployContractRequest
		err := c.ShouldBindJSON(&req)
		// An empty body falls back to the configured base URI, so EOF error is ok.
		if err != nil && !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		contract, err := contracts.Deploy(c.Request.Context(), req.BaseURI)
		if err != nil {
			if errors.Is(err, contractspkg.ErrContractAlreadyDeployed) {
				c.JSON(http.StatusConflict, errorBody(err.Error()))
				return
			}
			if errors.Is(err, db.ErrLockBusy) {
				c.JSON(http.StatusServiceUnavailable, errorBody("another chain operation is in progress, retry shortly"))
				return
			}
			logger.ErrorContext(c.Request.Context(), "deploy contract failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("deploy contract failed"))
			return
		}
		c.JSON(http.StatusCreated, contractResponseFromStore(contract))
	}
}

func handleActiveContract(logger *slog.Logger, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		contract, err := contracts.ActiveContract(c.Request.Context())
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorBody("active contract not found"))
			return
		}
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "get active contract failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("get active contract failed"))
			return
		}
		c.JSON(http.StatusOK, contractResponseFromStore(contract))
	}
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
