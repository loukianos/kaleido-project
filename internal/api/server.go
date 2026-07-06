package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ReadinessCheck struct {
	Name  string
	Check func(context.Context) error
}

type Options struct {
	ReadinessChecks []ReadinessCheck
	Contracts       ContractsService
	Loans           LoansService
	// SignerAddress is the platform signing key's address, surfaced so callers can originate notes into platform custody (warehouse loans).
	SignerAddress string
}

// @title			Loan Note API
// @version		0.1.0
// @description	API for managing ERC-721-backed loan notes
func New(version string, logger *slog.Logger, opts Options) http.Handler {
	if opts.Loans == nil || opts.Contracts == nil {
		panic("api: Options.Loans and Options.Contracts are required")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(requestLogger(logger), gin.Recovery())
	addRoutes(router, version, logger, opts)
	return router
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "http request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote", c.Request.RemoteAddr),
		)
	}
}

type serviceInfoResponse struct {
	Service       string `json:"service"`
	Version       string `json:"version"`
	SignerAddress string `json:"signer_address" example:"0x627306090abaB3A6e1400e9345bC60c78a8BEf57"`
}

type healthResponse struct {
	Status string `json:"status" example:"ok"`
}

type readinessResponse struct {
	Status         string            `json:"status" example:"ready"`
	StartedAt      string            `json:"started_at" example:"2026-01-02T15:04:05Z"`
	ActiveContract string            `json:"active_contract,omitempty"`
	Checks         map[string]string `json:"checks"`
}

// handleIndex reports service metadata.
//
//	@Summary	Service metadata
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	serviceInfoResponse
//	@Router		/ [get]
func handleIndex(version, signerAddress string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, serviceInfoResponse{
			Service:       "kaleido-project-api",
			Version:       version,
			SignerAddress: signerAddress,
		})
	}
}

// handleHealthz is the liveness probe.
//
//	@Summary	Liveness probe
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	healthResponse
//	@Router		/healthz [get]
func handleHealthz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, healthResponse{Status: "ok"})
	}
}

// handleReady is the readiness probe; it runs the injected dependency checks.
//
//	@Summary	Readiness probe
//	@Tags		system
//	@Produce	json
//	@Success	200	{object}	readinessResponse
//	@Failure	503	{object}	readinessResponse	"One or more readiness checks failed"
//	@Router		/ready [get]
func handleReady(logger *slog.Logger, startedAt time.Time, readinessChecks []ReadinessCheck, contracts ContractsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		checks := map[string]string{"api": "ok"}
		ready := true
		for _, check := range readinessChecks {
			if err := check.Check(ctx); err != nil {
				ready = false
				checks[check.Name] = "error"
				logger.WarnContext(ctx, "readiness check failed",
					slog.String("check", check.Name),
					slog.String("error", err.Error()),
				)
				continue
			}
			checks[check.Name] = "ok"
		}

		statusCode := http.StatusOK
		status := "ready"
		if !ready {
			statusCode = http.StatusServiceUnavailable
			status = "not_ready"
		}

		body := readinessResponse{
			Status:    status,
			StartedAt: startedAt.Format(time.RFC3339),
			Checks:    checks,
		}
		contract, err := contracts.ActiveContract(ctx)
		if err == nil {
			body.ActiveContract = contract.Address
		} else if !errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "active contract lookup failed",
				slog.String("error", err.Error()),
			)
		}

		c.JSON(statusCode, body)
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func errorBody(msg string) errorResponse {
	return errorResponse{Error: msg}
}
