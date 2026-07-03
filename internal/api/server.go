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
}

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

func handleIndex(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"service": "kaleido-project-api",
			"version": version,
		})
	}
}

func handleHealthz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}

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

		body := map[string]any{
			"status":     status,
			"started_at": startedAt.Format(time.RFC3339),
			"checks":     checks,
		}
		contract, err := contracts.ActiveContract(ctx)
		if err == nil {
			body["active_contract"] = contract.Address
		} else if !errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(ctx, "active contract lookup failed",
				slog.String("error", err.Error()),
			)
		}

		c.JSON(statusCode, body)
	}
}

func errorBody(msg string) map[string]string {
	return map[string]string{"error": msg}
}
