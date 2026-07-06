package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
)

type operationResponse struct {
	ID            int64  `json:"id"`
	Kind          string `json:"kind"`
	Status        string `json:"status"`
	Attempts      int32  `json:"attempts"`
	Error         string `json:"error,omitempty"`
	TxHash        string `json:"tx_hash,omitempty"`
	SignerAddress string `json:"signer_address,omitempty"`
	ContractID    int64  `json:"contract_id,omitempty"`
	LoanID        int64  `json:"loan_id,omitempty"`
}

// handleGetOperation reads one chain operation, for servicer-side visibility into retries and failures.
//
//	@Summary		Get chain operation
//	@Description	Operational view of a chain write: status, attempts, last error, and transaction hash. Lenders poll their loan instead; this endpoint is for the servicer.
//	@Tags			operations
//	@Produce		json
//	@Param			id	path	int	true	"Operation ID"
//	@Security		BearerAuth
//	@Success		200	{object}	operationResponse
//	@Failure		404	{object}	errorResponse
//	@Router			/operations/{id} [get]
func handleGetOperation(logger *slog.Logger, service LoansService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := pathID(c)
		if !ok {
			return
		}

		op, err := service.Operation(c.Request.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorBody("operation not found"))
			return
		}
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "get operation failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("get operation failed"))
			return
		}

		c.JSON(http.StatusOK, operationResponseFromStore(op))
	}
}

func operationResponseFromStore(op db.ChainOperation) operationResponse {
	return operationResponse{
		ID:            op.ID,
		Kind:          op.Kind,
		Status:        op.Status,
		Attempts:      op.Attempts,
		Error:         stringFromNull(op.Error),
		TxHash:        stringFromNull(op.TxHash),
		SignerAddress: stringFromNull(op.SignerAddress),
		ContractID:    int64FromNull(op.ContractID),
		LoanID:        int64FromNull(op.LoanID),
	}
}
