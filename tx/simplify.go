package tx

import "github.com/shopspring/decimal"

// SimplifyMultiLocationParaId Simplify paraId to VersionedMultiLocation
func SimplifyMultiLocationParaId(paraId uint32) *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V3: &V1MultiLocation{
			Interior: V0MultiLocation{
				X1: &XCMJunction{
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

// SimplifyUnlimitedWeight Simplify unlimited weight
func SimplifyUnlimitedWeight() *Weight {
	var unlimited string
	return &Weight{Unlimited: &unlimited}
}
