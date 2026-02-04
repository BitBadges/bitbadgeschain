# IBC v2 (Async Packets) Implementation Guide

## Overview

IBC v2 is a newer protocol version that supports asynchronous packet handling. Unlike IBC Classic (v1), IBC v2 allows contracts to write acknowledgements asynchronously after packet receipt.

## Current Status

Currently, `ChannelKeeperV2` is implemented as a no-op adapter that returns an error. This is sufficient for IBC Classic operations but prevents IBC v2 functionality.

## Implementation Requirements

To support IBC v2, you need to:

### 1. Understand IBC v2 Architecture

**Key Differences from IBC Classic:**
- IBC v2 uses `clientID` and `sequence` instead of `portID`, `channelID`, and `sequence`
- Packets can be acknowledged asynchronously (after `OnRecvPacket` completes)
- Requires packet tracking infrastructure to map `clientID/sequence` → `PacketI`

### 2. Implement Packet Tracking

IBC v2 requires storing packets when they're received so they can be retrieved later by `clientID/sequence`:

```go
// Store IBC v2 packet when received
func (k *Keeper) StoreIBCV2Packet(ctx sdk.Context, clientID string, sequence uint64, packet channeltypes.Packet) {
    // Store packet indexed by clientID/sequence
    // This allows WriteAcknowledgement to retrieve it later
}
```

### 3. Implement ChannelKeeperV2.WriteAcknowledgement

The current implementation needs to:

1. **Retrieve the packet** from `clientID/sequence`:
   ```go
   packet, err := retrievePacketByClientIDAndSequence(ctx, clientID, sequence)
   if err != nil {
       return err
   }
   ```

2. **Convert `channeltypesv2.Acknowledgement` to `ibcexported.Acknowledgement`**:
   ```go
   // channeltypesv2.Acknowledgement has AppAcknowledgements field
   // Need to convert to ibcexported.Acknowledgement interface
   standardAck := convertV2AckToStandardAck(v2Ack)
   ```

3. **Call IBC v10's WriteAcknowledgement**:
   ```go
   return channelKeeper.WriteAcknowledgement(ctx, packet, standardAck)
   ```

### 4. Integration Points

**Where packets are received:**
- `OnRecvPacket` in IBC modules needs to detect IBC v2 packets
- Store packets for async acknowledgement if needed

**Where acknowledgements are written:**
- `ChannelKeeperV2.WriteAcknowledgement` is called by wasm contracts
- Must retrieve stored packet and write acknowledgement

### 5. Example Implementation

```go
type channelKeeperV2Adapter struct {
    channelKeeper *channelkeeper.Keeper
    wasmKeeper    *wasmkeeper.Keeper  // Access to async packet storage
}

func (a *channelKeeperV2Adapter) WriteAcknowledgement(
    ctx sdk.Context,
    clientID string,
    sequence uint64,
    v2Ack channeltypesv2.Acknowledgement,
) error {
    // Step 1: Retrieve packet from wasmd's async packet storage
    // wasmd stores async packets using portID/channelID/sequence
    // But IBC v2 provides clientID/sequence
    // Need to map clientID → portID/channelID
    
    // Option A: Use wasmd's LoadAsyncAckPacket if it supports clientID lookup
    // Option B: Maintain your own mapping from clientID/sequence → packet
    
    // For now, this is the challenge - IBC v10 doesn't provide
    // a direct way to get packet from clientID/sequence
    
    // Step 2: Convert v2 acknowledgement to standard acknowledgement
    if len(v2Ack.AppAcknowledgements) != 1 {
        return errorsmod.Wrapf(
            channeltypes.ErrInvalidAcknowledgement,
            "IBC v2 acknowledgement must have exactly one app acknowledgement",
        )
    }
    
    // Create standard acknowledgement
    standardAck := channeltypes.NewResultAcknowledgement(v2Ack.AppAcknowledgements[0])
    
    // Step 3: Retrieve packet (this is the missing piece)
    // Need to implement packet retrieval from clientID/sequence
    packet, err := a.retrievePacketByClientID(ctx, clientID, sequence)
    if err != nil {
        return err
    }
    
    // Step 4: Write acknowledgement using IBC v10's standard method
    return a.channelKeeper.WriteAcknowledgement(ctx, packet, standardAck)
}
```

### 6. Challenges

**Main Challenge: Packet Retrieval**
- IBC v10's `ChannelKeeper` doesn't provide a method to get packet from `clientID/sequence`
- IBC v2 uses `clientID` but IBC Classic uses `portID/channelID`
- Need to maintain a mapping: `clientID/sequence` → `PacketI`

**Possible Solutions:**

1. **Use wasmd's async packet storage:**
   - wasmd already stores async packets
   - But it uses `portID/channelID/sequence`, not `clientID/sequence`
   - Need to map `clientID` → `portID/channelID`

2. **Maintain your own packet tracking:**
   - Store packets when received in IBC v2 channels
   - Index by `clientID/sequence`
   - Retrieve when `WriteAcknowledgement` is called

3. **Query IBC v2 packet storage:**
   - IBC v10 might have v2-specific storage
   - Check `modules/core/04-channel/v2/` for packet storage methods

### 7. IBC v2 Channel Setup

IBC v2 channels are created differently than IBC Classic:
- Use IBC v2 channel handshake
- Channels are identified by `clientID` instead of `portID/channelID`
- Requires IBC v2-compatible relayer

### 8. Testing

When implementing IBC v2 support:
1. Test async packet receipt
2. Test async acknowledgement writing
3. Test packet retrieval from `clientID/sequence`
4. Test error cases (packet not found, invalid acknowledgement)

## References

- [IBC v2 Protocol Specification](https://github.com/cosmos/ibc/tree/main/spec/IBC_V2)
- [IBC-Go v10 Documentation](https://docs.cosmos.network/ibc/v10.1.x/intro)
- [wasmd IBC v2 Handler](https://github.com/CosmWasm/wasmd/blob/main/x/wasm/keeper/handler_plugin.go)

## Current Recommendation

**For now:** Keep the no-op adapter that returns an error. IBC v2 is still experimental and requires:
1. IBC v2-compatible relayers
2. Additional packet tracking infrastructure
3. Testing and validation

**When to implement:**
- When you have a specific use case requiring IBC v2
- When IBC v2 tooling and relayers are mature
- When you can test end-to-end IBC v2 flows

