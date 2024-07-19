package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestZKPs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
	}
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(100),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ZkProofs = []*types.ZkProof{
		{
			VerificationKey: `{
				"protocol": "groth16",
				"curve": "bn128",
				"nPublic": 1,
				"vk_alpha_1": [
				 "13148620221535972250452852483621388707005895520407749254029942589019911968317",
				 "2889632388761537408217987365790126432839324967233397666357486962226656274862",
				 "1"
				],
				"vk_beta_2": [
				 [
					"235184211515192335800888829014779038630276494620926581575684953647609832787",
					"19688718941701106378739125377381200373190897484213597668400002191692567452234"
				 ],
				 [
					"17139505730823840138430760175182215851804207584515613019250263963726263545682",
					"15827693620421453696738510788300059898248795618831839628175348577712566794715"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_gamma_2": [
				 [
					"10857046999023057135944570762232829481370756359578518086990519993285655852781",
					"11559732032986387107991004021392285783925812861821192530917403151452391805634"
				 ],
				 [
					"8495653923123431417604973247489272438418190587263600148770280649306958101930",
					"4082367875863433681332203403145435568316851327593401208105741076214120093531"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_delta_2": [
				 [
					"503765186455257630293828867624988330153366884012427107660916834216629962594",
					"5063908373717690658921653190209557568433126368699956121277602772003448036384"
				 ],
				 [
					"7089519436043338270412729805469570497303116366059156310430246336448604006110",
					"4176347866589911546633049464699627263382512831672401904553239775361215104922"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_alphabeta_12": [
				 [
					[
					 "18604433221051551034480587482664214870599340609016463817742300813893366628865",
					 "18085175124799896104522525897296636509228330955796915514958149085721253919374"
					],
					[
					 "5818719788508620172684530200840038363514092465014332706277430885575255653155",
					 "13271739325997130171594168052085563392477243703137959518318871241167867548541"
					],
					[
					 "15990090000399229437316874893245659161037348656905875554257382661511921214897",
					 "21037562072546073120182607002350315755086517195889115785617784603842266232387"
					]
				 ],
				 [
					[
					 "2203738935953083333113402844434965170714091346913547128007885238438368497527",
					 "12254441210030526054660165902009315439490317537751453916247356602839811164652"
					],
					[
					 "18042350999254486153299086364647673457109976987505413008188358179267817106631",
					 "20440311422039242859067644437230958996502643417500127952743189403997461532235"
					],
					[
					 "3918110739729735412812703372321634476922749810769795477323113295317098326469",
					 "659595211998342112032231124403258238789761295830862227926617204748855071177"
					]
				 ]
				],
				"IC": [
				 [
					"12090554620905913999269511296228069288923128009598574755312442527252764243545",
					"20477553370846440301046582608494182925724854012349177770405179934985848699502",
					"1"
				 ],
				 [
					"21363409161656661765284862807955826021689494498345713938600256958010467481555",
					"6970854698956478603764888286629038914534847038473858754877750290063128523222",
					"1"
				 ]
				]
			 }`,
			Uri:        "",
			CustomData: "",
		},
	}
	collectionsToCreate[0].Transfers = []*types.Transfer{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				ZkProofSolutions: []*types.ZkProofSolution{
					{
						Proof: `{
							"pi_a": [
							 "15280415993796597666813636493043784530691061101504618964413605220785439881335",
							 "675563788996052368699527111569929717352978256610453492975963571932469065914",
							 "1"
							],
							"pi_b": [
							 [
								"13271940717404841760534790023945884845854169265312102895125858252700907827176",
								"591918590265443693515148586084694063044314117360908071728035775984509485618"
							 ],
							 [
								"1321184000365240114236816479208164217123473043443213174227113077515267780427",
								"9588815572669609102065114656724798374518140131116623313758405358680261499082"
							 ],
							 [
								"1",
								"0"
							 ]
							],
							"pi_c": [
							 "20270283025961949711733235392141938263044103713056685530615033679931520723664",
							 "8871872070133185879918010372399365027368084761154313139928014267382902159344",
							 "1"
							],
							"protocol": "groth16",
							"curve": "bn128"
						 }`,
						PublicInputs: `[
							"33"
						 ]`,
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				ZkProofSolutions: []*types.ZkProofSolution{
					{
						Proof: `{
							"pi_a": [
							 "15280415993796597666813636493043784530691061101504618964413605220785439881335",
							 "675563788996052368699527111569929717352978256610453492975963571932469065914",
							 "1"
							],
							"pi_b": [
							 [
								"13271940717404841760534790023945884845854169265312102895125858252700907827176",
								"591918590265443693515148586084694063044314117360908071728035775984509485618"
							 ],
							 [
								"1321184000365240114236816479208164217123473043443213174227113077515267780427",
								"9588815572669609102065114656724798374518140131116623313758405358680261499082"
							 ],
							 [
								"1",
								"0"
							 ]
							],
							"pi_c": [
							 "20270283025961949711733235392141938263044103713056685530615033679931520723664",
							 "8871872070133185879918010372399365027368084761154313139928014267382902159344",
							 "1"
							],
							"protocol": "groth16",
							"curve": "bn128"
						 }`,
						PublicInputs: `[
							"33"
						 ]`,
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge") //should be marked as used
}

func (suite *TestSuite) TestZKPsInvalidVerificationKey() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
	}
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(100),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ZkProofs = []*types.ZkProof{
		{
			//Note abc123 at vk_alpha_1
			VerificationKey: `{
				"protocol": "groth16",
				"curve": "bn128",
				"nPublic": 1,
				"vk_alpha_1": [
				 "13148620221535972250452852483621388707005895520407749254029942589019911968317abc123",
				 "2889632388761537408217987365790126432839324967233397666357486962226656274862",
				 "1"
				],
				"vk_beta_2": [
				 [
					"235184211515192335800888829014779038630276494620926581575684953647609832787",
					"19688718941701106378739125377381200373190897484213597668400002191692567452234"
				 ],
				 [
					"17139505730823840138430760175182215851804207584515613019250263963726263545682",
					"15827693620421453696738510788300059898248795618831839628175348577712566794715"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_gamma_2": [
				 [
					"10857046999023057135944570762232829481370756359578518086990519993285655852781",
					"11559732032986387107991004021392285783925812861821192530917403151452391805634"
				 ],
				 [
					"8495653923123431417604973247489272438418190587263600148770280649306958101930",
					"4082367875863433681332203403145435568316851327593401208105741076214120093531"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_delta_2": [
				 [
					"503765186455257630293828867624988330153366884012427107660916834216629962594",
					"5063908373717690658921653190209557568433126368699956121277602772003448036384"
				 ],
				 [
					"7089519436043338270412729805469570497303116366059156310430246336448604006110",
					"4176347866589911546633049464699627263382512831672401904553239775361215104922"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_alphabeta_12": [
				 [
					[
					 "18604433221051551034480587482664214870599340609016463817742300813893366628865",
					 "18085175124799896104522525897296636509228330955796915514958149085721253919374"
					],
					[
					 "5818719788508620172684530200840038363514092465014332706277430885575255653155",
					 "13271739325997130171594168052085563392477243703137959518318871241167867548541"
					],
					[
					 "15990090000399229437316874893245659161037348656905875554257382661511921214897",
					 "21037562072546073120182607002350315755086517195889115785617784603842266232387"
					]
				 ],
				 [
					[
					 "2203738935953083333113402844434965170714091346913547128007885238438368497527",
					 "12254441210030526054660165902009315439490317537751453916247356602839811164652"
					],
					[
					 "18042350999254486153299086364647673457109976987505413008188358179267817106631",
					 "20440311422039242859067644437230958996502643417500127952743189403997461532235"
					],
					[
					 "3918110739729735412812703372321634476922749810769795477323113295317098326469",
					 "659595211998342112032231124403258238789761295830862227926617204748855071177"
					]
				 ]
				],
				"IC": [
				 [
					"12090554620905913999269511296228069288923128009598574755312442527252764243545",
					"20477553370846440301046582608494182925724854012349177770405179934985848699502",
					"1"
				 ],
				 [
					"21363409161656661765284862807955826021689494498345713938600256958010467481555",
					"6970854698956478603764888286629038914534847038473858754877750290063128523222",
					"1"
				 ]
				]
			 }`,
			Uri:        "",
			CustomData: "",
		},
	}
	collectionsToCreate[0].Transfers = []*types.Transfer{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				ZkProofSolutions: []*types.ZkProofSolution{
					{
						Proof: `{
							"pi_a": [
							 "15280415993796597666813636493043784530691061101504618964413605220785439881335",
							 "675563788996052368699527111569929717352978256610453492975963571932469065914",
							 "1"
							],
							"pi_b": [
							 [
								"13271940717404841760534790023945884845854169265312102895125858252700907827176",
								"591918590265443693515148586084694063044314117360908071728035775984509485618"
							 ],
							 [
								"1321184000365240114236816479208164217123473043443213174227113077515267780427",
								"9588815572669609102065114656724798374518140131116623313758405358680261499082"
							 ],
							 [
								"1",
								"0"
							 ]
							],
							"pi_c": [
							 "20270283025961949711733235392141938263044103713056685530615033679931520723664",
							 "8871872070133185879918010372399365027368084761154313139928014267382902159344",
							 "1"
							],
							"protocol": "groth16",
							"curve": "bn128"
						 }`,
						PublicInputs: `[
							"33"
						 ]`,
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge")
}

func (suite *TestSuite) TestZKPsInvalidProof() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
	}
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(100),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ZkProofs = []*types.ZkProof{
		{
			VerificationKey: `{
				"protocol": "groth16",
				"curve": "bn128",
				"nPublic": 1,
				"vk_alpha_1": [
				 "13148620221535972250452852483621388707005895520407749254029942589019911968317",
				 "2889632388761537408217987365790126432839324967233397666357486962226656274862",
				 "1"
				],
				"vk_beta_2": [
				 [
					"235184211515192335800888829014779038630276494620926581575684953647609832787",
					"19688718941701106378739125377381200373190897484213597668400002191692567452234"
				 ],
				 [
					"17139505730823840138430760175182215851804207584515613019250263963726263545682",
					"15827693620421453696738510788300059898248795618831839628175348577712566794715"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_gamma_2": [
				 [
					"10857046999023057135944570762232829481370756359578518086990519993285655852781",
					"11559732032986387107991004021392285783925812861821192530917403151452391805634"
				 ],
				 [
					"8495653923123431417604973247489272438418190587263600148770280649306958101930",
					"4082367875863433681332203403145435568316851327593401208105741076214120093531"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_delta_2": [
				 [
					"503765186455257630293828867624988330153366884012427107660916834216629962594",
					"5063908373717690658921653190209557568433126368699956121277602772003448036384"
				 ],
				 [
					"7089519436043338270412729805469570497303116366059156310430246336448604006110",
					"4176347866589911546633049464699627263382512831672401904553239775361215104922"
				 ],
				 [
					"1",
					"0"
				 ]
				],
				"vk_alphabeta_12": [
				 [
					[
					 "18604433221051551034480587482664214870599340609016463817742300813893366628865",
					 "18085175124799896104522525897296636509228330955796915514958149085721253919374"
					],
					[
					 "5818719788508620172684530200840038363514092465014332706277430885575255653155",
					 "13271739325997130171594168052085563392477243703137959518318871241167867548541"
					],
					[
					 "15990090000399229437316874893245659161037348656905875554257382661511921214897",
					 "21037562072546073120182607002350315755086517195889115785617784603842266232387"
					]
				 ],
				 [
					[
					 "2203738935953083333113402844434965170714091346913547128007885238438368497527",
					 "12254441210030526054660165902009315439490317537751453916247356602839811164652"
					],
					[
					 "18042350999254486153299086364647673457109976987505413008188358179267817106631",
					 "20440311422039242859067644437230958996502643417500127952743189403997461532235"
					],
					[
					 "3918110739729735412812703372321634476922749810769795477323113295317098326469",
					 "659595211998342112032231124403258238789761295830862227926617204748855071177"
					]
				 ]
				],
				"IC": [
				 [
					"12090554620905913999269511296228069288923128009598574755312442527252764243545",
					"20477553370846440301046582608494182925724854012349177770405179934985848699502",
					"1"
				 ],
				 [
					"21363409161656661765284862807955826021689494498345713938600256958010467481555",
					"6970854698956478603764888286629038914534847038473858754877750290063128523222",
					"1"
				 ]
				]
			 }`,
			Uri:        "",
			CustomData: "",
		},
	}
	collectionsToCreate[0].Transfers = []*types.Transfer{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				ZkProofSolutions: []*types.ZkProofSolution{
					{
						//Note abc123 at end pi_a
						Proof: `{
							"pi_a": [
							 "15280415993796597666813636493043784530691061101504618964413605220785439881335abc123",
							 "675563788996052368699527111569929717352978256610453492975963571932469065914",
							 "1"
							],
							"pi_b": [
							 [
								"13271940717404841760534790023945884845854169265312102895125858252700907827176",
								"591918590265443693515148586084694063044314117360908071728035775984509485618"
							 ],
							 [
								"1321184000365240114236816479208164217123473043443213174227113077515267780427",
								"9588815572669609102065114656724798374518140131116623313758405358680261499082"
							 ],
							 [
								"1",
								"0"
							 ]
							],
							"pi_c": [
							 "20270283025961949711733235392141938263044103713056685530615033679931520723664",
							 "8871872070133185879918010372399365027368084761154313139928014267382902159344",
							 "1"
							],
							"protocol": "groth16",
							"curve": "bn128"
						 }`,
						PublicInputs: `[
							"33"
						 ]`,
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge") //should be marked as used
}
