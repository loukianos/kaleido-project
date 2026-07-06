package loans

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"

	db "kaleido-project/db/sqlc"
	"kaleido-project/internal/chainop"
	contractpkg "kaleido-project/internal/contracts"
	"kaleido-project/internal/eth"
	"kaleido-project/internal/identity"
)

const (
	originateOperationKind     = "originate"
	settleOperationKind        = "settle"
	transferOperationKind      = "transfer"
	markDefaultedOperationKind = "mark_defaulted"
)

var (
	ErrNoActiveContract           = errors.New("no active contract deployed for this chain")
	ErrInvalidAddress             = errors.New("invalid ethereum address")
	ErrInvalidLender              = errors.New("exactly one of lender_address and lender_subject is required")
	ErrInvalidTransferTarget      = errors.New("exactly one of to_address and to_subject is required")
	ErrLoanNotActive              = errors.New("loan is not active")
	ErrLoanNotTransferable        = errors.New("loan is not transferable")
	ErrNotNoteOwner               = errors.New("transfer requires the note owner's signature")
	ErrLoanMissingToken           = errors.New("loan is missing token id")
	ErrLoanMissingContract        = errors.New("loan is missing contract")
	ErrLoanOriginatedEventMissing = errors.New("loan originated event missing from receipt")
	ErrDuplicateExternalRef       = errors.New("repayment with this external ref already recorded")
)

// IdentityResolver maps subjects to custodial identities and signing keys; identity.Service satisfies it.
type IdentityResolver interface {
	ResolveLender(ctx context.Context, subject string) (db.Identity, *eth.Signer, error)
	SignerForAddress(ctx context.Context, address common.Address) (*eth.Signer, error)
	Identity(ctx context.Context, id int64) (db.Identity, error)
}

type Service struct {
	repo       *Repository
	chain      eth.Writer
	submitter  *chainop.Submitter
	identities IdentityResolver
}

type OriginateRequest struct {
	BorrowerRef string
	// LenderAddress names an external lender wallet; LenderSubject names a custodial identity whose key is provisioned on demand.
	// Exactly one must be set.
	LenderAddress  string
	LenderSubject  string
	PrincipalMinor int64
	APRBps         uint16
	TermDays       int64
	// ContractID selects the loan series to originate into; nil means the chain's active contract.
	ContractID *int64
}

type OriginateResult struct {
	Loan          db.Loan
	LenderSubject string
	OperationID   int64
	TxHash        string
}

type RepaymentRequest struct {
	AmountMinor int64
	ExternalRef string
}

type RepaymentResult struct {
	Loan                  db.Loan
	Repayment             db.Repayment
	SettlementOperationID int64
	SettlementTxHash      string
}

type TransferRequest struct {
	// ToAddress names an external wallet; ToSubject names a custodial identity whose key is provisioned on demand.
	// Exactly one must be set.
	ToAddress string
	ToSubject string
}

type TransferResult struct {
	Loan          db.Loan
	LenderSubject string
	OperationID   int64
	TxHash        string
}

type DefaultResult struct {
	Loan        db.Loan
	OperationID int64
	TxHash      string
}

type ReadResult struct {
	Loan          db.Loan
	LenderSubject string
	OwnerAddress  string
	MintTxHash    string
}

type ListRequest struct {
	Lender string
	Status string
	Limit  int32
	Offset int32
}

func NewService(repo *Repository, chain eth.Writer, locks *db.LockManager, lockHolder string, identities IdentityResolver) *Service {
	return &Service{
		repo:       repo,
		chain:      chain,
		submitter:  chainop.NewSubmitter(chain, locks, lockHolder, repo),
		identities: identities,
	}
}

func (s *Service) Get(ctx context.Context, id int64) (ReadResult, error) {
	loan, err := s.repo.Loan(ctx, id)
	if err != nil {
		return ReadResult{}, err
	}
	result := ReadResult{Loan: loan}

	if loan.MintOperationID != nil {
		op, err := s.repo.Operation(ctx, *loan.MintOperationID)
		if err != nil {
			return ReadResult{}, fmt.Errorf("get mint operation: %w", err)
		}
		if op.TxHash != nil {
			result.MintTxHash = *op.TxHash
		}
	}

	// The token is burned on settlement, so ownerOf would revert for repaid loans.
	if loan.TokenID != nil && loan.Status != LoanStatusRepaid {
		owner, err := s.ownerOf(ctx, loan)
		if err != nil {
			return ReadResult{}, err
		}
		result.OwnerAddress = owner.Hex()
	}

	if loan.LenderIdentityID != nil {
		ident, err := s.identities.Identity(ctx, *loan.LenderIdentityID)
		if err != nil {
			return ReadResult{}, fmt.Errorf("get lender identity: %w", err)
		}
		result.LenderSubject = ident.Subject
	}
	return result, nil
}

