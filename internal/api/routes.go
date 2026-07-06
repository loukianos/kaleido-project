package api

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"kaleido-project/internal/auth"

	// Imported for its side effect: registering the generated OpenAPI
	// definition with the swag runtime, where gin-swagger reads it.
	_ "kaleido-project/docs"
)

// addRoutes maps the entire API surface in one place.
// System endpoints are public; everything else requires a bearer token, with mutations gated by role and loan reads scoped to the caller.
func addRoutes(router *gin.Engine, version string, logger *slog.Logger, opts Options) {
	startedAt := time.Now().UTC()

	router.GET("/", handleIndex(version, opts.SignerAddress))
	router.GET("/healthz", handleHealthz())
	router.GET("/ready", handleReady(logger, startedAt, opts.ReadinessChecks, opts.Contracts))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authed := router.Group("", requireAuth(logger, opts.Verifier))
	admin := authed.Group("", requireRole(auth.RoleAdmin))
	servicer := authed.Group("", requireRole(auth.RoleServicer))

	admin.POST("/admin/contracts/deploy", handleDeployContract(logger, opts.Contracts))
	admin.POST("/admin/contracts/:id/activate", handleActivateContract(logger, opts.Contracts))
	authed.GET("/contracts", handleListContracts(logger, opts.Contracts))
	authed.GET("/contracts/active", handleActiveContract(logger, opts.Contracts))
	authed.GET("/contracts/:id", handleGetContract(logger, opts.Contracts))

	servicer.POST("/loans", handleCreateLoan(logger, opts.Loans))
	authed.GET("/loans", handleListLoans(logger, opts.Loans, opts.Identities))
	authed.GET("/loans/:id", handleGetLoan(logger, opts.Loans, opts.Identities))
	authed.POST("/loans/:id/transfer", handleTransferLoan(logger, opts.Loans, opts.Identities))
	servicer.POST("/loans/:id/default", handleDefaultLoan(logger, opts.Loans))
	servicer.POST("/loans/:id/repayments", handleCreateRepayment(logger, opts.Loans))
	authed.GET("/loans/:id/repayments", handleListRepayments(logger, opts.Loans, opts.Identities))
	authed.GET("/loans/:id/terms", handleLoanTerms(logger, opts.Loans, opts.Identities))
}
