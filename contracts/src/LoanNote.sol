// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {ERC721URIStorage, ERC721} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract LoanNote is ERC721URIStorage, Ownable {
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

    constructor(address initialOwner) ERC721("LoanNote", "LOAN") Ownable(initialOwner) {}

    function originate(address lender, uint256 principal, uint16 aprBps, uint64 maturity, string calldata tokenURI_)
        external
        onlyOwner
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

    function markDefaulted(uint256 tokenId) external onlyOwner {
        _requireOwned(tokenId);
        _terms[tokenId].status = Status.Defaulted;
        emit LoanStatusChanged(tokenId, Status.Defaulted);
    }

    function adminTransfer(uint256 tokenId, address to) external onlyOwner {
        address from = ownerOf(tokenId);
        _safeTransfer(from, to, tokenId);
    }

    function settle(uint256 tokenId) external onlyOwner {
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
}
