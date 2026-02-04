package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

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

// IBC v10: channelID parameter added
func (w *customHooksWrapper) OnRecvPacket(ctx sdk.Context, channelID string, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// If custom hooks exist, use them; otherwise go straight to app
	if w.customHooks != nil {
		// Create a minimal IBCMiddleware wrapper for custom hooks
		// Custom hooks expect an IBCMiddleware, so we need to provide one
		// We pass nil for ICS4Middleware since custom hooks don't use it for OnRecvPacket
		minimalIM := ibchooks.NewIBCMiddleware(w.app, nil)
		return w.customHooks.OnRecvPacketOverride(minimalIM, ctx, channelID, packet, relayer)
	}
	return w.app.OnRecvPacket(ctx, channelID, packet, relayer)
}

// Delegate all other IBCModule methods to the underlying app
// IBC v10: capabilities removed
func (w *customHooksWrapper) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, counterparty channeltypes.Counterparty, version string) (string, error) {
	return w.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, counterparty, version)
}

// IBC v10: capabilities removed
func (w *customHooksWrapper) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, counterparty channeltypes.Counterparty, counterpartyVersion string) (string, error) {
	return w.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, counterparty, counterpartyVersion)
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

// IBC v10: packetID parameter added
func (w *customHooksWrapper) OnAcknowledgementPacket(ctx sdk.Context, packetID string, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return w.app.OnAcknowledgementPacket(ctx, packetID, packet, acknowledgement, relayer)
}

// IBC v10: packetID parameter added
func (w *customHooksWrapper) OnTimeoutPacket(ctx sdk.Context, packetID string, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return w.app.OnTimeoutPacket(ctx, packetID, packet, relayer)
}
