package loans

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestIsUniqueViolation(t *testing.T) {
	duplicateRef := &pgconn.PgError{Code: "23505", ConstraintName: "repayments_external_ref_unique"}

	require.True(t, isUniqueViolation(duplicateRef, "repayments_external_ref_unique"))
	require.True(t, isUniqueViolation(fmt.Errorf("create repayment: %w", duplicateRef), "repayments_external_ref_unique"))

	otherConstraint := &pgconn.PgError{Code: "23505", ConstraintName: "loans_contract_token_unique"}
	require.False(t, isUniqueViolation(otherConstraint, "repayments_external_ref_unique"))

	notNullViolation := &pgconn.PgError{Code: "23502", ConstraintName: "repayments_external_ref_unique"}
	require.False(t, isUniqueViolation(notNullViolation, "repayments_external_ref_unique"))

	require.False(t, isUniqueViolation(errors.New("connection reset"), "repayments_external_ref_unique"))
}
