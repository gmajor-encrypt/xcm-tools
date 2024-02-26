package tx

import "github.com/shopspring/decimal"

// SimplifyMultiLocationParaId Simplify paraId to VersionedMultiLocation
func SimplifyMultiLocationParaId(paraId uint32) *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V3: &V3MultiLocation{
			Interior: V3MultiLocationJunctions{
				X1: &XCMJunctionV3{
					Parachain: &paraId,
				},
			},
			Parents: 0,
		},
	}
}

// SimplifyMultiLocationRelayChain Simplify parent relay chain to VersionedMultiLocation
func SimplifyMultiLocationRelayChain() *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V2: &V1MultiLocation{
			Interior: V0MultiLocation{
				Here: "NULL",
			},
			Parents: 1,
		},
	}
}

// SimplifyMultiLocationAccountId32 Simplify accountId to VersionedMultiLocation
func SimplifyMultiLocationAccountId32(accountId string) *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V2: &V1MultiLocation{
			Interior: V0MultiLocation{
				X1: &XCMJunction{
					AccountId32: &XCMJunctionAccountId32{
						Network: Enum{
							"Any": "NULL",
						},
						Id: accountId,
					},
				},
			},
			Parents: 0,
		}}
}

// SimplifyMultiAssets Simplify sovereignty token to MultiAssets
func SimplifyMultiAssets(amount decimal.Decimal) *MultiAssets {
	return &MultiAssets{
		V2: []V2MultiAssets{
			{
				Id: AssetsId{
					Concrete: &V1MultiLocation{
						Interior: V0MultiLocation{Here: "NULL"},
						Parents:  1,
					},
				},
				Fun: AssetsFun{
					Fungible: &amount,
				},
			},
		},
	}
}

func SimplifyEthereumAssets(chainId uint64, tokenContract string, amount decimal.Decimal) V3MultiAssets {
	return V3MultiAssets{
		Id: AssetsIdV3{
			Concrete: &V3MultiLocation{
				Interior: V3MultiLocationJunctions{
					X2: map[string]XCMJunctionV3{
						"col0": {GlobalConsensus: &GlobalConsensusNetworkId{Ethereum: &chainId}},
						"col1": {AccountKey20: &XCMJunctionV3AccountKey20{Key: tokenContract}},
					},
				},
				Parents: 2,
			},
		},
		Fun: AssetsFun{
			Fungible: &amount,
		},
	}
}

// SimplifyUnlimitedWeight Simplify unlimited weight
func SimplifyUnlimitedWeight() *Weight {
	var unlimited string
	return &Weight{Unlimited: &unlimited}
}