func (s *Service) List(ctx context.Context, req ListRequest) ([]db.Loan, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	return s.repo.ListLoans(ctx, db.ListLoansParams{
		Lender:      req.Lender,
		Status:      req.Status,
		LimitCount:  req.Limit,
		OffsetCount: req.Offset,
	})
}

func (s *Service) Terms(ctx context.Context, id int64) (LoanTerms, error) {
	loan, err := s.repo.Loan(ctx, id)
	if err != nil {
		return LoanTerms{}, err
	}
	return LoanTerms{
		PrincipalMinor:   loan.PrincipalMinor,
		APRBps:           uint16(loan.AprBps),
		TermDays:         loan.TermDays,
		InterestDueMinor: loan.InterestDueMinor,
		TotalDueMinor:    loan.TotalDueMinor,
	}, nil
}

func (s *Service) Originate(ctx context.Context, req OriginateRequest) (OriginateResult, error) {
	terms, err := NewLoanTerms(req.PrincipalMinor, req.APRBps, req.TermDays)
	if err != nil {
		return OriginateResult{}, err
	}

	var (
		lender           common.Address
		lenderIdentityID *int64
		lenderSubject    string
	)
	switch {
	case req.LenderSubject != "" && req.LenderAddress != "", req.LenderSubject == "" && req.LenderAddress == "":
		return OriginateResult{}, ErrInvalidLender
	case req.LenderSubject != "":
		ident, signer, err := s.identities.ResolveLender(ctx, req.LenderSubject)
		if err != nil {
			return OriginateResult{}, fmt.Errorf("resolve lender identity: %w", err)
		}
		lender = signer.Address()
		lenderIdentityID = db.Ptr(ident.ID)
		lenderSubject = ident.Subject
	default:
		lender, err = parseAddress(req.LenderAddress)
		if err != nil {
			return OriginateResult{}, err
		}
	}

	var contract db.Contract
	if req.ContractID != nil {
		contract, err = s.repo.Contract(ctx, *req.ContractID)
		if errors.Is(err, pgx.ErrNoRows) {
			return OriginateResult{}, contractpkg.ErrContractNotFound
		}
		if err != nil {
			return OriginateResult{}, fmt.Errorf("get contract: %w", err)
		}
		if contract.ChainID != s.chain.ChainID().Int64() {
			return OriginateResult{}, contractpkg.ErrContractNotFound
		}
	} else {
		contract, err = s.repo.ActiveContract(ctx, s.chain.ChainID().Int64())
		if errors.Is(err, pgx.ErrNoRows) {
			return OriginateResult{}, ErrNoActiveContract
		}
		if err != nil {
			return OriginateResult{}, fmt.Errorf("get active contract: %w", err)
		}
	}

	loan, op, err := s.repo.CreateOrigination(ctx, CreateOriginationParams{
		ContractID:       contract.ID,
		BorrowerRef:      req.BorrowerRef,
		LenderAddress:    lender.Hex(),
		LenderIdentityID: lenderIdentityID,
		PrincipalMinor:   terms.PrincipalMinor,
		APRBps:           int32(terms.APRBps),
		TermDays:         terms.TermDays,
		InterestDueMinor: terms.InterestDueMinor,
		TotalDueMinor:    terms.TotalDueMinor,
	})
	if err != nil {
		return OriginateResult{}, err
	}

	metadataURI := fmt.Sprintf("%s%d/terms", contract.BaseUri, loan.ID)

	var note *contractpkg.LoanNote
	txHash, receipt, err := s.submitOperation(ctx, op.ID, contract.Address, "originate",
		func(auth *bind.TransactOpts, n *contractpkg.LoanNote) (*types.Transaction, error) {
			note = n
			maturity := uint64(time.Now().UTC().Add(time.Duration(terms.TermDays) * 24 * time.Hour).Unix())
			return n.Originate(auth, lender, big.NewInt(terms.PrincipalMinor), terms.APRBps, maturity, metadataURI)
		})
	if err != nil {
		return OriginateResult{}, err
	}

	tokenID, err := parseOriginatedTokenID(note, receipt)
	if err != nil {
		return OriginateResult{}, s.submitter.Retryable(ctx, op.ID, err)
	}

	loan, err = s.repo.ApplyOrigination(ctx, loan.ID, tokenID.String(), op.ID, contract.ID)
	if err != nil {
		return OriginateResult{}, err
	}

	return OriginateResult{Loan: loan, LenderSubject: lenderSubject, OperationID: op.ID, TxHash: txHash}, nil
}

