import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const LoanNoteModule = buildModule("LoanNoteModule", (m) => {
  const owner = m.getAccount(0);

  const loanNote = m.contract("LoanNote", [owner]);

  return { loanNote };
});

export default LoanNoteModule;
