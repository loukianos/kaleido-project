package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
	contractspkg "kaleido-project/internal/contracts"
)

type ContractsService interface {
	Deploy(ctx context.Context, baseURI string, activate bool) (db.Contract, error)
	ActiveContract(context.Context) (db.Contract, error)
	ListContracts(context.Context) ([]db.Contract, error)
	Contract(ctx context.Context, id int64) (db.Contract, error)
	Activate(ctx context.Context, id int64) (db.Contract, error)
}

type deployContractRequest struct {
	BaseURI string `json:"base_uri"`
	// Activate makes the new contract the default for new originations, replacing the current default.
	// The chain's first contract always becomes the default regardless of this flag.
	Activate bool `json:"activate"`
}

// handleDeployContract deploys a new LoanNote contract instance.
//
//	@Summary		Deploy a LoanNote contract
//	@Description	Deploys a new LoanNote contract instance; each instance is its own loan series. The chain's first contract becomes the origination default; later deploys only take over the default when activate is true.
//	@Tags			contracts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		deployContractRequest	false	"Optional base URI override and activation flag. Empty body uses the configured default base URI"
//	@Success		201		{object}	contractResponse
//	@Failure		400		{object}	errorResponse	"Invalid JSON body"
//	@Failure		503		{object}	errorResponse	"Another chain operation is in progress, retry shortly"
//	@Security		BearerAuth
//	@Router			/admin/contracts/deploy [post]
func handleDeployContract(logger *slog.Logger, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req deployContractRequest
		err := c.ShouldBindJSON(&req)
		// An empty body falls back to the configured base URI, so EOF error is ok.
		if err != nil && !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, errorBody("invalid json body"))
			return
		}

		contract, err := contracts.Deploy(c.Request.Context(), req.BaseURI, req.Activate)
		if err != nil {
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

// handleActiveContract reads the active contract's metadata.
//
//	@Summary	Active contract
//	@Tags		contracts
//	@Produce	json
//	@Success	200	{object}	contractResponse
//	@Failure	404	{object}	errorResponse
//	@Security	BearerAuth
//	@Router		/contracts/active [get]
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

// handleListContracts lists every contract deployed on this chain.
//
//	@Summary	List contracts
//	@Tags		contracts
//	@Produce	json
//	@Success	200	{object}	contractsListResponse
//	@Security	BearerAuth
//	@Router		/contracts [get]
func handleListContracts(logger *slog.Logger, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := contracts.ListContracts(c.Request.Context())
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "list contracts failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("list contracts failed"))
			return
		}
		body := contractsListResponse{Contracts: make([]contractResponse, 0, len(list))}
		for _, contract := range list {
			body.Contracts = append(body.Contracts, contractResponseFromStore(contract))
		}
		c.JSON(http.StatusOK, body)
	}
}

// handleGetContract reads one contract's metadata.
//
//	@Summary	Get contract
//	@Tags		contracts
//	@Produce	json
//	@Param		id	path		int	true	"Contract ID"
//	@Success	200	{object}	contractResponse
//	@Failure	400	{object}	errorResponse
//	@Failure	404	{object}	errorResponse
//	@Security	BearerAuth
//	@Router		/contracts/{id} [get]
func handleGetContract(logger *slog.Logger, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := contractIDParam(c)
		if !ok {
			return
		}
		contract, err := contracts.Contract(c.Request.Context(), id)
		if errors.Is(err, contractspkg.ErrContractNotFound) {
			c.JSON(http.StatusNotFound, errorBody(err.Error()))
			return
		}
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "get contract failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("get contract failed"))
			return
		}
		c.JSON(http.StatusOK, contractResponseFromStore(contract))
	}
}

// handleActivateContract makes a contract the default for new originations.
//
//	@Summary		Activate contract
//	@Description	Makes the contract the default for new originations, replacing the current default.
//	@Tags			contracts
//	@Produce		json
//	@Param			id	path		int	true	"Contract ID"
//	@Success		200	{object}	contractResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/admin/contracts/{id}/activate [post]
func handleActivateContract(logger *slog.Logger, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := contractIDParam(c)
		if !ok {
			return
		}
		contract, err := contracts.Activate(c.Request.Context(), id)
		if errors.Is(err, contractspkg.ErrContractNotFound) {
			c.JSON(http.StatusNotFound, errorBody(err.Error()))
			return
		}
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "activate contract failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("activate contract failed"))
			return
		}
		c.JSON(http.StatusOK, contractResponseFromStore(contract))
	}
}

func contractIDParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, errorBody("invalid id"))
		return 0, false
	}
	return id, true
}

type contractResponse struct {
	ID           int64  `json:"id"`
	ChainID      int64  `json:"chain_id"`
	Address      string `json:"address"`
	DeployTxHash string `json:"deploy_tx_hash,omitempty"`
	BaseURI      string `json:"base_uri"`
	Active       bool   `json:"active"`
}

type contractsListResponse struct {
	Contracts []contractResponse `json:"contracts"`
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
