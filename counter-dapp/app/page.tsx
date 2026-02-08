"use client";

import { ConnectButton } from "@rainbow-me/rainbowkit";
import { useReadContract, useWriteContract, useWaitForTransactionReceipt, useBlockNumber } from "wagmi";
import { useState, useEffect } from "react";
import deployedContract from "@/contracts/deployed.json";

export default function Home() {
  const [contractAddress, setContractAddress] = useState<string>("");
  const [contractAbi, setContractAbi] = useState<any[]>([]);
  const [txStatus, setTxStatus] = useState<string>("");

  useEffect(() => {
    // Load deployed contract info
    if (deployedContract.address) {
      setContractAddress(deployedContract.address);
      setContractAbi(deployedContract.abi);
    }
  }, []);

  // Watch block number for polling indication
  const { data: blockNumber } = useBlockNumber({ watch: true });

  const { data: count, refetch: refetchCount } = useReadContract({
    address: contractAddress as `0x${string}`,
    abi: contractAbi,
    functionName: "getCount",
    query: {
      enabled: !!contractAddress,
      refetchInterval: 3000, // Poll every 3 seconds
    },
  });

  const { writeContract, data: hash, isPending, error: writeError, reset: resetWrite } = useWriteContract();
  const { isLoading: isConfirming, isSuccess, error: txError } = useWaitForTransactionReceipt({
    hash,
    confirmations: 1,
    query: {
      refetchInterval: 2000, // Poll for receipt every 2 seconds
    },
  });

  const handleIncrement = () => {
    if (!contractAddress) {
      alert("Contract not deployed. Please run: bun run deploy-contract");
      return;
    }

    setTxStatus("Sending transaction...");
    resetWrite();

    writeContract(
      {
        address: contractAddress as `0x${string}`,
        abi: contractAbi,
        functionName: "increment",
        gas: BigInt(100000), // Set explicit gas limit
      },
      {
        onSuccess: () => {
          setTxStatus("Transaction sent, waiting for confirmation...");
        },
        onError: (error) => {
          console.error("Write contract error:", error);
          setTxStatus(`Error: ${error.message}`);
        },
      }
    );
  };

  useEffect(() => {
    if (isSuccess) {
      setTxStatus("Transaction confirmed!");
      refetchCount();
      // Clear success message after 5 seconds
      setTimeout(() => setTxStatus(""), 5000);
    }
  }, [isSuccess, refetchCount]);

  useEffect(() => {
    if (isConfirming) {
      setTxStatus("Waiting for confirmation...");
    }
  }, [isConfirming]);

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-2xl mx-auto">
          <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8">
            <div className="flex justify-between items-center mb-8">
              <h1 className="text-4xl font-bold text-gray-900 dark:text-white">
                Counter dApp
              </h1>
              <ConnectButton />
            </div>

            <div className="mt-8 space-y-6">
              {!contractAddress ? (
                <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
                  <p className="text-yellow-800 dark:text-yellow-200">
                    Contract not deployed. Please deploy the contract first:
                  </p>
                  <code className="block mt-2 p-2 bg-yellow-100 dark:bg-yellow-900/40 rounded text-sm">
                    bun run deploy-contract
                  </code>
                </div>
              ) : (
                <>
                  <div className="text-center">
                    <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Contract Address
                    </p>
                    <p className="font-mono text-xs text-gray-800 dark:text-gray-200 break-all">
                      {contractAddress}
                    </p>
                  </div>

                  <div className="bg-indigo-50 dark:bg-indigo-900/20 rounded-lg p-8 text-center">
                    <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Current Count
                    </p>
                    <p className="text-6xl font-bold text-indigo-600 dark:text-indigo-400">
                      {count !== undefined ? count.toString() : "..."}
                    </p>
                  </div>

                  <button
                    onClick={handleIncrement}
                    disabled={isPending || isConfirming || !contractAddress}
                    className="w-full bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-semibold py-4 px-6 rounded-lg transition-colors duration-200 shadow-lg hover:shadow-xl"
                  >
                    {isPending
                      ? "Confirming..."
                      : isConfirming
                      ? "Processing..."
                      : "Increment Counter"}
                  </button>

                  {txStatus && !writeError && !txError && !isSuccess && (
                    <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                      <p className="text-blue-800 dark:text-blue-200">
                        ⏳ {txStatus}
                      </p>
                      {hash && (
                        <p className="text-blue-600 dark:text-blue-400 text-xs mt-1 font-mono break-all">
                          Tx: {hash}
                        </p>
                      )}
                    </div>
                  )}
                  {writeError && (
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                      <p className="text-red-800 dark:text-red-200 font-semibold">
                        ❌ Transaction Error
                      </p>
                      <p className="text-red-700 dark:text-red-300 text-sm mt-1 break-all">
                        {writeError.message}
                      </p>
                    </div>
                  )}
                  {txError && (
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                      <p className="text-red-800 dark:text-red-200 font-semibold">
                        ❌ Transaction Failed
                      </p>
                      <p className="text-red-700 dark:text-red-300 text-sm mt-1 break-all">
                        {txError.message}
                      </p>
                    </div>
                  )}
                  {isSuccess && (
                    <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
                      <p className="text-green-800 dark:text-green-200">
                        ✅ Transaction successful! Count updated.
                      </p>
                      {hash && (
                        <p className="text-green-600 dark:text-green-400 text-xs mt-1 font-mono break-all">
                          Tx: {hash}
                        </p>
                      )}
                    </div>
                  )}
                </>
              )}
            </div>
          </div>

          <div className="mt-8 text-center text-sm text-gray-600 dark:text-gray-400">
            <p>Built for BitBadges Chain</p>
            <p className="mt-1">EVM RPC: http://localhost:8545</p>
            {blockNumber && (
              <p className="mt-1 text-xs">Block: {blockNumber.toString()}</p>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}

