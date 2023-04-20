# bitbadgeschain
**bitbadgeschain** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli). BitBadges offers an open-source, community-driven suite of tools focused on cross-chain issuance of digital tokens (badges). This blockchain is the core of the BitBadges ecosystem.

See the [BitBadges documentation](https://docs.bitbadges.io/overview) to learn about BitBadges and the BitBadges blockchain.

See the [Cosmos SDK Docs](https://docs.cosmos.network) to learn about the Cosmos SDK and Tendermint. This repository
follows the Cosmos SDK's [directory structure](https://docs.cosmos.network/master/building-modules/module-manager.html#directory-structure).

## Acknowledgements
We want to acknowledge the [Evmos](
    https://github.com/evmos
) project for maintaining and open-sourcing their software, which BitBadges 
was built upon. 

## Get started
```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend
See the [BitBadges documentation](https://docs.bitbadges.io/overview) and [BitBadges website](https://bitbadges.io) for more information.

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/bitbadges/bitbadgeschain@latest! | sudo bash
```
`bitbadges/bitbadgeschain` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more about Ignite CLI and Cosmos SDK

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
