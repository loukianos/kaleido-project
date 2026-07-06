package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/auth"
	"kaleido-project/internal/eth"
	"kaleido-project/internal/loans"
)

// IdentityService resolves authenticated principals to lender identities; identity.Service satisfies it.
type IdentityService interface {
	ResolveIdentity(ctx context.Context, issuer, subject string) (db.Identity, error)
	OnboardLender(ctx context.Context, issuer, subject string) (db.Identity, *eth.Signer, error)
}

const principalContextKey = "auth.principal"

// requireAuth validates the bearer token and stashes the principal for handlers.
func requireAuth(logger *slog.Logger, verifier auth.Verifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorBody("bearer token required"))
			return
		}
		principal, err := verifier.Verify(c.Request.Context(), token)
		if err != nil {
			logger.WarnContext(c.Request.Context(), "token verification failed", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorBody("invalid bearer token"))
			return
		}
		c.Set(principalContextKey, principal)
		c.Next()
	}
}

// requireRole rejects principals without the given realm role.
func requireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !principalFrom(c).HasRole(role) {
			c.AbortWithStatusJSON(http.StatusForbidden, errorBody(role+" role required"))
			return
		}
		c.Next()
	}
}

func principalFrom(c *gin.Context) auth.Principal {
	principal, _ := c.MustGet(principalContextKey).(auth.Principal)
	return principal
}

// callerFor resolves the request principal into the loans domain's Caller, creating lender identities on first authenticated sight.
// It writes the error response and returns ok=false on failure.
func callerFor(c *gin.Context, logger *slog.Logger, identities IdentityService) (loans.Caller, bool) {
	principal := principalFrom(c)
	if principal.HasRole(auth.RoleServicer) {
		return loans.Caller{Servicer: true}, true
	}
	ident, err := identities.ResolveIdentity(c.Request.Context(), principal.Issuer, principal.Subject)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "resolve caller identity failed", "error", err)
		c.JSON(http.StatusInternalServerError, errorBody("resolve caller identity failed"))
		return loans.Caller{}, false
	}
	return loans.Caller{IdentityID: ident.ID}, true
}

// canReadLoan scopes reads: the servicer sees everything, a lender sees only loans they currently hold.
func canReadLoan(caller loans.Caller, loan db.Loan) bool {
	if caller.Servicer {
		return true
	}
	return loan.LenderIdentityID != nil && *loan.LenderIdentityID == caller.IdentityID
}
