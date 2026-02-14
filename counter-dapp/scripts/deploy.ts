import { ethers } from "ethers";
import * as fs from "fs";
import * as path from "path";
import { COUNTER_ABI, COUNTER_BYTECODE } from "./compile-contract.js";

async function deploy() {
  // Try multiple RPC endpoints - EVM JSON-RPC might be on different ports
  const RPC_URLS = [
    process.env.RPC_URL,
    "http://localhost:8545", // Standard EVM JSON-RPC port
    "http://localhost:26657", // Tendermint RPC (might support EVM)
  ].filter(Boolean) as string[];
  
  const PRIVATE_KEY = process.env.PRIVATE_KEY || "";
  
  if (!PRIVATE_KEY) {
    console.error("Please set PRIVATE_KEY environment variable");
    console.log("You can get a private key from: bitbadgeschaind keys export alice --keyring-backend test --unarmored-hex");
    process.exit(1);
  }

  // Try to connect to any available RPC
  let provider: ethers.JsonRpcProvider | null = null;
  let workingRpc = "";
  
  for (const rpcUrl of RPC_URLS) {
    try {
      const testProvider = new ethers.JsonRpcProvider(rpcUrl);
      const blockNumber = await testProvider.getBlockNumber();
      provider = testProvider;
      workingRpc = rpcUrl;
      console.log(`Connected to RPC: ${rpcUrl}`);
      console.log(`Current block: ${blockNumber}`);
      break;
    } catch (error) {
      console.log(`Failed to connect to ${rpcUrl}, trying next...`);
    }
  }
  
  if (!provider) {
    console.error("Failed to connect to any RPC endpoint.");
    console.error("Make sure your chain is running and EVM JSON-RPC is enabled.");
    console.error("Tried:", RPC_URLS.join(", "));
    process.exit(1);
  }
  
  const wallet = new ethers.Wallet(PRIVATE_KEY, provider);
  
  // Check network
  try {
    const network = await provider.getNetwork();
    console.log("Network Chain ID:", network.chainId.toString());
  } catch (error) {
    console.warn("Could not get network info, continuing...");
  }

  console.log("Deploying Counter contract...");
  console.log("Deployer address:", wallet.address);
  
  // Check balance
  const balance = await provider.getBalance(wallet.address);
  console.log("Deployer balance:", ethers.formatEther(balance), "ETH/BADGE");
  
  if (balance === 0n) {
    console.warn("Warning: Deployer has no balance. Transactions may fail.");
    console.warn("You may need to fund this address with BADGE tokens.");
  }

  // Deploy contract
  console.log("Deploying contract...");
  const factory = new ethers.ContractFactory(COUNTER_ABI, COUNTER_BYTECODE, wallet);
  const contract = await factory.deploy(0); // Start with count = 0

  console.log("Waiting for deployment confirmation...");
  await contract.waitForDeployment();
  const address = await contract.getAddress();

  console.log("Contract deployed at:", address);

  // Save contract address and ABI
  const contractInfo = {
    address,
    abi: COUNTER_ABI,
    deployedAt: new Date().toISOString(),
  };

  const contractPath = path.join(process.cwd(), "contracts", "deployed.json");
  fs.writeFileSync(contractPath, JSON.stringify(contractInfo, null, 2));

  console.log("Contract info saved to:", contractPath);
  console.log("\n✅ Deployment complete!");
  console.log("Contract Address:", address);
}

deploy().catch((error) => {
  console.error("Deployment failed:", error);
  process.exit(1);
});

