package api

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Imported for its side effect: registering the generated OpenAPI
	// definition with the swag runtime, where gin-swagger reads it.
	_ "kaleido-project/docs"
)

// addRoutes maps the entire API surface in one place.
func addRoutes(router *gin.Engine, version string, logger *slog.Logger, opts Options) {
	startedAt := time.Now().UTC()

	router.GET("/", handleIndex(version, opts.SignerAddress))
	router.GET("/healthz", handleHealthz())
	router.GET("/ready", handleReady(logger, startedAt, opts.ReadinessChecks, opts.Contracts))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/admin/contracts/deploy", handleDeployContract(logger, opts.Contracts))
	router.POST("/admin/contracts/:id/activate", handleActivateContract(logger, opts.Contracts))
	router.GET("/contracts", handleListContracts(logger, opts.Contracts))
	router.GET("/contracts/active", handleActiveContract(logger, opts.Contracts))
	router.GET("/contracts/:id", handleGetContract(logger, opts.Contracts))
	router.POST("/loans", handleCreateLoan(logger, opts.Loans))
	router.GET("/loans", handleListLoans(logger, opts.Loans))
	router.GET("/loans/:id", handleGetLoan(logger, opts.Loans))
	router.POST("/loans/:id/transfer", handleTransferLoan(logger, opts.Loans))
	router.POST("/loans/:id/default", handleDefaultLoan(logger, opts.Loans))
	router.POST("/loans/:id/repayments", handleCreateRepayment(logger, opts.Loans))
	router.GET("/loans/:id/repayments", handleListRepayments(logger, opts.Loans))
	router.GET("/loans/:id/terms", handleLoanTerms(logger, opts.Loans))
}
