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
	grantRoleOperationKind     = "grant_role"
)

var (
	ErrNoActiveContract           = errors.New("no active contract deployed for this chain")
	ErrInvalidAddress             = errors.New("invalid ethereum address")
	ErrInvalidLender              = errors.New("exactly one of lender_address and lender_subject is required")
	ErrInvalidTransferTarget      = errors.New("exactly one of to_address and to_subject is required")
	ErrLoanNotActive              = errors.New("loan is not active")
	ErrLoanNotTransferable        = errors.New("loan is not transferable")
	ErrNotNoteOwner               = errors.New("transfer requires the note owner's signature")
	// ErrOperationPending reports a journaled chain write the platform will retry; handlers translate it to 202 with the affected resources.
	ErrOperationPending = errors.New("chain operation pending; the platform will retry, poll the loan for progress")
	ErrLoanMissingToken           = errors.New("loan is missing token id")
	ErrLoanMissingContract        = errors.New("loan is missing contract")
	ErrLoanOriginatedEventMissing = errors.New("loan originated event missing from receipt")
	ErrDuplicateExternalRef       = errors.New("repayment with this external ref already recorded")
)

// IdentityResolver maps subjects to custodial identities and signing keys; identity.Service satisfies it.
// LenderAddress is lookup-only: subjects named in request bodies must already be onboarded, so the loans domain can never fabricate identities.
type IdentityResolver interface {
	LenderAddress(ctx context.Context, issuer, subject string) (db.Identity, common.Address, error)
	SignerForIdentity(ctx context.Context, identityID int64) (*eth.Signer, error)
	Identity(ctx context.Context, id int64) (db.Identity, error)
}

// Caller is the authenticated principal a handler resolved, as the loans domain sees it.
type Caller struct {
	// IdentityID is the caller's lender identity; zero when the caller is the servicer.
	IdentityID int64
	Servicer   bool
}

type Service struct {
	repo       *Repository
	chain      eth.Writer
	submitter  *chainop.Submitter
	identities IdentityResolver
	// issuer qualifies subjects named in request bodies; authenticated callers carry their own issuer, and the API accepts a single one.
	issuer string
	// platformPool are the signers for servicer operations: the primary key plus the role-granted pool, so concurrent writes don't serialize on one nonce sequence.
	platformPool []*eth.Signer
}

type OriginateRequest struct {
	BorrowerRef string
	// LenderAddress names an external lender wallet; LenderSubject names an onboarded custodial identity.
	// Exactly one must be set.
	LenderAddress  string
	LenderSubject  string
	PrincipalMinor int64
	APRBps         uint16
	TermDays       int64
	// ContractID selects the loan series to originate into; nil means the chain's active contract.
	ContractID *int64
	// ExternalRef is the client's idempotency key: a retried request with the same ref returns the existing loan instead of minting a sibling.
	ExternalRef string
}

type OriginateResult struct {
	Loan          db.Loan
	LenderSubject string
	OperationID   int64
	TxHash        string
	// Existing reports an idempotent replay: the loan already existed under the request's external_ref.
	Existing bool
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
	// MintSignerAddress is the pool key that signed the mint, surfaced for auditability.
	MintSignerAddress string
}

type ListRequest struct {
	Lender string
	// LenderIdentityID scopes the list to one lender's loans; zero means unscoped (servicer view).
	LenderIdentityID int64
	Status           string
	Limit            int32
	Offset           int32
}

func NewService(repo *Repository, chain eth.Writer, locks *db.LockManager, lockHolder string, identities IdentityResolver, issuer string, poolSigners []*eth.Signer) *Service {
	return &Service{
		repo:         repo,
		chain:        chain,
		submitter:    chainop.NewSubmitter(chain, locks, lockHolder, repo),
		identities:   identities,
		issuer:       issuer,
		platformPool: append([]*eth.Signer{chain.DefaultSigner()}, poolSigners...),
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
		if op.SignerAddress != nil {
			result.MintSignerAddress = *op.SignerAddress
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
		Lender:           req.Lender,
		LenderIdentityID: req.LenderIdentityID,
		Status:           req.Status,
		LimitCount:       req.Limit,
		OffsetCount:      req.Offset,
	})
}

