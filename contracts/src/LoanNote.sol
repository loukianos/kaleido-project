// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {ERC721URIStorage, ERC721} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";

/// Each token is a lender's claim to repayment.
/// The lender named at origination owns the token outright: transfers use the standard ERC-721 owner-signed paths, and no role can move a note it does not hold.
/// The servicer's powers are limited to lifecycle transitions (settle, default), which never reassign the claim.
contract LoanNote is ERC721URIStorage, AccessControl {
    /// May originate loans (mint new notes).
    bytes32 public constant ORIGINATOR_ROLE = keccak256("ORIGINATOR_ROLE");
    /// May settle (burn on final repayment) and mark loans defaulted.
    bytes32 public constant SERVICER_ROLE = keccak256("SERVICER_ROLE");

    enum Status {
        Active,
        Repaid,
        Defaulted
    }

    struct LoanTerms {
        uint256 principal; // cents
        uint16 aprBps; // basis points
        uint64 maturity; // unix timestamp (seconds) the loan matures
        Status status;
    }

    uint256 private _nextId;
    mapping(uint256 tokenId => LoanTerms) private _terms;

    event LoanOriginated(
        uint256 indexed tokenId,
        address indexed lender,
        uint256 principal,
        uint16 aprBps,
        uint64 maturity
    );
    event LoanStatusChanged(uint256 indexed tokenId, Status status);
    event LoanSettled(uint256 indexed tokenId);

    /// The admin can grant and revoke roles.
    /// It starts with both business roles so a single platform key can operate the contract; splitting them across keys is a grant away.
    constructor(address admin) ERC721("LoanNote", "LOAN") {
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ORIGINATOR_ROLE, admin);
        _grantRole(SERVICER_ROLE, admin);
    }

    function originate(address lender, uint256 principal, uint16 aprBps, uint64 maturity, string calldata tokenURI_)
        external
        onlyRole(ORIGINATOR_ROLE)
        returns (uint256 tokenId)
    {
        tokenId = _nextId++;
        _terms[tokenId] = LoanTerms({
            principal: principal,
            aprBps: aprBps,
            maturity: maturity,
            status: Status.Active
        });
        _safeMint(lender, tokenId);
        _setTokenURI(tokenId, tokenURI_);
        emit LoanOriginated(tokenId, lender, principal, aprBps, maturity);
    }

    function markDefaulted(uint256 tokenId) external onlyRole(SERVICER_ROLE) {
        _requireOwned(tokenId);
        _terms[tokenId].status = Status.Defaulted;
        emit LoanStatusChanged(tokenId, Status.Defaulted);
    }

    /// Final repayment extinguishes the claim, so settlement burns the note regardless of who holds it.
    /// This is the one intentional exception to "only the owner touches the token".
    function settle(uint256 tokenId) external onlyRole(SERVICER_ROLE) {
        _requireOwned(tokenId);
        emit LoanStatusChanged(tokenId, Status.Repaid);
        emit LoanSettled(tokenId);
        _burn(tokenId);
        delete _terms[tokenId];
    }

    function terms(uint256 tokenId) external view returns (LoanTerms memory) {
        _requireOwned(tokenId);
        return _terms[tokenId];
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721URIStorage, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
