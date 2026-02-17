package tokenization

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// IBCModule implements the ICS26 interface for interchain accounts host chains
type IBCModule struct {
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the associated keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
// IBC v10: capabilities removed
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	// Require portID is the portID module is bound to
	boundPort := im.keeper.GetPort(ctx)
	if boundPort != portID {
		return "", errorsmod.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if version != types.Version {
		return "", errorsmod.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", version, types.Version)
	}

	// IBC v10: capabilities removed - no need to claim capability

	return version, nil
}

// OnChanOpenTry implements the IBCModule interface
// IBC v10: capabilities removed
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// Require portID is the portID module is bound to
	boundPort := im.keeper.GetPort(ctx)
	if boundPort != portID {
		return "", errorsmod.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if counterpartyVersion != types.Version {
		return "", errorsmod.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: got: %s, expected %s", counterpartyVersion, types.Version)
	}

	// IBC v10: capabilities removed - no need to claim or authenticate capability

	return types.Version, nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	_,
	counterpartyVersion string,
) error {
	if counterpartyVersion != types.Version {
		return errorsmod.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, types.Version)
	}
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for channels
	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
// IBC v10: channelID parameter added
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	channelID string,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// this line is used by starport scaffolding # oracle/packet/module/recv

	var modulePacketData types.TokenizationPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error()))
	}

	// Dispatch packet
	switch packet := modulePacketData.Packet.(type) {
	// ICQ: Ownership Query - verify token ownership for cross-chain queries (single tokenId/time)
	case *types.TokenizationPacketData_OwnershipQuery:
		return im.handleOwnershipQuery(ctx, modulePacket, packet.OwnershipQuery)

	// ICQ: Bulk Ownership Query - handle multiple ownership queries in one packet
	case *types.TokenizationPacketData_BulkOwnershipQuery:
		return im.handleBulkOwnershipQuery(ctx, modulePacket, packet.BulkOwnershipQuery)

	// ICQ: Full Balance Query - return complete UserBalanceStore
	case *types.TokenizationPacketData_FullBalanceQuery:
		return im.handleFullBalanceQuery(ctx, modulePacket, packet.FullBalanceQuery)

	// ICQ: Response packets should not be received here (they are sent as acknowledgements)
	case *types.TokenizationPacketData_OwnershipQueryResponse,
		*types.TokenizationPacketData_BulkOwnershipQueryResponse,
		*types.TokenizationPacketData_FullBalanceQueryResponse:
		err := errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "response packets should be handled via acknowledgement flow")
		return channeltypes.NewErrorAcknowledgement(err)

	// this line is used by starport scaffolding # ibc/packet/module/recv
	default:
		err := errorsmod.Wrapf(types.ErrUnrecognizedPacketType, "packet type: %T", packet)
		return channeltypes.NewErrorAcknowledgement(err)
	}
}

// OnAcknowledgementPacket implements the IBCModule interface
// IBC v10: channelID parameter added
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	channelID string,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}

	// this line is used by starport scaffolding # oracle/packet/module/ack

	var modulePacketData types.TokenizationPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	}

	// Dispatch packet
	switch packet := modulePacketData.Packet.(type) {
	// this line is used by starport scaffolding # ibc/packet/module/ack
	default:
		errMsg := fmt.Sprintf("unrecognized %s packet type: %T", types.ModuleName, packet)
		return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	// var eventType string
	// // Handle acknowledgement
	// switch resp := ack.Response.(type) {
	// case *channeltypes.Acknowledgement_Result:
	// 	ctx.EventManager().EmitEvent(
	// 		sdk.NewEvent(
	// 			types.AttributeKeyAck,
	// 			sdk.NewAttribute(types.AttributeKeyAckSuccess, string(resp.Result)),
	// 		),
	// 	)
	// case *channeltypes.Acknowledgement_Error:
	// 	ctx.EventManager().EmitEvent(
	// 		sdk.NewEvent(
	// 			types.AttributeKeyAck,
	// 			sdk.NewAttribute(types.AttributeKeyAckError, resp.Error),
	// 		),
	// 	)
	// }

	// return nil
}

// handleOwnershipQuery processes an ICQ ownership query and returns the response as an acknowledgement
func (im IBCModule) handleOwnershipQuery(
	ctx sdk.Context,
	packet channeltypes.Packet,
	query *types.OwnershipQueryPacket,
) ibcexported.Acknowledgement {
	// Emit event for the incoming query
	im.keeper.EmitICQRequestEvent(ctx, query)

	// Process the ownership query
	response := im.keeper.ProcessOwnershipQuery(ctx, query)

	// Emit event for the response
	im.keeper.EmitICQResponseEvent(ctx, response)

	// Marshal the response as acknowledgement data
	ackData, err := types.ModuleCdc.MarshalJSON(response)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(
			errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal ICQ response: %s", err.Error()),
		)
	}

	return channeltypes.NewResultAcknowledgement(ackData)
}

// handleBulkOwnershipQuery processes a bulk ICQ ownership query and returns responses as an acknowledgement
func (im IBCModule) handleBulkOwnershipQuery(
	ctx sdk.Context,
	packet channeltypes.Packet,
	bulk *types.BulkOwnershipQueryPacket,
) ibcexported.Acknowledgement {
	// Emit event for each query in the bulk request
	for _, query := range bulk.Queries {
		im.keeper.EmitICQRequestEvent(ctx, query)
	}

	// Process the bulk ownership query
	response := im.keeper.ProcessBulkOwnershipQuery(ctx, bulk)

	// Emit events for each response
	for _, resp := range response.Responses {
		im.keeper.EmitICQResponseEvent(ctx, resp)
	}

	// Marshal the response as acknowledgement data
	ackData, err := types.ModuleCdc.MarshalJSON(response)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(
			errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal bulk ICQ response: %s", err.Error()),
		)
	}

	return channeltypes.NewResultAcknowledgement(ackData)
}

// handleFullBalanceQuery processes an ICQ full balance query and returns the complete UserBalanceStore
func (im IBCModule) handleFullBalanceQuery(
	ctx sdk.Context,
	packet channeltypes.Packet,
	query *types.FullBalanceQueryPacket,
) ibcexported.Acknowledgement {
	// Emit event for the incoming query
	im.keeper.EmitFullBalanceQueryRequestEvent(ctx, query)

	// Process the full balance query
	response := im.keeper.ProcessFullBalanceQuery(ctx, query)

	// Emit event for the response
	im.keeper.EmitFullBalanceQueryResponseEvent(ctx, response)

	// Marshal the response as acknowledgement data
	ackData, err := types.ModuleCdc.MarshalJSON(response)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(
			errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal full balance ICQ response: %s", err.Error()),
		)
	}

	return channeltypes.NewResultAcknowledgement(ackData)
}

// OnTimeoutPacket implements the IBCModule interface
// IBC v10: channelID parameter added
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	channelID string,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var modulePacketData types.TokenizationPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	}

	// Dispatch packet
	switch packet := modulePacketData.Packet.(type) {
	// this line is used by starport scaffolding # ibc/packet/module/timeout
	default:
		errMsg := fmt.Sprintf("unrecognized %s packet type: %T", types.ModuleName, packet)
		return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
