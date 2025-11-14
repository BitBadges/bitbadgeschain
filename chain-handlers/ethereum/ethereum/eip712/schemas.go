package eip712

/*
	These are used as fully populated examples to generate EIP712 types.
	This is because the EIP712 type generation code expects all values to be populated an  non-optional.

	We want to make sure the type generation includes all default values and empty values, even for optional fields.
	This is because that is what the SDK does.
*/

// getMerkleChallengeSchema returns the schema for a Merkle challenge
func getMerkleChallengeSchema() string {
	return `{
		"root": "",
		"expectedProofLength": "",
		"useCreatorAddressAsLeaf": false,
		"leafSigner": "",
		"maxUsesPerLeaf": "",
		"challengeTrackerId": "",
		"uri": "",
		"customData": ""
	}`
}

// getETHSignatureChallengeSchema returns the schema for an ETH signature challenge
func getETHSignatureChallengeSchema() string {
	return `{
		"signer": "",
		"challengeTrackerId": "",
		"uri": "",
		"customData": ""
	}`
}

// getMustOwnTokensSchema returns the schema for must own tokens criteria
func getMustOwnTokensSchema() string {
	return `{
		"collectionId": "",
		"amountRange": ` + getUintRangeSchema() + `,
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"tokenIds": [` + getUintRangeSchema() + `],
		"overrideWithCurrentTime": false,
		"mustSatisfyForAllAssets": false,
		"ownershipCheckParty": ""
	}`
}

// getCoinTransferSchema returns the schema for coin transfers
func getCoinTransferSchema() string {
	return `{
		"to": "",
		"overrideFromWithApproverAddress": false,
		"overrideToWithInitiator": false,
		"coins": [
			{
				"amount": "",
				"denom": ""
			}
		]
	}`
}

// getBalanceSchema returns the schema for a balance
func getBalanceSchema() string {
	return `{
		"amount": "",
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"tokenIds": [` + getUintRangeSchema() + `]
	}`
}

// getPredeterminedBalancesSchema returns the schema for predetermined balances
func getPredeterminedBalancesSchema() string {
	return `{
		"manualBalances": [
			{
				"balances": [
					` + getBalanceSchema() + `
				]
			}
		],
		"incrementedBalances": {
			"startBalances": [
				` + getBalanceSchema() + `
			],
			"incrementTokenIdsBy": "",
			"allowOverrideWithAnyValidToken": false,
			"durationFromTimestamp": "",
			"incrementOwnershipTimesBy": "",
			"allowOverrideTimestamp": true,
			"recurringOwnershipTimes": { 
				"startTime": "",
				"intervalLength": "",
				"chargePeriodLength": ""
			}
		},
		"orderCalculationMethod": {
			"useOverallNumTransfers": false,
			"usePerToAddressNumTransfers": false,
			"usePerFromAddressNumTransfers": false,
			"usePerInitiatedByAddressNumTransfers": false,
			"challengeTrackerId": "",
			"useMerkleChallengeLeafIndex": false
		}
	}`
}

// getApprovalAmountsSchema returns the schema for approval amounts
func getApprovalAmountsSchema() string {
	return `{
		"overallApprovalAmount": "",
		"perToAddressApprovalAmount": "",
		"perFromAddressApprovalAmount": "",
		"amountTrackerId": "",
		"perInitiatedByAddressApprovalAmount": "",
		"resetTimeIntervals": {
			"startTime": "",
			"intervalLength": ""
		}
	}`
}

// getMaxNumTransfersSchema returns the schema for max number of transfers
func getMaxNumTransfersSchema() string {
	return `{
		"overallMaxNumTransfers": "",
		"perToAddressMaxNumTransfers": "",
		"perFromAddressMaxNumTransfers": "",
		"amountTrackerId": "",
		"perInitiatedByAddressMaxNumTransfers": "",
		"resetTimeIntervals": {
			"startTime": "",
			"intervalLength": ""
		}
	}`
}

// getAutoDeletionOptionsSchema returns the schema for auto deletion options
func getAutoDeletionOptionsSchema() string {
	return `{ 
		"afterOneUse": true, 
		"afterOverallMaxNumTransfers": false, 
		"allowCounterpartyPurge": false, 
		"allowPurgeIfExpired": false 
	}`
}

// getDynamicStoreChallengeSchema returns the schema for dynamic store challenges
func getDynamicStoreChallengeSchema() string {
	return `{ "storeId": "" }`
}