func (s *Service) RecordRepayment(ctx context.Context, loanID int64, req RepaymentRequest) (RepaymentResult, error) {
	if req.AmountMinor <= 0 {
		return RepaymentResult{}, ErrInvalidAmount
	}

	txResult, err := s.repo.RecordRepayment(ctx, loanID, req)
	if err != nil {
		return RepaymentResult{}, err
	}

	result := RepaymentResult{
		Loan:                  txResult.Loan,
		Repayment:             txResult.Repayment,
		SettlementOperationID: txResult.SettlementOperation.ID,
	}
	if txResult.SettlementOperation.ID == 0 {
		return result, nil
	}

	txHash, err := s.settle(ctx, txResult.Loan, txResult.SettlementOperation)
	if err != nil {
		return result, err
	}
	loan, err := s.repo.ApplySettlement(ctx, txResult.Loan.ID, txResult.SettlementOperation.ID)
	if err != nil {
		return result, err
	}
	result.Loan = loan
	result.SettlementTxHash = txHash
	return result, nil
}

func (s *Service) ListRepayments(ctx context.Context, loanID int64) ([]db.Repayment, error) {
	if _, err := s.repo.Loan(ctx, loanID); err != nil {
		return nil, err
	}
	return s.repo.Repayments(ctx, loanID)
}

func (s *Service) Transfer(ctx context.Context, loanID int64, req TransferRequest) (TransferResult, error) {
	var (
		to           common.Address
		toIdentityID *int64
		toSubject    string
	)
	switch {
	case req.ToSubject != "" && req.ToAddress != "", req.ToSubject == "" && req.ToAddress == "":
		return TransferResult{}, ErrInvalidTransferTarget
	case req.ToSubject != "":
		ident, signer, err := s.identities.ResolveLender(ctx, req.ToSubject)
		if err != nil {
			return TransferResult{}, fmt.Errorf("resolve transfer target identity: %w", err)
		}
		to = signer.Address()
		toIdentityID = db.Ptr(ident.ID)
		toSubject = ident.Subject
	default:
		var err error
		to, err = parseAddress(req.ToAddress)
		if err != nil {
			return TransferResult{}, err
		}
	}

	loan, err := s.repo.Loan(ctx, loanID)
	if err != nil {
		return TransferResult{}, err
	}
	if loan.Status != LoanStatusActive && loan.Status != LoanStatusDefaulted {
		return TransferResult{}, ErrLoanNotTransferable
	}
	if loan.TokenID == nil {
		return TransferResult{}, ErrLoanMissingToken
	}
	contract, err := s.contractForLoan(ctx, loan)
	if err != nil {
		return TransferResult{}, err
	}

	// Transfers are owner-signed ERC-721 transferFrom: the contract has no admin path to move a note.
	// The signer is whoever holds the note: a custodial lender's key when we custody it, the platform key for warehouse notes.
	// Externally held notes can't be moved by the API at all; their owner transfers on-chain directly.
	owner, err := s.ownerOf(ctx, loan)
	if err != nil {
		return TransferResult{}, err
	}
	signer := s.chain.DefaultSigner()
	if owner != signer.Address() {
		signer, err = s.identities.SignerForAddress(ctx, owner)
		if errors.Is(err, identity.ErrNoCustodialKey) {
			return TransferResult{}, fmt.Errorf("%w: note is held by %s", ErrNotNoteOwner, owner.Hex())
		}
		if err != nil {
			return TransferResult{}, fmt.Errorf("resolve owner signer: %w", err)
		}
	}

	op, err := s.repo.CreateLoanOperation(ctx, transferOperationKind, contract.ID, loan.ID)
	if err != nil {
		return TransferResult{}, fmt.Errorf("create transfer operation: %w", err)
	}

	tokenID, err := parseTokenID(loan)
	if err != nil {
		return TransferResult{}, s.submitter.Retryable(ctx, op.ID, err)
	}
	txHash, _, err := s.submitOperationAs(ctx, signer, op.ID, contract.Address, "transfer",
		func(auth *bind.TransactOpts, note *contractpkg.LoanNote) (*types.Transaction, error) {
			return note.SafeTransferFrom(auth, owner, to, tokenID)
		})
	if err != nil {
		return TransferResult{}, err
	}

	loan, err = s.repo.ApplyTransfer(ctx, loan.ID, to.Hex(), toIdentityID, op.ID)
	if err != nil {
		return TransferResult{}, err
	}

	return TransferResult{Loan: loan, LenderSubject: toSubject, OperationID: op.ID, TxHash: txHash}, nil
}

