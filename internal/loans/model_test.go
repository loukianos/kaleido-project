package loans

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoanTerms(t *testing.T) {
	tests := map[string]struct {
		principal int64
		aprBps    uint16
		termDays  int64
		want      LoanTerms
		wantErr   error
	}{
		"zero apr": {
			principal: 100_00,
			aprBps:    0,
			termDays:  365,
			want: LoanTerms{
				PrincipalMinor:   100_00,
				APRBps:           0,
				TermDays:         365,
				InterestDueMinor: 0,
				TotalDueMinor:    100_00,
			},
		},
		"one day term floors interest": {
			principal: 100_00,
			aprBps:    1000,
			termDays:  1,
			want: LoanTerms{
				PrincipalMinor:   100_00,
				APRBps:           1000,
				TermDays:         1,
				InterestDueMinor: 2,
				TotalDueMinor:    100_02,
			},
		},
		"normal interest": {
			principal: 10_000_00,
			aprBps:    800,
			termDays:  365,
			want: LoanTerms{
				PrincipalMinor:   10_000_00,
				APRBps:           800,
				TermDays:         365,
				InterestDueMinor: 80_000,
				TotalDueMinor:    10_800_00,
			},
		},
		"reject zero principal": {
			principal: 0,
			aprBps:    800,
			termDays:  365,
			wantErr:   ErrInvalidAmount,
		},
		"reject zero term": {
			principal: 100_00,
			aprBps:    800,
			termDays:  0,
			wantErr:   ErrInvalidTerm,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewLoanTerms(tc.principal, tc.aprBps, tc.termDays)
			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr != nil {
				return
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestNewLoanTermsDetectsOverflow(t *testing.T) {
	_, err := NewLoanTerms(math.MaxInt64, math.MaxUint16, math.MaxInt64)
	require.ErrorIs(t, err, ErrAmountOverflow)
}

func TestApplyRepayment(t *testing.T) {
	tests := map[string]struct {
		outstanding int64
		amount      int64
		want        int64
		wantErr     error
	}{
		"partial": {
			outstanding: 500,
			amount:      125,
			want:        375,
		},
		"exact final": {
			outstanding: 500,
			amount:      500,
			want:        0,
		},
		"reject overpayment": {
			outstanding: 500,
			amount:      501,
			want:        500,
			wantErr:     ErrOverpayment,
		},
		"reject zero amount": {
			outstanding: 500,
			amount:      0,
			want:        500,
			wantErr:     ErrInvalidAmount,
		},
		"reject negative amount": {
			outstanding: 500,
			amount:      -1,
			want:        500,
			wantErr:     ErrInvalidAmount,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ApplyRepayment(tc.outstanding, tc.amount)
			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}
