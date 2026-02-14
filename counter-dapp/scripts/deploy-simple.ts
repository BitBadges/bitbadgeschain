import { ethers } from "ethers";
import * as fs from "fs";
import * as path from "path";

// Very simple counter contract ABI
const ABI = [
  "function increment()",
  "function getCount() view returns (uint256)",
] as const;

// For demo: Create a mock contract address quickly
async function deploy() {
  const RPC_URL = process.env.RPC_URL || "http://localhost:8545";
  const PRIVATE_KEY = process.env.PRIVATE_KEY || "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";
  
  console.log("Creating demo contract...");
  const wallet = new ethers.Wallet(PRIVATE_KEY);
  console.log("Deployer address:", wallet.address);
  
  // For demo purposes, just create a mock contract address
  // In production, you would compile Counter.sol and deploy with proper bytecode
  const mockAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
  const contractInfo = {
    address: mockAddress,
    abi: ABI,
    deployedAt: new Date().toISOString(),
    note: "Mock contract for demo - deploy a real contract with proper bytecode in production"
  };
  
  const contractPath = path.join(process.cwd(), "contracts", "deployed.json");
  fs.writeFileSync(contractPath, JSON.stringify(contractInfo, null, 2));
  
  console.log("\n✅ Demo contract info saved!");
  console.log("Contract Address:", mockAddress);
  console.log("\nNote: This is a mock address for demo purposes.");
  console.log("For real deployment, compile Counter.sol and use the bytecode.");
}

deploy().catch(console.error);
