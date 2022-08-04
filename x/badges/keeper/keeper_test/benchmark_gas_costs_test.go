package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rodaine/table"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)


func (suite *TestSuite) TestGasCosts() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	for i := 0; i < 1; i++ {
		badgesToCreate := []BadgesToCreate{
			{
				Badge: types.MsgNewBadge{
					Uri:          validUri,
					Permissions:  46,
					SubassetUris: validUri,
				},
				Amount:  1,
				Creator: bob,
			},
		}
		tbl := table.New("Function Name", "Gas Consumed")


		startGas := suite.ctx.GasMeter().GasConsumed();
		CreateBadges(suite, wctx, badgesToCreate)
		endGas := suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("CreateBadge", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadge", endGas - startGas)
		

		startGas = suite.ctx.GasMeter().GasConsumed();
		err := CreateSubBadges(suite, wctx, bob, 0, []uint64{1000000}, []uint64{1})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("CreateSubBadge 1 (Supply 10000)", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{10000})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("CreateSubBadge 10000 (Supply 1)", endGas - startGas)


		startGas = suite.ctx.GasMeter().GasConsumed();
		err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{10000})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("CreateSubBadge 10000 (Supply 10000)", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = FreezeAddresses(suite, wctx, bob, []uint64{100000}, 0, 0, true)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Freeze 1 Address", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadge", endGas - startGas)

		addressesThroughTenThousand := make([]uint64, 10000)
		for j := 0; j < 10000; j++ {
			addressesThroughTenThousand[j] = uint64(j)
		}
		startGas = suite.ctx.GasMeter().GasConsumed();
		err = FreezeAddresses(suite, wctx, bob, addressesThroughTenThousand, 0, 0, true)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Freeze 10000 Addresses", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadge", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = FreezeAddresses(suite, wctx, bob, []uint64{100000}, 0, 0, false)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Unfreeze 1 Address", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadge", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = FreezeAddresses(suite, wctx, bob, addressesThroughTenThousand, 0, 0, false)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Unfreeze 10000 Address", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadge", endGas - startGas)


		startGas = suite.ctx.GasMeter().GasConsumed();
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{1}, 0, types.SubbadgeRange{Start: 0, End: 0})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Pending 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, alice, true, 0, 0, 0, 0)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("HandlePendingTransfer - 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{1}, 0, types.SubbadgeRange{Start: 0, End: 999})
		suite.Require().Nil(err, "Error transferring badge")	
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Pending 1000 Diff IDs", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, alice, true, 0, 0, 0, 999)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("HandlePendingTransfer - 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)


		startGas = suite.ctx.GasMeter().GasConsumed();
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{1}, 0, types.SubbadgeRange{Start: 0, End: 0})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Pending 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, alice, false, 0, 0, 0, 0)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("HandlePendingTransfer (Reject) - 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{1}, 0, types.SubbadgeRange{Start: 1000, End: 1999})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Pending 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, alice, false, 0, 0, 0, 999)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("HandlePendingTransfer (Reject) - 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, types.SubbadgeRange{Start: 0, End: 0})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("RequestTransferBadge - 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, bob, true, 0, 0, 0, 0)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Accept RequestTransfer - 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, types.SubbadgeRange{Start: 0, End: 999})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Request Transfer- 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, bob, true, 0, 0, 0, 999)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Accept RequestTransfer - 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)


		startGas = suite.ctx.GasMeter().GasConsumed();
		err = RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, types.SubbadgeRange{Start: 0, End: 0})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("RequestTransferBadge - 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, bob, false, 0, 0, 0, 0)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Reject RequestTransfer - 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, types.SubbadgeRange{Start: 0, End: 999})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Request Transfer- 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = HandlePendingTransfers(suite, wctx, bob, false, 0, 0, 0, 999)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Reject RequestTransfer - 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = RevokeBadges(suite, wctx, bob, []uint64{ firstAccountNumCreated + 1}, []uint64{ 1 }, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Revoke Badge 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		for j := 0; j < 1000; j++ {
			err = RevokeBadges(suite, wctx, bob, []uint64{ firstAccountNumCreated + 1}, []uint64{ 1 }, 0, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Revoke Badge 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)


		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = SetApproval(suite, wctx, bob, 10000, firstAccountNumCreated + 1, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Set Approval 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = SetApproval(suite, wctx, bob, 10, firstAccountNumCreated + 1, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Remove Approval 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)


		startGas = suite.ctx.GasMeter().GasConsumed();
		// for j := 0; j < 1000; j++ {
		// 	err = SetApproval(suite, wctx, bob, 10, firstAccountNumCreated + 1 + uint64(j), 0, 0)
		// 	suite.Require().Nil(err, "Error transferring badge")
		// }
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("Set Approval 1000", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)	
		
		badgesToCreate = []BadgesToCreate{
			{
				Badge: types.MsgNewBadge{
					Uri:          validUri,
					Permissions:  62,
					SubassetUris: validUri,
				},
				Amount:  1,
				Creator: bob,
			},
		}

		startGas = suite.ctx.GasMeter().GasConsumed();
		CreateBadges(suite, wctx, badgesToCreate)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("CreateBadge", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadge(suite, wctx, 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadge", endGas - startGas)
		

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = CreateSubBadges(suite, wctx, bob, 1, []uint64{10000}, []uint64{10000})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("CreateSubBadge 10000 (Supply 10000)", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{1}, 1, types.SubbadgeRange{Start: 0, End: 0})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Forceful 1", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{1}, 1, types.SubbadgeRange{Start: 0, End: 999})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Forceful 1000 (Same Address, Diff Subbadge ID)", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		addresses := []uint64{}
		for i := 0; i < 1000; i++ {
			addresses = append(addresses, firstAccountNumCreated + 1 + uint64(i))
		}
		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, addresses, []uint64{1}, 1, types.SubbadgeRange{Start: 0, End: 0})
		suite.Require().Nil(err, "Error transferring badge")
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Forceful 1000 (Different Addresses, Same SubbadgeId)", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)


		err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, addresses, []uint64{1}, 1, types.SubbadgeRange{Start: 0, End: 999})
		suite.Require().Nil(err, "Error transferring badge")
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("TransferBadge - Forceful 1000 (Different Addresses, Diff SubbadgeId)", endGas - startGas)

		startGas = suite.ctx.GasMeter().GasConsumed();
		_, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
		endGas = suite.ctx.GasMeter().GasConsumed();
		tbl.AddRow("GetBadgeBalance", endGas - startGas)


		firstColumnFormatter := func(format string, vals ...interface{}) string {
			return ""
		}
		
		tbl.WithFirstColumnFormatter(firstColumnFormatter).Print()

	}
}