package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	customhooks "github.com/bitbadges/bitbadgeschain/x/custom-hooks"
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
)

// customHooksWrapper wraps the app to execute custom hooks before calling the underlying app
// It implements IBCModule interface by delegating all methods to the underlying app
type customHooksWrapper struct {
	app         porttypes.IBCModule
	customHooks *customhooks.CustomHooks
}

var _ porttypes.IBCModule = (*customHooksWrapper)(nil)

func (w *customHooksWrapper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// If custom hooks exist, use them; otherwise go straight to app
	if w.customHooks != nil {
		// Create a minimal IBCMiddleware wrapper for custom hooks
		// Custom hooks expect an IBCMiddleware, so we need to provide one
		// We pass nil for ICS4Middleware since custom hooks don't use it for OnRecvPacket
		minimalIM := ibchooks.NewIBCMiddleware(w.app, nil)
		return w.customHooks.OnRecvPacketOverride(minimalIM, ctx, packet, relayer)
	}
	return w.app.OnRecvPacket(ctx, packet, relayer)
}

// Delegate all other IBCModule methods to the underlying app
func (w *customHooksWrapper) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return w.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (w *customHooksWrapper) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (string, error) {
	return w.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

func (w *customHooksWrapper) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return w.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (w *customHooksWrapper) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return w.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (w *customHooksWrapper) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return w.app.OnChanCloseInit(ctx, portID, channelID)
}

func (w *customHooksWrapper) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return w.app.OnChanCloseConfirm(ctx, portID, channelID)
}

func (w *customHooksWrapper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return w.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (w *customHooksWrapper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return w.app.OnTimeoutPacket(ctx, packet, relayer)
}
