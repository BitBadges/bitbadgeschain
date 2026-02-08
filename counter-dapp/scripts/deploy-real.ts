import { ethers } from "ethers";
import * as fs from "fs";
import * as path from "path";

async function deploy() {
  // Try port 8545 first (EVM JSON-RPC), fallback to 26657
  const RPC_URL = process.env.RPC_URL || "http://localhost:8545";
  // Use alice's key by default (has funds from genesis)
  // Get it with: bitbadgeschaind keys export alice --keyring-backend test --unarmored-hex
  const PRIVATE_KEY = process.env.PRIVATE_KEY || "";
  
  if (!PRIVATE_KEY) {
    console.error("❌ PRIVATE_KEY not set!");
    console.error("Get a private key with funds:");
    console.error("  bitbadgeschaind keys export alice --keyring-backend test --unarmored-hex");
    console.error("Then set: PRIVATE_KEY=0x... bun run deploy");
    process.exit(1);
  }
  
  console.log("Connecting to RPC:", RPC_URL);
  // Use static network to avoid auto-detection issues
  const provider = new ethers.JsonRpcProvider(RPC_URL, {
    chainId: 90123,
    name: "bitbadges-1"
  }, { staticNetwork: true });
  const wallet = new ethers.Wallet(PRIVATE_KEY, provider);
  
  console.log("Deployer address:", wallet.address);
  
  // Try to get balance with timeout
  let balance = 0n;
  try {
    balance = await Promise.race([
      provider.getBalance(wallet.address),
      new Promise<bigint>((_, reject) => setTimeout(() => reject(new Error("Timeout")), 5000))
    ]);
    console.log("Balance:", ethers.formatEther(balance), "BADGE");
  } catch (error: any) {
    console.log("⚠️  Could not check balance:", error.message);
    console.log("Continuing anyway...");
  }
  
  // Read compiled contract
  const binPath = path.join(process.cwd(), "contracts", "Counter.bin");
  const abiPath = path.join(process.cwd(), "contracts", "Counter.abi.json");
  
  if (!fs.existsSync(binPath) || !fs.existsSync(abiPath)) {
    console.error("❌ Contract not compiled.");
    console.error("Run: bun run compile");
    process.exit(1);
  }
  
  const bytecode = "0x" + fs.readFileSync(binPath, "utf-8").trim();
  const abi = JSON.parse(fs.readFileSync(abiPath, "utf-8"));
  
  console.log("\n📦 Deploying contract...");
  console.log("Bytecode length:", bytecode.length, "chars");
  
  if (balance === 0n) {
    console.warn("⚠️  Warning: No balance detected. Transaction may fail.");
    console.warn("You may need to fund the deployer address with BADGE tokens.");
  }
  
  try {
    const factory = new ethers.ContractFactory(abi, bytecode, wallet);
    console.log("Creating contract instance...");
    const contract = await factory.deploy(0); // Start with count = 0
    
    console.log("⏳ Waiting for deployment transaction...");
    await contract.waitForDeployment();
    const address = await contract.getAddress();
    
    console.log("\n✅ Contract deployed successfully!");
    console.log("Contract Address:", address);
    
    // Verify deployment by calling getCount
    try {
      const initialCount = await contract.getCount();
      console.log("Initial count:", initialCount.toString());
    } catch (e) {
      console.warn("Could not verify contract (this is ok)");
    }
    
    // Save contract info
    const contractInfo = {
      address,
      abi,
      deployedAt: new Date().toISOString(),
    };
    
    const contractPath = path.join(process.cwd(), "contracts", "deployed.json");
    fs.writeFileSync(contractPath, JSON.stringify(contractInfo, null, 2));
    console.log("\n💾 Contract info saved to:", contractPath);
    
  } catch (error: any) {
    console.error("\n❌ Deployment failed!");
    console.error("Error:", error.message);
    if (error.message.includes("balance") || error.message.includes("insufficient funds") || error.message.includes("gas")) {
      console.error("\n💡 Solution: Fund the deployer address with BADGE tokens.");
      console.error("Deployer address:", wallet.address);
      console.error("\nYou can send tokens using:");
      console.error(`bitbadgeschaind tx bank send alice ${wallet.address} 1000000000000000000ubadge --keyring-backend test --chain-id bitbadges-1 --yes`);
    }
    if (error.message.includes("revert") || error.message.includes("execution reverted")) {
      console.error("\n💡 The contract bytecode might be invalid. Try recompiling.");
    }
    process.exit(1);
  }
}

deploy().catch(console.error);
