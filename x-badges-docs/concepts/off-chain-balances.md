# Off-Chain Balances

> **⚠️ DEPRECATED**: Off-chain balance types have been deprecated and are no longer recommended for new collections. Always use `"Standard"` balance type for new collections.

## Overview

Off-chain balances were a legacy feature that allowed badge collections to store balance information and handle transferability outside of the blockchain. This approach has been deprecated in favor of off-chain only claims, address lists, and other off-chain features.

## Deprecated Balance Types

### Off-Chain - Indexed

```json
"balancesType": "Off-Chain - Indexed"
```

This balance type was used for collections where all balances were stored off-chain in a single JSON file that was fetched by the client.

### Off-Chain - Non-Indexed

```json
"balancesType": "Off-Chain - Non-Indexed"
```

This balance type was used for collections where balance information was stored off-chain in a dynamic on-demand manner (e.g. /balances/{address} endpoint).

## Deprecated Metadata Structure

Collections using off-chain balance types previously used the `offChainBalancesMetadataTimeline` field to store metadata about where and how to fetch balance information:

```json
{
    "offChainBalancesMetadataTimeline": [
        {
            "offChainBalancesMetadata": {
                "uri": "https://example.com/balances",
                "customData": ""
            },
            "timelineTimes": [
                {
                    "start": 1,
                    "end": 1000
                }
            ]
        }
    ]
}
```

This field has been deprecated and should not be used for new collections.

## Migration to Standard Balances

### For New Collections

Always use the "Standard" balance type for new collections:

```json
{
    "balancesType": "Standard",
    "offChainBalancesMetadataTimeline": []
}
```