// getBaseApprovalCriteriaSchema returns the base approval criteria schema without collection-specific fields
func getBaseApprovalCriteriaSchema() string {
	return `{
		"mustOwnTokens": [
			` + getMustOwnTokensSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getCollectionApprovalCriteriaSchema returns the schema for collection approval criteria
func getCollectionApprovalCriteriaSchema() string {
	return `{
		"mustOwnTokens": [
			` + getMustOwnTokensSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"requireToEqualsInitiatedBy": false,
		"requireFromEqualsInitiatedBy": false,
		"requireToDoesNotEqualInitiatedBy": false,
		"requireFromDoesNotEqualInitiatedBy": false,
		"overridesFromOutgoingApprovals": false,
		"userRoyalties": {
			"percentage": "",
			"payoutAddress": ""
		},
		"overridesToIncomingApprovals": false,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getOutgoingApprovalCriteriaSchema returns the schema for outgoing approval criteria
func getOutgoingApprovalCriteriaSchema() string {
	return `{
		"mustOwnTokens": [
			` + getMustOwnTokensSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"requireToEqualsInitiatedBy": false,
		"requireToDoesNotEqualInitiatedBy": false,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getIncomingApprovalCriteriaSchema returns the schema for incoming approval criteria
func getIncomingApprovalCriteriaSchema() string {
	return `{
		"mustOwnTokens": [
			` + getMustOwnTokensSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"requireFromEqualsInitiatedBy": false,
		"requireFromDoesNotEqualInitiatedBy": false,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getUintRangeSchema returns the schema for a uint range
func getUintRangeSchema() string {
	return `{"start": "", "end": ""}`
}

// getCollectionApprovalSchema returns the schema for collection approval
func getCollectionApprovalSchema() string {
	return `{
		"fromListId": "",
		"toListId": "",
		"initiatedByListId": "",
		"transferTimes": [` + getUintRangeSchema() + `],
		"tokenIds": [` + getUintRangeSchema() + `],
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"uri": "",
		"customData": "",
		"approvalId": "",
		"version": "0",
		"approvalCriteria": ` + getCollectionApprovalCriteriaSchema() + `
	}`
}

// getIncomingApprovalSchema returns the schema for an incoming approval object
func getIncomingApprovalSchema() string {
	return `{
		"fromListId": "",
		"initiatedByListId": "",
		"transferTimes": [` + getUintRangeSchema() + `],
		"tokenIds": [` + getUintRangeSchema() + `],
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"uri": "",
		"customData": "",
		"approvalId": "",
		"version": "0",
		"approvalCriteria": ` + getIncomingApprovalCriteriaSchema() + `
	}`
}

// getOutgoingApprovalSchema returns the schema for an outgoing approval object
func getOutgoingApprovalSchema() string {
	return `{
		"toListId": "",
		"initiatedByListId": "",
		"transferTimes": [` + getUintRangeSchema() + `],
		"tokenIds": [` + getUintRangeSchema() + `],
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"uri": "",
		"customData": "",
		"approvalId": "",
		"version": "0",
		"approvalCriteria": ` + getOutgoingApprovalCriteriaSchema() + `
	}`
}

// getCollectionPermissionsSchema returns the schema for collection permissions
func getCollectionPermissionsSchema() string {
	return `{
		"canDeleteCollection": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canArchiveCollection": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateStandards": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateCustomData": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateManager": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateCollectionMetadata": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateValidTokenIds": [{"tokenIds": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateTokenMetadata": [{"tokenIds": [` + getUintRangeSchema() + `], "timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateCollectionApprovals": [{"fromListId": "", "toListId": "", "initiatedByListId": "", "transferTimes": [` + getUintRangeSchema() + `], "tokenIds": [` + getUintRangeSchema() + `], "ownershipTimes": [` + getUintRangeSchema() + `], "approvalId": "", "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}]
	}`
}

// getUserPermissionsSchema returns the schema for user permissions
func getUserPermissionsSchema() string {
	return `{
		"canUpdateOutgoingApprovals": [{"toListId": "", "initiatedByListId": "", "transferTimes": [` + getUintRangeSchema() + `], "tokenIds": [` + getUintRangeSchema() + `], "ownershipTimes": [` + getUintRangeSchema() + `], "approvalId": "", "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateIncomingApprovals": [{"fromListId": "", "initiatedByListId": "", "transferTimes": [` + getUintRangeSchema() + `], "tokenIds": [` + getUintRangeSchema() + `], "ownershipTimes": [` + getUintRangeSchema() + `], "approvalId": "", "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateAutoApproveAllIncomingTransfers": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}]
	}`
}

// GetSchemas returns all the schemas for the EIP712 types
func GetSchemas() []string {
	schemas := make([]string, 0)

	schemas = append(schemas, `{
		"type": "maps/CreateMap",
		"value": {
			"creator": "",
			"mapId": "",
					"inheritManagerTimelineFrom": "",
		"managerTimeline": [
			{
				"manager": "",
				"timelineTimes": [`+getUintRangeSchema()+`]
			}
		],
		"updateCriteria": {
			"managerOnly": false,
			"collectionId": "",
			"creatorOnly": false,
			"firstComeFirstServe": false
		},
		"valueOptions": {
			"noDuplicates": false,
			"permanentOnceSet": false,
			"expectUint": false,
			"expectBoolean": false,
			"expectAddress": false,
			"expectUri": false
		},
		"defaultValue": "",
		"metadataTimeline": [
			{
				"metadata": {
					"uri": "",
					"customData": ""
				},
				"timelineTimes": [`+getUintRangeSchema()+`]
			}
		],
		"permissions": {
			"canUpdateMetadata": [
				{
					"timelineTimes": [`+getUintRangeSchema()+`],
					"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
					"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateManager": [
				{
					"timelineTimes": [`+getUintRangeSchema()+`],
					"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
					"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canDeleteMap": [
				{
					"timelineTimes": [`+getUintRangeSchema()+`],
					"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
					"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "maps/UpdateMap",
		"value": {
			"creator": "",
			"mapId": "",
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateMetadataTimeline": false,
			"metadataTimeline": [
				{
					"metadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updatePermissions": false,
			"permissions": {
				"canUpdateMetadata": [
					{
						"timelineTimes": [`+getUintRangeSchema()+`],
						"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
						"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [`+getUintRangeSchema()+`],
						"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
						"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
					}
				],
				"canDeleteMap": [
					{
						"timelineTimes": [`+getUintRangeSchema()+`],
						"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
						"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
					}
				]
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "maps/DeleteMap",
		"value": {
			"creator": "",
			"mapId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "maps/SetValue",
		"value": {
			"creator": "",
			"mapId": "",
			"key": "",
			"value": "",
			"options": {
				"useMostRecentCollectionId": false
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/CreateProtocol",
		"value": {
			"creator": "",
			"name": "",
			"uri": "",
			"customData": "",
			"isFrozen": false
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/DeleteProtocol",
		"value": {
			"creator": "",
			"name": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/SetCollectionForProtocol",
		"value": {
			"creator": "",
			"name": "",
			"collectionId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/UpdateProtocol",
		"value": {
			"creator": "",
			"name": "",
			"uri": "",
			"customData": "",
			"isFrozen": false
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateAddressLists",
		"value": {
			"creator": "",
			"addressLists": [
				{
					"listId": "",
					"addresses": [],
					"whitelist": false,
					"uri": "",
					"customData": "",
					"createdBy": ""
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteCollection",
		"value": {
			"creator": "",
			"collectionId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/TransferTokens",
		"value": {
			"creator": "",
			"collectionId": "",
			"transfers": [
				{
					"from": "",
					"toAddresses": [],
					"balances": [
						`+getBalanceSchema()+`
					],
					"precalculateBalancesFromApproval": {
						"approvalId": "",
						"approvalLevel": "",
						"approverAddress": "",
						"version": "0"
					},
					"merkleProofs": [
						{
							"leaf": "",
							"aunts": [
								{
									"aunt": "",
									"onRight": false
								}
							],
							"leafSignature": ""
						}
					],
					"ethSignatureProofs": [
						{
							"nonce": "",
							"signature": ""
						}
					],
					"memo": "",
					"prioritizedApprovals": [
						{
							"approvalId": "",
							"approvalLevel": "",
							"approverAddress": "",
							"version": "0"
						}
					],
					"onlyCheckPrioritizedCollectionApprovals": false,
					"onlyCheckPrioritizedIncomingApprovals": false,
					"onlyCheckPrioritizedOutgoingApprovals": false,
					"precalculationOptions": {
						"overrideTimestamp": "0",
						"tokenIdsOverride": [`+getUintRangeSchema()+`]
					},
					"affiliateAddress": "",
					"numAttempts": "0"
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UniversalUpdateCollection",
		"value": {
			"creator": "",
			"collectionId": "",
			"defaultBalances": {
				"balances": [
					`+getBalanceSchema()+`
				],
				"incomingApprovals": [
					`+getIncomingApprovalSchema()+`
				],
				"outgoingApprovals": [
					`+getOutgoingApprovalSchema()+`
				],
				"userPermissions": `+getUserPermissionsSchema()+`,
				"autoApproveSelfInitiatedIncomingTransfers": true,
				"autoApproveSelfInitiatedOutgoingTransfers": true,
				"autoApproveAllIncomingTransfers": true
			},
			"updateValidTokenIds": false,
			"validTokenIds": [`+getUintRangeSchema()+`],
			"updateCollectionPermissions": false,
			"collectionPermissions": `+getCollectionPermissionsSchema()+`,
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateTokenMetadataTimeline": false,
			"tokenMetadataTimeline": [
				{
					"tokenMetadata": [
						{
							"uri": "",
							"customData": "",
							"tokenIds": [`+getUintRangeSchema()+`]
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"mintEscrowCoinsToTransfer": [
				{
					"amount": "",
					"denom": ""
				}
			],
			"cosmosCoinWrapperPathsToAdd": [
				{
					"denom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"symbol": "",
					"denomUnits": [
						{
							"decimals": "0",
							"symbol": "",
							"isDefaultDisplay": false
						}
					],
					"allowOverrideWithAnyValidToken": false,
					"allowCosmosWrapping": false
				}
			],
			"cosmosCoinBackedPathsToAdd": [
				{
					"ibcDenom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"ibcAmount": ""
				}
			],
			"invariants": {
				"noCustomOwnershipTimes": false,
				"maxSupplyPerId": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateUserApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"updateOutgoingApprovals": false,
			"outgoingApprovals": [
				`+getOutgoingApprovalSchema()+`
			],
			"updateIncomingApprovals": false,
			"incomingApprovals": [
				`+getIncomingApprovalSchema()+`
			],
			"updateAutoApproveSelfInitiatedOutgoingTransfers": false,
			"autoApproveSelfInitiatedOutgoingTransfers": false,
			"updateAutoApproveSelfInitiatedIncomingTransfers": false,
			"autoApproveSelfInitiatedIncomingTransfers": false,
			"updateAutoApproveAllIncomingTransfers": false,
			"autoApproveAllIncomingTransfers": false,
			"updateUserPermissions": false,
			"userPermissions": `+getUserPermissionsSchema()+`
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateCollection",
		"value": {
			"creator": "",
			"defaultBalances": {
										"balances": [
					`+getBalanceSchema()+`
				],
				"incomingApprovals": [
					`+getIncomingApprovalSchema()+`
				],
				"outgoingApprovals": [
					`+getOutgoingApprovalSchema()+`
				],
				"userPermissions": `+getUserPermissionsSchema()+`,
				"autoApproveSelfInitiatedIncomingTransfers": true,
				"autoApproveSelfInitiatedOutgoingTransfers": true,
				"autoApproveAllIncomingTransfers": true
			},
			"validTokenIds": [`+getUintRangeSchema()+`],
			"collectionPermissions": `+getCollectionPermissionsSchema()+`,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"tokenMetadataTimeline": [
				{
					"tokenMetadata": [
						{
							"uri": "",
							"customData": "",
							"tokenIds": [`+getUintRangeSchema()+`]
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"mintEscrowCoinsToTransfer": [
				{
					"amount": "",
					"denom": ""
				}
			],
			"cosmosCoinWrapperPathsToAdd": [
				{
					"denom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"symbol": "",
					"denomUnits": [
						{
							"decimals": "0",
							"symbol": "",
							"isDefaultDisplay": false
						}
					],
					"allowOverrideWithAnyValidToken": false,
					"allowCosmosWrapping": false
				}
			],
			"cosmosCoinBackedPathsToAdd": [
				{
					"ibcDenom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"ibcAmount": ""
				}
			],
			"invariants": {
				"noCustomOwnershipTimes": false,
				"maxSupplyPerId": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateCollection",
		"value": {
			"creator": "",
			"collectionId": "",
			"updateValidTokenIds": false,
			"validTokenIds": [`+getUintRangeSchema()+`],
			"updateCollectionPermissions": false,
			"collectionPermissions": `+getCollectionPermissionsSchema()+`,
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateTokenMetadataTimeline": false,
			"tokenMetadataTimeline": [
				{
					"tokenMetadata": [
						{
							"uri": "",
							"customData": "",
							"tokenIds": [`+getUintRangeSchema()+`]
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"mintEscrowCoinsToTransfer": [
				{
					"amount": "",
					"denom": ""
				}
			],
			"cosmosCoinWrapperPathsToAdd": [
				{
					"denom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"symbol": "",
					"denomUnits": [
						{
							"decimals": "0",
							"symbol": "",
							"isDefaultDisplay": false
						}
					],
					"allowOverrideWithAnyValidToken": false,
					"allowCosmosWrapping": false
				}
			],
			"cosmosCoinBackedPathsToAdd": [
				{
					"ibcDenom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"ibcAmount": ""
				}
			],
			"invariants": {
				"noCustomOwnershipTimes": false,
				"maxSupplyPerId": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateDynamicStore",
		"value": {
			"creator": "",
			"defaultValue": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateDynamicStore",
		"value": {
			"creator": "",
			"storeId": "",
			"defaultValue": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteDynamicStore",
		"value": {
			"creator": "",
			"storeId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetDynamicStoreValue",
		"value": {
			"creator": "",
			"storeId": "",
			"address": "",
			"value": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/IncrementStoreValue",
		"value": {
			"creator": "",
			"storeId": "",
			"address": "",
			"amount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DecrementStoreValue",
		"value": {
			"creator": "",
			"storeId": "",
			"address": "",
			"amount": "",
			"setToZeroOnUnderflow": false
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetIncomingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approval": `+getIncomingApprovalSchema()+`
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteIncomingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approvalId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetOutgoingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approval": `+getOutgoingApprovalSchema()+`
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteOutgoingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approvalId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/PurgeApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"purgeExpired": false,
			"approverAddress": "",
			"purgeCounterpartyApprovals": false,
			"approvalsToPurge": [
				{
					"approvalId": "",
					"approvalLevel": "",
					"approverAddress": "",
					"version": "0"
				}
			]
		}
	}`)

	// UniversalUpdateCollection helper message types
	schemas = append(schemas, `{
		"type": "badges/SetValidTokenIds",
		"value": {
			"creator": "",
			"collectionId": "",
			"validTokenIds": [`+getUintRangeSchema()+`],
			"canUpdateValidTokenIds": [{"tokenIds": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetManager",
		"value": {
			"creator": "",
			"collectionId": "",
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateManager": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetCollectionMetadata",
		"value": {
			"creator": "",
			"collectionId": "",
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateCollectionMetadata": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetTokenMetadata",
		"value": {
			"creator": "",
			"collectionId": "",
			"tokenMetadataTimeline": [
				{
					"tokenMetadata": [
						{
							"uri": "",
							"customData": ""
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateTokenMetadata": [{"tokenIds": [`+getUintRangeSchema()+`], "timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetCustomData",
		"value": {
			"creator": "",
			"collectionId": "",
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateCustomData": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetStandards",
		"value": {
			"creator": "",
			"collectionId": "",
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateStandards": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetCollectionApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"canUpdateCollectionApprovals": [{"fromListId": "", "toListId": "", "initiatedByListId": "", "transferTimes": [`+getUintRangeSchema()+`], "tokenIds": [`+getUintRangeSchema()+`], "ownershipTimes": [`+getUintRangeSchema()+`], "approvalId": "", "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetIsArchived",
		"value": {
			"creator": "",
			"collectionId": "",
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canArchiveCollection": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	// GAMM transaction schemas
	schemas = append(schemas, `{
		"type": "gamm/JoinPool",
		"value": {
			"sender": "",
			"poolId": "0",
			"shareOutAmount": "",
			"tokenInMaxs": [
				{
					"amount": "",
					"denom": ""
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/ExitPool",
		"value": {
			"sender": "",
			"poolId": "0",
			"shareInAmount": "",
			"tokenOutMins": [
				{
					"amount": "",
					"denom": ""
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/SwapExactAmountIn",
		"value": {
			"sender": "",
			"routes": [
				{
					"poolId": "0",
					"tokenOutDenom": ""
				}
			],
			"tokenIn": {
				"amount": "",
				"denom": ""
			},
			"tokenOutMinAmount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/SwapExactAmountInWithIBCTransfer",
		"value": {
			"sender": "",
			"routes": [
				{
					"poolId": "0",
					"tokenOutDenom": ""
				}
			],
			"tokenIn": {
				"amount": "",
				"denom": ""
			},
			"tokenOutMinAmount": "",
			"ibcTransferInfo": {
				"sourceChannel": "",
				"receiver": "",
				"memo": "",
				"timeoutTimestamp": "0"
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/SwapExactAmountOut",
		"value": {
			"sender": "",
			"routes": [
				{
					"poolId": "0",
					"tokenInDenom": ""
				}
			],
			"tokenInMaxAmount": "",
			"tokenOut": {
				"amount": "",
				"denom": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/JoinSwapExternAmountIn",
		"value": {
			"sender": "",
			"poolId": "0",
			"tokenIn": {
				"amount": "",
				"denom": ""
			},
			"shareOutMinAmount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/JoinSwapShareAmountOut",
		"value": {
			"sender": "",
			"poolId": "0",
			"tokenInDenom": "",
			"shareOutAmount": "",
			"tokenInMaxAmount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/ExitSwapShareAmountIn",
		"value": {
			"sender": "",
			"poolId": "0",
			"tokenOutDenom": "",
			"shareInAmount": "",
			"tokenOutMinAmount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/ExitSwapExternAmountOut",
		"value": {
			"sender": "",
			"poolId": "0",
			"tokenOut": {
				"amount": "",
				"denom": ""
			},
			"shareInMaxAmount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "gamm/CreateBalancerPool",
		"value": {
			"sender": "",
			"poolParams": {
				"swapFee": "",
				"exitFee": ""
			},
			"poolAssets": [
				{
					"token": {
						"amount": "",
						"denom": ""
					},
					"weight": ""
				}
			]
		}
	}`)

	// Cosmos SDK Group module schemas
	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgCreateGroup",
		"value": {
			"admin": "",
			"members": [{"address": "", "weight": "0", "metadata": ""}],
			"metadata": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgUpdateGroupMembers",
		"value": {
			"admin": "",
			"group_id": "0",
			"member_updates": [{"address": "", "weight": "0", "metadata": ""}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgUpdateGroupAdmin",
		"value": {
			"admin": "",
			"group_id": "0",
			"new_admin": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgUpdateGroupMetadata",
		"value": {
			"admin": "",
			"group_id": "0",
			"metadata": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgCreateGroupPolicy",
		"value": {
			"admin": "",
			"group_id": "0",
			"metadata": "",
			"decision_policy": {
				"type": "/cosmos.group.v1.ThresholdDecisionPolicy",
				"value": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgUpdateGroupPolicyAdmin",
		"value": {
			"admin": "",
			"group_policy_address": "",
			"new_admin": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgCreateGroupWithPolicy",
		"value": {
			"admin": "",
			"members": [{"address": "", "weight": "0", "metadata": ""}],
			"group_metadata": "",
			"group_policy_metadata": "",
			"group_policy_as_admin": false,
			"decision_policy": {
				"type": "/cosmos.group.v1.ThresholdDecisionPolicy",
				"value": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgUpdateGroupDecisionPolicy",
		"value": {
			"admin": "",
			"group_policy_address": "",
			"decision_policy": {
				"type": "/cosmos.group.v1.ThresholdDecisionPolicy",
				"value": ""
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/MsgUpdateGroupPolicyMetadata",
		"value": {
			"admin": "",
			"group_policy_address": "",
			"metadata": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/group/MsgSubmitProposal",
		"value": {
			"group_policy_address": "",
			"proposers": [],
			"metadata": "",
			"messages": [],
			"exec": 0,
			"title": "",
			"summary": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/group/MsgWithdrawProposal",
		"value": {
			"proposal_id": "0",
			"address": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/group/MsgVote",
		"value": {
			"proposal_id": "0",
			"voter": "",
			"option": 0,
			"metadata": "",
			"exec": 0
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/group/MsgExec",
		"value": {
			"proposal_id": "0",
			"executor": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "cosmos-sdk/group/MsgLeaveGroup",
		"value": {
			"address": "",
			"group_id": "0"
		}
	}`)

	return schemas
}
