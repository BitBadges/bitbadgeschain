"use client";

import { RainbowKitProvider } from "@rainbow-me/rainbowkit";
import "@rainbow-me/rainbowkit/styles.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { defineChain } from "viem";
import { createConfig, http, WagmiProvider } from "wagmi";

// Create a custom chain config for BitBadges
// ubadge is the base unit - if it has 18 decimals in EVM context
const bitbadgesChain = defineChain({
  id: 50025, // BitBadges Testnet EVM Chain ID (for local development)
  name: "BitBadges",
  nativeCurrency: {
    decimals: 18,
    name: "Badge",
    symbol: "BADGE",
  },
  rpcUrls: {
    default: {
      http: ["http://localhost:8545"],
    },
  },
  blockExplorers: {
    default: {
      name: "Local",
      url: "http://localhost:26657",
    },
  },
  testnet: true,
});

// Create config with proper polling settings for local chain
const config = createConfig({
  chains: [bitbadgesChain],
  transports: {
    [bitbadgesChain.id]: http("http://localhost:8545", {
      // Polling interval for transaction receipts (2 seconds)
      retryCount: 3,
      retryDelay: 1000,
    }),
  },
  ssr: true,
});

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            // Refetch every 4 seconds for live updates
            refetchInterval: 4000,
            staleTime: 2000,
          },
        },
      })
  );

  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider>{children}</RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}

