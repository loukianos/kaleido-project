package loans

import (
	"errors"
	"math"
	"math/big"
)

var (
	ErrInvalidAmount  = errors.New("amount must be positive")
	ErrInvalidTerm    = errors.New("term days must be positive")
	ErrOverpayment    = errors.New("repayment exceeds outstanding balance")
	ErrAmountOverflow = errors.New("amount exceeds int64 range")
)

const (
	LoanStatusOriginating = "originating"
	LoanStatusActive      = "active"
	LoanStatusSettling    = "settling"
	LoanStatusRepaid      = "repaid"
	LoanStatusDefaulted   = "defaulted"
)

type LoanTerms struct {
	PrincipalMinor   int64
	APRBps           uint16
	TermDays         int64
	InterestDueMinor int64
	TotalDueMinor    int64
}

func NewLoanTerms(principalMinor int64, aprBps uint16, termDays int64) (LoanTerms, error) {
	if principalMinor <= 0 {
		return LoanTerms{}, ErrInvalidAmount
	}
	if termDays <= 0 {
		return LoanTerms{}, ErrInvalidTerm
	}

	interestDue, err := InterestDue(principalMinor, aprBps, termDays)
	if err != nil {
		return LoanTerms{}, err
	}
	if interestDue > math.MaxInt64-principalMinor {
		return LoanTerms{}, ErrAmountOverflow
	}

	return LoanTerms{
		PrincipalMinor:   principalMinor,
		APRBps:           aprBps,
		TermDays:         termDays,
		InterestDueMinor: interestDue,
		TotalDueMinor:    principalMinor + interestDue,
	}, nil
}

func InterestDue(principalMinor int64, aprBps uint16, termDays int64) (int64, error) {
	if principalMinor <= 0 {
		return 0, ErrInvalidAmount
	}
	if termDays <= 0 {
		return 0, ErrInvalidTerm
	}
	if aprBps == 0 {
		return 0, nil
	}

	numerator := big.NewInt(principalMinor)
	numerator.Mul(numerator, big.NewInt(int64(aprBps)))
	numerator.Mul(numerator, big.NewInt(termDays))
	numerator.Div(numerator, big.NewInt(365*10000))

	if !numerator.IsInt64() {
		return 0, ErrAmountOverflow
	}
	return numerator.Int64(), nil
}

func ApplyRepayment(outstandingMinor int64, amountMinor int64) (int64, error) {
	if outstandingMinor < 0 {
		return 0, ErrInvalidAmount
	}
	if amountMinor <= 0 {
		return outstandingMinor, ErrInvalidAmount
	}
	if amountMinor > outstandingMinor {
		return outstandingMinor, ErrOverpayment
	}
	return outstandingMinor - amountMinor, nil
}
