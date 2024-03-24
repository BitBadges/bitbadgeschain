# BitBadges Blockchian

**bitbadgeschain** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli). BitBadges offers an open-source, community-driven suite of tools focused on cross-chain issuance of digital tokens (badges). This blockchain is the core of the BitBadges ecosystem.

See the [BitBadges documentation](https://docs.bitbadges.io/overview) to learn about BitBadges and the BitBadges blockchain.

See the [Cosmos SDK Docs](https://docs.cosmos.network) to learn about the Cosmos SDK and Tendermint. This repository
follows the Cosmos SDK's [directory structure](https://docs.cosmos.network/master/building-modules/module-manager.html#directory-structure).

## Building with Makefile

The following instructions are for Ubuntu 23.10. If you are using a different operating system, you may need to modify the commands.

If you do not have the following dependencies installed, you will need to install them before you can build the blockchain.

```bash
sudo apt-get install git curl make build-essential gcc
```

To build the BitBadges blockchain from source, run the following:

```bash
snap install go --classic # Install Go 1.21
```

To build the BitBadges blockchain from source, run the following:

```bash
make build-all
# OR
make build-linux/amd64
make build-darwin/amd64
make build-linux/arm64
```

For building linux/arm64, you will need to have a cross-compilation toolchain installed. You can install it with the following command:

```bash
sudo apt-get install gcc-aarch64-linux-gnu
```

For building darwin/amd64, you will need to have the o64-clang cross-compilation toolchain installed. We refer you to the [https://github.com/tpoechtrager/osxcross](https://github.com/tpoechtrager/osxcross) project for more information.

## Building / Serving With Ignite CLI

This blockchain was also built using the Ignite CLI. To build and serve the blockchain, download the Ignite CLI from the [Ignite CLI website](https://ignite.com/cli).

Then, run the following commands to build and serve the blockchain:

```
ignite chain init --skip-proto
ignite chain build --skip-proto
ignite chain serve --skip-proto
```

You will have to use the --skip-proto flag because we manually correct a query file to fix a small bug in the generated code.

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

## Web Frontend

See the [BitBadges documentation](https://docs.bitbadges.io/overview) and [BitBadges website](https://bitbadges.io) for more information.

## Release

To release a new version of the blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

## Development

The BitBadges blockchain is open-source and community-driven. We welcome contributions from the community. To contribute, fork this repository and submit a pull request.

Couple of development notes:

-   The `x` directory contains the modules of the blockchain. Each module is a separate directory.
-   The `proto` directory contains the protobuf files for the modules. These files are used to generate the Go code for the modules. This has typically been done with `ignite generate proto-go`.
-   The `chain-handlers` directory contains the handlers for the blockchain. These handlers are used to handle the signature logic for each respective blockchain that is supported. Ethereum uses EIP712 signatures. Solana and Bitcoin use a JSON schema with everything alphabetically sorted. Cosmos uses typical Cosmos signatures. Learn more on the BitBadges documentation and via bitbadgesjs.

## Learn more about Ignite CLI and Cosmos SDK

-   [Ignite CLI](https://ignite.com/cli)
-   [Tutorials](https://docs.ignite.com/guide)
-   [Ignite CLI docs](https://docs.ignite.com)
-   [Cosmos SDK docs](https://docs.cosmos.network)
-   [Developer Chat](https://discord.gg/ignite)

## Additional Information

The code related to handling EIP-712 signatures is forked and adapted from the ethermint repository [here](https://github.com/evmos/ethermint)
licensed under LGPL-3.0. This repository is also licensed under LGPL-3.0.


