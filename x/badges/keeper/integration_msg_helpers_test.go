package keeper_test

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

/* Query helpers */

func GetBadge(suite *TestSuite, ctx context.Context, id uint64) (types.BitBadge, error) {
	res, err := suite.app.BadgesKeeper.GetBadge(ctx, &types.QueryGetBadgeRequest{Id: uint64(id)})
	if err != nil {
		return types.BitBadge{}, err
	}

	return *res.Badge, nil
}

func GetUserBalance(suite *TestSuite, ctx context.Context, badgeId uint64, address uint64) (types.UserBalanceInfo, error) {
	res, err := suite.app.BadgesKeeper.GetBalance(ctx, &types.QueryGetBalanceRequest{
		BadgeId: uint64(badgeId),
		Address: uint64(address),
	})
	if err != nil {
		return types.UserBalanceInfo{}, err
	}

	return *res.BalanceInfo, nil
}

/* Msg helpers */

type BadgesToCreate struct {
	Badge   types.MsgNewBadge
	Amount  uint64
	Creator string
}

func CreateBadges(suite *TestSuite, ctx context.Context, badgesToCreate []BadgesToCreate) error {
	for _, badgeToCreate := range badgesToCreate {
		for i := 0; i < int(badgeToCreate.Amount); i++ {
			msg := types.NewMsgNewBadge(badgeToCreate.Creator, badgeToCreate.Badge.Standard, badgeToCreate.Badge.DefaultSubassetSupply, badgeToCreate.Badge.SubassetSupplysAndAmounts, badgeToCreate.Badge.Uri, badgeToCreate.Badge.Permissions, badgeToCreate.Badge.FreezeAddressRanges, badgeToCreate.Badge.ArbitraryBytes, badgeToCreate.Badge.WhitelistedRecipients)
			_, err := suite.msgServer.NewBadge(ctx, msg)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func CreateSubBadges(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, supplysAndAmounts []*types.SubassetSupplyAndAmount) error {
	msg := types.NewMsgNewSubBadge(creator, badgeId, supplysAndAmounts)
	_, err := suite.msgServer.NewSubBadge(ctx, msg)
	return err
}

func RequestTransferBadge(suite *TestSuite, ctx context.Context, creator string, from uint64, amount uint64, badgeId uint64, subbadgeRange []*types.IdRange, expirationTime uint64, cantCancelBeforeTime uint64) error {
	msg := types.NewMsgRequestTransferBadge(creator, from, amount, badgeId, subbadgeRange, expirationTime, cantCancelBeforeTime)
	_, err := suite.msgServer.RequestTransferBadge(ctx, msg)
	return err
}

func RevokeBadges(suite *TestSuite, ctx context.Context, creator string, addresses []uint64, amounts []uint64, badgeId uint64, subbadgeRange []*types.IdRange) error {
	msg := types.NewMsgRevokeBadge(creator, addresses, amounts, badgeId, subbadgeRange)
	_, err := suite.msgServer.RevokeBadge(ctx, msg)
	return err
}

func TransferBadge(suite *TestSuite, ctx context.Context, creator string, from uint64, to []uint64, amounts []uint64, badgeId uint64, subbadgeRange []*types.IdRange, expirationTime uint64, cantCancelBeforeTime uint64) error {
	msg := types.NewMsgTransferBadge(creator, from, to, amounts, badgeId, subbadgeRange, expirationTime, cantCancelBeforeTime)
	_, err := suite.msgServer.TransferBadge(ctx, msg)
	return err
}

func SetApproval(suite *TestSuite, ctx context.Context, creator string, amount uint64, address uint64, badgeId uint64, subbadgeRange []*types.IdRange) error {
	msg := types.NewMsgSetApproval(creator, amount, address, badgeId, subbadgeRange)
	_, err := suite.msgServer.SetApproval(ctx, msg)
	return err
}

func HandlePendingTransfers(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, nonceRanges []*types.IdRange, actions []uint64) error {
	msg := types.NewMsgHandlePendingTransfer(creator, badgeId, nonceRanges, actions)
	_, err := suite.msgServer.HandlePendingTransfer(ctx, msg)
	return err
}

func FreezeAddresses(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, add bool, addresses []*types.IdRange) error {
	msg := types.NewMsgFreezeAddress(creator, badgeId, add, addresses)
	_, err := suite.msgServer.FreezeAddress(ctx, msg)
	return err
}

func RequestTransferManager(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, add bool) error {
	msg := types.NewMsgRequestTransferManager(creator, badgeId, add)
	_, err := suite.msgServer.RequestTransferManager(ctx, msg)
	return err
}

func TransferManager(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, address uint64) error {
	msg := types.NewMsgTransferManager(creator, badgeId, address)
	_, err := suite.msgServer.TransferManager(ctx, msg)
	return err
}

func UpdateURIs(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, uri *types.UriObject) error {
	msg := types.NewMsgUpdateUris(creator, badgeId, uri)
	_, err := suite.msgServer.UpdateUris(ctx, msg)
	return err
}

func UpdatePermissions(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, permissions uint64) error {
	msg := types.NewMsgUpdatePermissions(creator, badgeId, permissions)
	_, err := suite.msgServer.UpdatePermissions(ctx, msg)
	return err
}

func UpdateBytes(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, bytes string) error {
	msg := types.NewMsgUpdateBytes(creator, badgeId, bytes)
	_, err := suite.msgServer.UpdateBytes(ctx, msg)
	return err
}

func SelfDestructBadge(suite *TestSuite, ctx context.Context, creator string, badgeId uint64) error {
	msg := types.NewMsgSelfDestructBadge(creator, badgeId)
	_, err := suite.msgServer.SelfDestructBadge(ctx, msg)
	return err
}

func PruneBalances(suite *TestSuite, ctx context.Context, creator string, addresses []uint64, badgeIds []uint64) error {
	msg := types.NewMsgPruneBalances(creator, badgeIds, addresses)
	_, err := suite.msgServer.PruneBalances(ctx, msg)
	return err
}

func RegisterAddresses(suite *TestSuite, ctx context.Context, creator string, addresses []string) error {
	msg := types.NewMsgRegisterAddresses(creator, addresses)
	_, err := suite.msgServer.RegisterAddresses(ctx, msg)
	return err
}