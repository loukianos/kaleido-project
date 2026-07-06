package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"kaleido-project/internal/auth"
)

type onboardResponse struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject" example:"080dfd67-ccc1-49c1-8a30-c292b3a4a7cf"`
	Address string `json:"address" example:"0x1C7021cc3fAa237D4B837268dc19c12B6003D449"`
}

// handleOnboardLender onboards the authenticated caller as a custodial lender.
//
//	@Summary		Onboard as a lender
//	@Description	Creates the caller's lender identity and eagerly provisions their custodial wallet. Idempotent: an already-onboarded lender gets their existing identity back. A lender must onboard before they can be named as lender_subject or to_subject.
//	@Tags			lenders
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	onboardResponse
//	@Failure		403	{object}	errorResponse	"Platform service accounts do not onboard as lenders"
//	@Router			/lenders/onboard [post]
func handleOnboardLender(logger *slog.Logger, identities IdentityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := principalFrom(c)
		// The platform's on-chain presence is the deployer key with contract roles; its service accounts don't double as custodial lenders.
		if principal.HasRole(auth.RoleServicer) || principal.HasRole(auth.RoleAdmin) {
			c.JSON(http.StatusForbidden, errorBody("platform service accounts do not onboard as lenders"))
			return
		}

		ident, signer, err := identities.OnboardLender(c.Request.Context(), principal.Issuer, principal.Subject)
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "onboard lender failed", "error", err)
			c.JSON(http.StatusInternalServerError, errorBody("onboard lender failed"))
			return
		}

		c.JSON(http.StatusOK, onboardResponse{
			Issuer:  ident.Issuer,
			Subject: ident.Subject,
			Address: signer.Address().Hex(),
		})
	}
}
