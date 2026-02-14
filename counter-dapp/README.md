# Counter dApp

A simple counter dApp built with Next.js, Tailwind CSS, React, and MetaMask integration for BitBadges Chain.

## Features

- 🔗 MetaMask wallet connection
- 📊 Simple counter contract interaction
- 🎨 Modern UI with Tailwind CSS
- ⚡ Built with Next.js 14 and React

## Prerequisites

- Bun installed
- BitBadges chain running locally (http://localhost:26657)
- MetaMask configured with BitBadges chain (Chain ID: 90123)

## Setup

1. Install dependencies:
```bash
bun install
```

2. Deploy the counter contract:
```bash
# First, get a private key from your chain
bitbadgeschaind keys export alice --keyring-backend test --unarmored-hex

# Then deploy (replace with your private key)
PRIVATE_KEY=your_private_key_here bun run deploy-contract
```

Or create a `.env` file:
```bash
cp .env.example .env
# Edit .env and add your PRIVATE_KEY
bun run deploy-contract
```

**Note:** Make sure your chain is running and EVM JSON-RPC is enabled. The script will try multiple RPC endpoints automatically.

3. Start the development server:
```bash
bun run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser

## MetaMask Setup

To connect MetaMask to your local BitBadges chain:

1. Open MetaMask
2. Go to Settings > Networks > Add Network
3. Add the following network:
   - Network Name: BitBadges Local
   - RPC URL: http://localhost:26657
   - Chain ID: 90123
   - Currency Symbol: BADGE

## Usage

1. Connect your MetaMask wallet
2. Click "Increment Counter" to increment the counter
3. The count will update automatically after the transaction is confirmed

## Project Structure

```
counter-dapp/
├── app/              # Next.js app directory
│   ├── page.tsx      # Main page component
│   ├── providers.tsx # Wallet providers setup
│   └── globals.css   # Global styles
├── contracts/        # Smart contracts
│   ├── Counter.sol   # Counter contract source
│   └── deployed.json # Deployed contract info
├── scripts/          # Deployment scripts
│   └── deploy.ts     # Contract deployment script
└── package.json      # Dependencies
```

## Troubleshooting

- **Contract not deployed**: Make sure you've run `bun run deploy-contract` first
- **Transaction fails**: Ensure you have enough BADGE tokens in your wallet
- **Can't connect wallet**: Make sure MetaMask is installed and configured with the correct network

