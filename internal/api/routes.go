package api

import (
	"log/slog"
	"net/http"
	"time"
)

// addRoutes maps the entire API surface in one place.
func addRoutes(mux *http.ServeMux, version string, logger *slog.Logger, opts Options) {
	startedAt := time.Now().UTC()

	mux.Handle("GET /{$}", handleIndex(version))
	mux.Handle("GET /healthz", handleHealthz())
	mux.Handle("GET /ready", handleReady(logger, startedAt, opts.ReadinessChecks, opts.Contracts))
	mux.Handle("POST /admin/contracts/deploy", handleDeployContract(logger, opts.Contracts))
	mux.Handle("GET /contracts/active", handleActiveContract(logger, opts.Contracts))
	mux.Handle("POST /loans", handleCreateLoan(logger, opts.Loans))
	mux.Handle("GET /loans", handleListLoans(logger, opts.Loans))
	mux.Handle("GET /loans/{id}", handleGetLoan(logger, opts.Loans))
	mux.Handle("POST /loans/{id}/transfer", handleTransferLoan(logger, opts.Loans))
	mux.Handle("POST /loans/{id}/default", handleDefaultLoan(logger, opts.Loans))
	mux.Handle("POST /loans/{id}/repayments", handleCreateRepayment(logger, opts.Loans))
	mux.Handle("GET /loans/{id}/repayments", handleListRepayments(logger, opts.Loans))
	mux.Handle("GET /loans/{id}/terms", handleLoanTerms(logger, opts.Loans))
}
