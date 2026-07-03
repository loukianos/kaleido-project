package api

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// addRoutes maps the entire API surface in one place.
func addRoutes(router *gin.Engine, version string, logger *slog.Logger, opts Options) {
	startedAt := time.Now().UTC()

	router.GET("/", handleIndex(version))
	router.GET("/healthz", handleHealthz())
	router.GET("/ready", handleReady(logger, startedAt, opts.ReadinessChecks, opts.Contracts))
	router.POST("/admin/contracts/deploy", handleDeployContract(logger, opts.Contracts))
	router.GET("/contracts/active", handleActiveContract(logger, opts.Contracts))
	router.POST("/loans", handleCreateLoan(logger, opts.Loans))
	router.GET("/loans", handleListLoans(logger, opts.Loans))
	router.GET("/loans/:id", handleGetLoan(logger, opts.Loans))
	router.POST("/loans/:id/transfer", handleTransferLoan(logger, opts.Loans))
	router.POST("/loans/:id/default", handleDefaultLoan(logger, opts.Loans))
	router.POST("/loans/:id/repayments", handleCreateRepayment(logger, opts.Loans))
	router.GET("/loans/:id/repayments", handleListRepayments(logger, opts.Loans))
	router.GET("/loans/:id/terms", handleLoanTerms(logger, opts.Loans))
}
