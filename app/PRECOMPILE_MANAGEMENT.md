# Precompile Registration and Enablement

This document explains how precompiles are registered and enabled in the BitBadges chain.

## Overview

Precompiles require **two steps** to be callable:
1. **Registration** - Makes the precompile available to the EVM keeper
2. **Enablement** - Activates the precompile so it can be called

## Precompile Types

### Default Cosmos Precompiles

These are provided by the `cosmos/evm` module and include:
- Staking precompile
- Distribution precompile
- Bank precompile
- Governance precompile
- And others (see `cosmos/evm/precompiles/types`)

**Registration**: Automatically registered via `WithStaticPrecompiles(DefaultStaticPrecompiles(...))` in `app/evm.go:164-176`

**Enablement**: Must be enabled via genesis `active_static_precompiles` array or governance

### Custom BitBadges Precompiles

Currently implemented:
- **0x0000000000000000000000000000000000001001** - Tokenization precompile
- **0x0000000000000000000000000000000000001002** - Gamm precompile
- **0x0000000000000000000000000000000000001003** - SendManager precompile

**Registration**: Registered in `app/evm.go:registerCustomPrecompiles()`

**Enablement**: Must be enabled via genesis `active_static_precompiles` array or upgrade handler

## Address Space

To prevent collisions, follow this convention:
- **0x0800-0x0806**: Reserved for default Cosmos precompiles
- **0x1001+**: Reserved for custom BitBadges precompiles
  - 0x1001: Tokenization precompile
  - 0x1002: Gamm precompile
  - 0x1003: SendManager precompile
  - 0x1004+: Available for future custom precompiles

## Registration Flow

### In `app/evm.go`

1. **Default precompiles** are registered via `WithStaticPrecompiles()` during EVM keeper creation
2. **Custom precompiles** are registered via `registerCustomPrecompiles()` after EVM keeper creation

```go
// Default precompiles (automatic)
app.EVMKeeper = configureEVMKeeper(evmkeeper.NewKeeper(...).WithStaticPrecompiles(
    precompiletypes.DefaultStaticPrecompiles(...),
))

// Custom precompiles (manual)
app.registerCustomPrecompiles()
```

### Adding a New Custom Precompile

1. **Define the address** in the precompile package:
   ```go
   const MyPrecompileAddress = "0x0000000000000000000000000000000000001003"
   ```

2. **Register in `app/evm.go:registerCustomPrecompiles()`**:
   ```go
   myPrecompile := myprecompile.NewPrecompile(app.MyKeeper)
   myPrecompileAddr := common.HexToAddress(myprecompile.MyPrecompileAddress)
   app.EVMKeeper.RegisterStaticPrecompile(myPrecompileAddr, myPrecompile)
   ```

3. **Add to `GetAllPrecompileAddresses()`** in `app/evm.go`:
   ```go
   addresses = append(addresses, common.HexToAddress(myprecompile.MyPrecompileAddress))
   ```

4. **Add to `GetAllCustomPrecompileAddresses()`** in `app/precompile_helpers.go`:
   ```go
   return []common.Address{
       // ... existing addresses ...
       common.HexToAddress(myprecompile.MyPrecompileAddress),
   }
   ```

5. **Update `config.yml`** to include in `active_static_precompiles`:
   ```yaml
   active_static_precompiles:
     - "0x0000000000000000000000000000000000001003"  # My precompile
   ```

6. **Update upgrade handler** if needed (see below)

## Enablement Flow

### Genesis (New Chains)

Add precompile addresses to `config.yml`:
```yaml
evm:
  params:
    active_static_precompiles:
      - "0x0000000000000000000000000000000000001001"  # Tokenization
      - "0x0000000000000000000000000000000000001002"  # Gamm
      - "0x0000000000000000000000000000000000001003"  # SendManager
```

### Upgrade Handler (Existing Chains)

Use the centralized helper in `app/upgrades/v24/upgrades.go`:
```go
customPrecompileAddresses := []common.Address{
    common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress),
    common.HexToAddress(gammprecompile.GammPrecompileAddress),
    common.HexToAddress(sendmanagerprecompile.SendManagerPrecompileAddress),
    // Add new precompile here
}

for _, addr := range customPrecompileAddresses {
    if err := evmKeeper.EnableStaticPrecompiles(sdkCtx, addr); err != nil {
        sdkCtx.Logger().Info("Precompile enable attempt", "error", err, "address", addr.Hex())
    }
}
```

### Tests

Use the helper function from `app/precompile_helpers.go`:
```go
import apphelpers "github.com/bitbadges/bitbadgeschain/app"

err := apphelpers.RegisterAndEnableAllPrecompiles(
    suite.Ctx,
    suite.EVMKeeper,
    suite.TokenizationKeeper,
    suite.GammKeeper,
    suite.SendmanagerKeeper,
)
suite.Require().NoError(err)
```

## Validation

### Address Collision Detection

The system automatically validates that custom precompile addresses don't collide:
- Called during app initialization in `app/evm.go`
- Uses `ValidateNoAddressCollisions()` from `app/precompile_helpers.go`
- Panics on startup if collisions are detected

### Helper Functions

**`GetAllCustomPrecompileAddresses()`**: Returns all custom precompile addresses
- Useful for validation and testing
- Located in `app/precompile_helpers.go`

**`GetAllPrecompileAddresses()`**: Returns all custom precompile addresses (app method)
- Used by `EnableAllPrecompiles()`
- Located in `app/evm.go`

**`EnableAllPrecompiles(ctx)`**: Enables all custom precompiles
- Idempotent (safe to call multiple times)
- Located in `app/evm.go`

## Current Status

✅ **Registered and Enabled**:
- Tokenization precompile (0x1001) - Create collections, transfer tokens, manage approvals
- Gamm precompile (0x1002) - AMM liquidity pool operations
- SendManager precompile (0x1003) - Send native Cosmos coins from EVM

✅ **Registered (via DefaultStaticPrecompiles)**:
- All default Cosmos precompiles (0x0800-0x0806)
  - Staking (0x0800), Distribution (0x0801), ICS20/IBC (0x0802)
  - Vesting (0x0803), Bank (0x0804), Governance (0x0805), Slashing (0x0806)

📚 **Documentation**:
- Tokenization: `contracts/docs/GETTING_STARTED.md`
- Gamm: `contracts/docs/GAMM_PRECOMPILE.md`
- SendManager: `contracts/docs/SENDMANAGER_PRECOMPILE.md`
- Bank (read-only): `contracts/docs/BANK_PRECOMPILE.md`

## Best Practices

1. **Always register custom precompiles** in `registerCustomPrecompiles()`
2. **Always add addresses** to `GetAllPrecompileAddresses()` and `GetAllCustomPrecompileAddresses()`
3. **Update config.yml** with new precompile addresses
4. **Update upgrade handlers** for existing chains
5. **Use helper functions** in tests for consistency
6. **Validate addresses** don't collide (automatic on startup)
7. **Document addresses** in code comments and config files

## Troubleshooting

### Precompile Not Callable

1. **Check registration**: Verify precompile is registered in `app/evm.go`
2. **Check enablement**: Verify address is in `active_static_precompiles` in genesis
3. **Check address**: Verify the address matches between registration and enablement
4. **Check logs**: Look for "Precompile enable attempt" messages in logs

### Address Collision Error

1. **Check all precompile addresses**: Use `GetAllCustomPrecompileAddresses()`
2. **Verify uniqueness**: Each address must be unique
3. **Check default precompiles**: Ensure custom addresses don't conflict with 0x0001-0x0009 range