func (s *Service) Default(ctx context.Context, loanID int64) (DefaultResult, error) {
	loan, err := s.repo.Loan(ctx, loanID)
	if err != nil {
		return DefaultResult{}, err
	}
	if loan.Status != LoanStatusActive {
		return DefaultResult{}, ErrLoanNotActive
	}
	if loan.TokenID == nil {
		return DefaultResult{}, ErrLoanMissingToken
	}
	contract, err := s.contractForLoan(ctx, loan)
	if err != nil {
		return DefaultResult{}, err
	}

	op, err := s.repo.CreateLoanOperation(ctx, markDefaultedOperationKind, contract.ID, loan.ID)
	if err != nil {
		return DefaultResult{}, fmt.Errorf("create mark defaulted operation: %w", err)
	}

	tokenID, err := parseTokenID(loan)
	if err != nil {
		return DefaultResult{}, s.submitter.Retryable(ctx, op.ID, err)
	}
	txHash, _, err := s.submitOperation(ctx, op.ID, contract.Address, "default",
		func(auth *bind.TransactOpts, note *contractpkg.LoanNote) (*types.Transaction, error) {
			return note.MarkDefaulted(auth, tokenID)
		})
	if err != nil {
		return DefaultResult{}, err
	}

	loan, err = s.repo.ApplyDefault(ctx, loan.ID, op.ID)
	if err != nil {
		return DefaultResult{}, err
	}

	return DefaultResult{Loan: loan, OperationID: op.ID, TxHash: txHash}, nil
}

func (s *Service) settle(ctx context.Context, loan db.Loan, op db.ChainOperation) (string, error) {
	tokenID, err := parseTokenID(loan)
	if err != nil {
		return "", s.submitter.Retryable(ctx, op.ID, err)
	}
	if op.ContractID == nil {
		return "", s.submitter.Retryable(ctx, op.ID, ErrLoanMissingContract)
	}
	contract, err := s.repo.Contract(ctx, *op.ContractID)
	if err != nil {
		return "", s.submitter.Retryable(ctx, op.ID, fmt.Errorf("get settle contract: %w", err))
	}

	txHash, _, err := s.submitOperation(ctx, op.ID, contract.Address, "settle",
		func(auth *bind.TransactOpts, note *contractpkg.LoanNote) (*types.Transaction, error) {
			return note.Settle(auth, tokenID)
		})
	return txHash, err
}

// submitOperation performs a single platform-signed write against the loan-note contract at contractAddress.
func (s *Service) submitOperation(
	ctx context.Context,
	opID int64,
	contractAddress string,
	action string,
	send func(*bind.TransactOpts, *contractpkg.LoanNote) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	return s.submitOperationAs(ctx, s.chain.DefaultSigner(), opID, contractAddress, action, send)
}

// submitOperationAs performs a single write signed by signer, binding the contract and delegating the tracked submission to the shared chainop submitter.
func (s *Service) submitOperationAs(
	ctx context.Context,
	signer *eth.Signer,
	opID int64,
	contractAddress string,
	action string,
	send func(*bind.TransactOpts, *contractpkg.LoanNote) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	return s.submitter.SubmitAs(ctx, signer, opID, action,
		func(auth *bind.TransactOpts, backend eth.ContractBackend) (*types.Transaction, error) {
			note, err := contractpkg.NewLoanNote(common.HexToAddress(contractAddress), backend)
			if err != nil {
				return nil, fmt.Errorf("bind loan note: %w", err)
			}
			return send(auth, note)
		})
}

func (s *Service) ownerOf(ctx context.Context, loan db.Loan) (common.Address, error) {
	tokenID, err := parseTokenID(loan)
	if err != nil {
		return common.Address{}, err
	}
	contract, err := s.contractForLoan(ctx, loan)
	if err != nil {
		return common.Address{}, err
	}
	backend, err := s.chain.Backend()
	if err != nil {
		return common.Address{}, err
	}
	note, err := contractpkg.NewLoanNote(common.HexToAddress(contract.Address), backend)
	if err != nil {
		return common.Address{}, fmt.Errorf("bind loan note: %w", err)
	}
	return note.OwnerOf(&bind.CallOpts{Context: ctx}, tokenID)
}

func (s *Service) contractForLoan(ctx context.Context, loan db.Loan) (db.Contract, error) {
	if loan.ContractID == nil {
		return db.Contract{}, ErrLoanMissingContract
	}
	return s.repo.Contract(ctx, *loan.ContractID)
}

func parseAddress(hexAddr string) (common.Address, error) {
	if !common.IsHexAddress(hexAddr) {
		return common.Address{}, ErrInvalidAddress
	}
	addr := common.HexToAddress(hexAddr)
	return addr, nil
}

func parseTokenID(loan db.Loan) (*big.Int, error) {
	var raw string
	if loan.TokenID != nil {
		raw = *loan.TokenID
	}
	tokenID, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return nil, fmt.Errorf("invalid token id %q", raw)
	}
	return tokenID, nil
}

func parseOriginatedTokenID(note *contractpkg.LoanNote, receipt *types.Receipt) (*big.Int, error) {
	for _, log := range receipt.Logs {
		if event, err := note.ParseLoanOriginated(*log); err == nil {
			return event.TokenId, nil
		}
	}
	return nil, ErrLoanOriginatedEventMissing
}
