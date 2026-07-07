import hardhatToolboxMochaEthers from "@nomicfoundation/hardhat-toolbox-mocha-ethers";
import { defineConfig } from "hardhat/config";
import dotenv from "dotenv";

dotenv.config({ path: "../.env" });

const RPC_URL = process.env.ETH_RPC_URL ?? "http://127.0.0.1:31545";
// Deterministic local-only fallback; set DEPLOYER_PRIVATE_KEY for real deploys.
const DEPLOYER_PRIVATE_KEY =
  process.env.DEPLOYER_PRIVATE_KEY ?? `0x${"1".repeat(64)}`;

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
