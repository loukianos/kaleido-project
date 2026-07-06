import hardhatToolboxMochaEthers from "@nomicfoundation/hardhat-toolbox-mocha-ethers";
import { defineConfig } from "hardhat/config";
import dotenv from "dotenv";

dotenv.config({ path: "../.env" });

const RPC_URL = process.env.ETH_RPC_URL ?? "http://127.0.0.1:31545";
// Throwaway dev key by default
const DEPLOYER_PRIVATE_KEY =
  process.env.DEPLOYER_PRIVATE_KEY ??
  "0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63";

export default defineConfig({
  plugins: [hardhatToolboxMochaEthers],
  solidity: {
    version: "0.8.28",
    settings: {
      optimizer: { enabled: true, runs: 200 },
      evmVersion: "prague",
      // Keep bytecode a pure function of source + settings so the committed
      // Go bindings can be diff-checked reproducibly across machines and CI.
      metadata: { bytecodeHash: "none" },
    },
  },
  paths: { sources: "src" },
  networks: {
    besu: {
      type: "http",
      chainType: "l1",
      url: RPC_URL,
      chainId: 1337,
      accounts: [DEPLOYER_PRIVATE_KEY],
    },
  },
});