// Loan reads the bare loan row, for authorization checks that don't need the chain-decorated view.
func (s *Service) Loan(ctx context.Context, id int64) (db.Loan, error) {
	return s.repo.Loan(ctx, id)
}

// Operation reads one chain operation, for servicer-side visibility into retries and failures.
func (s *Service) Operation(ctx context.Context, id int64) (db.ChainOperation, error) {
	return s.repo.Operation(ctx, id)
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

	if req.ExternalRef != "" {
		existing, err := s.repo.LoanByExternalRef(ctx, req.ExternalRef)
		if err == nil {
			return s.existingOriginationResult(existing), nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return OriginateResult{}, fmt.Errorf("look up external ref: %w", err)
		}
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
		ident, address, err := s.identities.LenderAddress(ctx, s.issuer, req.LenderSubject)
		if err != nil {
			return OriginateResult{}, fmt.Errorf("resolve lender identity: %w", err)
		}
		lender = address
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
		ExternalRef:      nullableString(req.ExternalRef),
	})
	if isUniqueViolation(err, "loans_external_ref_unique") {
		// A concurrent request with the same idempotency key won the insert; replay theirs.
		existing, getErr := s.repo.LoanByExternalRef(ctx, req.ExternalRef)
		if getErr != nil {
			return OriginateResult{}, fmt.Errorf("look up external ref after conflict: %w", getErr)
		}
		return s.existingOriginationResult(existing), nil
	}
	if err != nil {
		return OriginateResult{}, err
	}

	applied, txHash, err := s.driveOrigination(ctx, loan, contract, op.ID)
	if err != nil {
		if errors.Is(err, chainop.ErrPending) {
			return OriginateResult{Loan: loan, LenderSubject: lenderSubject, OperationID: op.ID}, fmt.Errorf("originate: %w", ErrOperationPending)
		}
		if errors.Is(err, chainop.ErrReverted) {
			if failed, failErr := s.repo.FailLoan(ctx, loan.ID); failErr == nil {
				loan = failed
			}
		}
		return OriginateResult{Loan: loan, OperationID: op.ID}, err
	}

	return OriginateResult{Loan: applied, LenderSubject: lenderSubject, OperationID: op.ID, TxHash: txHash}, nil
}

func (s *Service) existingOriginationResult(loan db.Loan) OriginateResult {
	result := OriginateResult{Loan: loan, Existing: true}
	if loan.MintOperationID != nil {
		result.OperationID = *loan.MintOperationID
	}
	return result
}

