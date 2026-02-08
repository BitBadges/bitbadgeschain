import * as fs from "fs";
import * as path from "path";
import solc from "solc";

async function compile() {
  const contractPath = path.join(process.cwd(), "contracts", "Counter.sol");
  const source = fs.readFileSync(contractPath, "utf-8");
  
  const input = {
    language: "Solidity",
    sources: {
      "Counter.sol": {
        content: source,
      },
    },
    settings: {
      outputSelection: {
        "*": {
          "*": ["abi", "evm.bytecode"],
        },
      },
    },
  };
  
  console.log("Compiling Counter.sol...");
  const output = JSON.parse(solc.compile(JSON.stringify(input)));
  
  if (output.errors) {
    const errors = output.errors.filter((e: any) => e.severity === "error");
    if (errors.length > 0) {
      console.error("Compilation errors:");
      errors.forEach((error: any) => console.error(error.message));
      process.exit(1);
    }
  }
  
  const contract = output.contracts["Counter.sol"]["Counter"];
  const abi = contract.abi;
  const bytecode = contract.evm.bytecode.object;
  
  // Save ABI
  const abiPath = path.join(process.cwd(), "contracts", "Counter.abi.json");
  fs.writeFileSync(abiPath, JSON.stringify(abi, null, 2));
  
  // Save bytecode
  const binPath = path.join(process.cwd(), "contracts", "Counter.bin");
  fs.writeFileSync(binPath, bytecode);
  
  console.log("✅ Compilation successful!");
  console.log("ABI saved to:", abiPath);
  console.log("Bytecode saved to:", binPath);
  
  return { abi, bytecode };
}

compile().catch(console.error);

