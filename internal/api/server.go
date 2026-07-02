package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

	mux := http.NewServeMux()
	addRoutes(mux, version, logger, opts)

	var handler http.Handler = mux
	handler = requestLogger(logger, handler)
	return handler
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		logger.LogAttrs(r.Context(), slog.LevelInfo, "http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote", r.RemoteAddr),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func handleIndex(version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"service": "kaleido-project-api",
			"version": version,
		})
	})
}

func handleHealthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
}

func handleReady(logger *slog.Logger, startedAt time.Time, readinessChecks []ReadinessCheck, contracts ContractsService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checks := map[string]string{"api": "ok"}
		ready := true
		for _, check := range readinessChecks {
			if err := check.Check(r.Context()); err != nil {
				ready = false
				checks[check.Name] = "error"
				logger.WarnContext(r.Context(), "readiness check failed",
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
		contract, err := contracts.ActiveContract(r.Context())
		if err == nil {
			body["active_contract"] = contract.Address
		} else if !errors.Is(err, pgx.ErrNoRows) {
			logger.WarnContext(r.Context(), "active contract lookup failed",
				slog.String("error", err.Error()),
			)
		}

		writeJSON(w, statusCode, body)
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func errorBody(msg string) map[string]string {
	return map[string]string{"error": msg}
}