// driveOrigination performs the chain half of an origination from durable state: mint, parse the token id, apply.
// The reconciler re-drives retryable originations through the same path, so everything it needs lives on the loan and contract rows.
// Maturity is recomputed at mint time, so a re-driven origination matures TermDays from when it actually mints.
func (s *Service) driveOrigination(ctx context.Context, loan db.Loan, contract db.Contract, opID int64) (db.Loan, string, error) {
	metadataURI := fmt.Sprintf("%s%d/terms", contract.BaseUri, loan.ID)
	lender := common.HexToAddress(loan.LenderAddress)

	var note *contractpkg.LoanNote
	txHash, receipt, err := s.submitOperation(ctx, opID, contract.Address, "originate",
		func(auth *bind.TransactOpts, n *contractpkg.LoanNote) (*types.Transaction, error) {
			note = n
			maturity := uint64(time.Now().UTC().Add(time.Duration(loan.TermDays) * 24 * time.Hour).Unix())
			return n.Originate(auth, lender, big.NewInt(loan.PrincipalMinor), uint16(loan.AprBps), maturity, metadataURI)
		})
	if err != nil {
		return loan, "", err
	}

	tokenID, err := parseOriginatedTokenID(note, receipt)
	if err != nil {
		return loan, "", s.submitter.Retryable(ctx, opID, err)
	}

	applied, err := s.repo.ApplyOrigination(ctx, loan.ID, tokenID.String(), opID, contract.ID)
	if err != nil {
		return loan, "", err
	}
	return applied, txHash, nil
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
		if errors.Is(err, chainop.ErrPending) {
			// The repayment is recorded; only the on-chain burn is still in flight.
			return result, fmt.Errorf("settle: %w", ErrOperationPending)
		}
		if errors.Is(err, chainop.ErrReverted) {
			if failed, failErr := s.repo.FailLoan(ctx, txResult.Loan.ID); failErr == nil {
				result.Loan = failed
			}
		}
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

func (s *Service) Transfer(ctx context.Context, loanID int64, req TransferRequest, caller Caller) (TransferResult, error) {
	var (
		to           common.Address
		toIdentityID *int64
		toSubject    string
	)
	switch {
	case req.ToSubject != "" && req.ToAddress != "", req.ToSubject == "" && req.ToAddress == "":
		return TransferResult{}, ErrInvalidTransferTarget
	case req.ToSubject != "":
		ident, address, err := s.identities.LenderAddress(ctx, s.issuer, req.ToSubject)
		if err != nil {
			return TransferResult{}, fmt.Errorf("resolve transfer target identity: %w", err)
		}
		to = address
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
	// The caller must be the note's owner: a lender signs with their custodial key, and the servicer signs with the platform key for warehouse notes it holds.
	// Externally held notes can't be moved by the API at all; their owner transfers on-chain directly.
	owner, err := s.ownerOf(ctx, loan)
	if err != nil {
		return TransferResult{}, err
	}
	signer, err := s.transferSigner(ctx, caller, owner)
	if err != nil {
		return TransferResult{}, err
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
		// Transfers are never re-driven: re-signing a user-initiated custody change later is the wrong default, so transient failures fail terminally and the lender re-requests.
		if errors.Is(err, chainop.ErrPending) {
			err = errors.Unwrap(err)
			return TransferResult{}, s.submitter.Failed(ctx, op.ID, err)
		}
		return TransferResult{}, err
	}

	loan, err = s.repo.ApplyTransfer(ctx, loan.ID, to.Hex(), toIdentityID, op.ID)
	if err != nil {
		return TransferResult{}, err
	}

	return TransferResult{Loan: loan, LenderSubject: toSubject, OperationID: op.ID, TxHash: txHash}, nil
}

// transferSigner returns the caller's signer when the caller actually owns the note, and ErrNotNoteOwner otherwise.
func (s *Service) transferSigner(ctx context.Context, caller Caller, owner common.Address) (*eth.Signer, error) {
	if caller.Servicer {
		signer := s.chain.DefaultSigner()
		if owner != signer.Address() {
			return nil, fmt.Errorf("%w: note is held by %s", ErrNotNoteOwner, owner.Hex())
		}
		return signer, nil
	}

	signer, err := s.identities.SignerForIdentity(ctx, caller.IdentityID)
	if errors.Is(err, identity.ErrNoCustodialKey) {
		return nil, fmt.Errorf("%w: note is held by %s", ErrNotNoteOwner, owner.Hex())
	}
	if err != nil {
		return nil, fmt.Errorf("resolve caller signer: %w", err)
	}
	if owner != signer.Address() {
		return nil, fmt.Errorf("%w: note is held by %s", ErrNotNoteOwner, owner.Hex())
	}
	return signer, nil
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
		if errors.Is(err, chainop.ErrPending) {
			return DefaultResult{Loan: loan, OperationID: op.ID}, fmt.Errorf("default: %w", ErrOperationPending)
		}
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
// Any free pool signer serves: they all hold the business roles, so servicer operations spread across nonce sequences instead of serializing.
func (s *Service) submitOperation(
	ctx context.Context,
	opID int64,
	contractAddress string,
	action string,
	send func(*bind.TransactOpts, *contractpkg.LoanNote) (*types.Transaction, error),
) (string, *types.Receipt, error) {
	return s.submitter.SubmitAsAny(ctx, s.platformPool, opID, action,
		bindNote(contractAddress, send))
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
		bindNote(contractAddress, send))
}

// bindNote adapts a LoanNote-typed send callback to the submitter's backend-typed one.
func bindNote(
	contractAddress string,
	send func(*bind.TransactOpts, *contractpkg.LoanNote) (*types.Transaction, error),
) func(*bind.TransactOpts, eth.ContractBackend) (*types.Transaction, error) {
	return func(auth *bind.TransactOpts, backend eth.ContractBackend) (*types.Transaction, error) {
		note, err := contractpkg.NewLoanNote(common.HexToAddress(contractAddress), backend)
		if err != nil {
			return nil, fmt.Errorf("bind loan note: %w", err)
		}
		return send(auth, note)
	}
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
