import { expect } from "chai";
import { network } from "hardhat";
import type { LoanNote } from "../types/ethers-contracts/LoanNote.js";
import type { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/types";

const { ethers } = await network.getOrCreate();

const PRINCIPAL = 10_000n; // cents
const APR_BPS = 500;
const MATURITY = 4_102_444_800n; // 2100-01-01, comfortably in the future
const TOKEN_URI = "http://localhost:8080/loans/1/terms";

describe("LoanNote", () => {
  let note: LoanNote;
  let admin: HardhatEthersSigner;
  let lender: HardhatEthersSigner;
  let buyer: HardhatEthersSigner;
  let outsider: HardhatEthersSigner;

  beforeEach(async () => {
    [admin, lender, buyer, outsider] = await ethers.getSigners();
    note = await ethers.deployContract("LoanNote", [admin.address]);
  });

  function originate(to = lender.address) {
    return note.originate(to, PRINCIPAL, APR_BPS, MATURITY, TOKEN_URI);
  }

  describe("roles", () => {
    it("grants admin and both business roles to the deployer-named admin", async () => {
      expect(await note.hasRole(await note.DEFAULT_ADMIN_ROLE(), admin.address)).to.equal(true);
      expect(await note.hasRole(await note.ORIGINATOR_ROLE(), admin.address)).to.equal(true);
      expect(await note.hasRole(await note.SERVICER_ROLE(), admin.address)).to.equal(true);
    });

    it("grants no roles to anyone else", async () => {
      for (const role of [
        await note.DEFAULT_ADMIN_ROLE(),
        await note.ORIGINATOR_ROLE(),
        await note.SERVICER_ROLE(),
      ]) {
        expect(await note.hasRole(role, lender.address)).to.equal(false);
      }
    });

    it("admin can grant a business role to another account", async () => {
      const servicerRole = await note.SERVICER_ROLE();
      await note.grantRole(servicerRole, outsider.address);
      expect(await note.hasRole(servicerRole, outsider.address)).to.equal(true);
    });
  });

  describe("originate", () => {
    it("mints the note to the lender, not the caller", async () => {
      await expect(originate())
        .to.emit(note, "LoanOriginated")
        .withArgs(0n, lender.address, PRINCIPAL, APR_BPS, MATURITY);
      expect(await note.ownerOf(0n)).to.equal(lender.address);
    });

    it("records the terms", async () => {
      await originate();
      const terms = await note.terms(0n);
      expect(terms.principal).to.equal(PRINCIPAL);
      expect(terms.aprBps).to.equal(APR_BPS);
      expect(terms.maturity).to.equal(MATURITY);
      expect(terms.status).to.equal(0n); // Active
    });

    it("rejects callers without ORIGINATOR_ROLE", async () => {
      await expect(
        note.connect(outsider).originate(lender.address, PRINCIPAL, APR_BPS, MATURITY, TOKEN_URI),
      )
        .to.be.revertedWithCustomError(note, "AccessControlUnauthorizedAccount")
        .withArgs(outsider.address, await note.ORIGINATOR_ROLE());
    });
  });

  describe("transfers", () => {
    beforeEach(async () => {
      await originate();
    });

    it("the owner can transfer their note", async () => {
      await note.connect(lender).transferFrom(lender.address, buyer.address, 0n);
      expect(await note.ownerOf(0n)).to.equal(buyer.address);
    });

    it("the platform cannot move a note it does not hold, roles or not", async () => {
      await expect(note.connect(admin).transferFrom(lender.address, buyer.address, 0n))
        .to.be.revertedWithCustomError(note, "ERC721InsufficientApproval")
        .withArgs(admin.address, 0n);
    });

    it("adminTransfer no longer exists", async () => {
      expect(note.interface.getFunction("adminTransfer")).to.equal(null);
    });

    it("a defaulted note is still an assignable claim", async () => {
      await note.markDefaulted(0n);
      await note.connect(lender).transferFrom(lender.address, buyer.address, 0n);
      expect(await note.ownerOf(0n)).to.equal(buyer.address);
    });
  });

  describe("settle", () => {
    beforeEach(async () => {
      await originate();
    });

    it("burns the note and clears the terms", async () => {
      await expect(note.settle(0n))
        .to.emit(note, "LoanStatusChanged")
        .withArgs(0n, 1n) // Repaid
        .and.to.emit(note, "LoanSettled")
        .withArgs(0n);
      await expect(note.ownerOf(0n)).to.be.revertedWithCustomError(note, "ERC721NonexistentToken");
      await expect(note.terms(0n)).to.be.revertedWithCustomError(note, "ERC721NonexistentToken");
    });

    it("burns regardless of who holds the note: final repayment extinguishes the claim", async () => {
      await note.connect(lender).transferFrom(lender.address, buyer.address, 0n);
      await note.settle(0n);
      await expect(note.ownerOf(0n)).to.be.revertedWithCustomError(note, "ERC721NonexistentToken");
    });

    it("rejects callers without SERVICER_ROLE, including the note's owner", async () => {
      await expect(note.connect(lender).settle(0n))
        .to.be.revertedWithCustomError(note, "AccessControlUnauthorizedAccount")
        .withArgs(lender.address, await note.SERVICER_ROLE());
    });
  });

  describe("markDefaulted", () => {
    beforeEach(async () => {
      await originate();
    });

    it("flags the terms and emits a status change", async () => {
      await expect(note.markDefaulted(0n))
        .to.emit(note, "LoanStatusChanged")
        .withArgs(0n, 2n); // Defaulted
      expect((await note.terms(0n)).status).to.equal(2n);
    });

    it("rejects callers without SERVICER_ROLE", async () => {
      await expect(note.connect(lender).markDefaulted(0n))
        .to.be.revertedWithCustomError(note, "AccessControlUnauthorizedAccount")
        .withArgs(lender.address, await note.SERVICER_ROLE());
    });
  });

  describe("supportsInterface", () => {
    it("reports ERC-721 and AccessControl", async () => {
      expect(await note.supportsInterface("0x80ac58cd")).to.equal(true); // ERC-721
      expect(await note.supportsInterface("0x7965db0b")).to.equal(true); // AccessControl
    });
  });
});
