package tx

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test ConvertMultiLocationParaId with custom account id
func Test_ConvertMultiLocationAccountId32(t *testing.T) {
	cases := []struct {
		accountId string
		expected  *VersionedMultiLocation
	}{
		{
			accountId: "da5ace7cd7253f7d2e0ce11c5fbc5c91c68665016a5df21572f11f9f94088748",
			expected: &VersionedMultiLocation{
				V2: &V1MultiLocation{Interior: V0MultiLocation{X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: Enum{"Any": "NULL"}, Id: "da5ace7cd7253f7d2e0ce11c5fbc5c91c68665016a5df21572f11f9f94088748"}}}, Parents: 0}},
		},
		{
			accountId: "e3e063eb9d79010fe24ff20fbf9e73ea8b5ec4abdbecbc4e2e4d852fd63cd982",
			expected: &VersionedMultiLocation{
				V2: &V1MultiLocation{Interior: V0MultiLocation{X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: Enum{"Any": "NULL"}, Id: "e3e063eb9d79010fe24ff20fbf9e73ea8b5ec4abdbecbc4e2e4d852fd63cd982"}}}, Parents: 0}},
		},
	}
	for _, v := range cases {
		assert.EqualValues(t, v.expected, SimplifyMultiLocationAccountId32(v.accountId))
	}
}

func Test_SimplifyMultiLocationParaId(t *testing.T) {
	cases := []struct {
		paraId   uint32
		expected *VersionedMultiLocation
	}{
		{
			paraId:   200,
			expected: &VersionedMultiLocation{V3: &V3MultiLocation{Interior: V3MultiLocationJunctions{X1: &XCMJunctionV3{Parachain: &[]uint32{200}[0]}}, Parents: 0}},
		},
		{
			paraId:   300,
			expected: &VersionedMultiLocation{V3: &V3MultiLocation{Interior: V3MultiLocationJunctions{X1: &XCMJunctionV3{Parachain: &[]uint32{300}[0]}}, Parents: 0}},
		},
	}
	for _, v := range cases {
		assert.EqualValues(t, v.expected, SimplifyMultiLocationParaId(v.paraId))
	}
}

func Test_SimplifyMultiLocationRelayChain(t *testing.T) {
	assert.EqualValues(t, SimplifyMultiLocationRelayChain(), &VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}})
}

func Test_SimplifyMultiAssets(t *testing.T) {
	// test with custom amount
	var toPtr = func(amount decimal.Decimal) *decimal.Decimal {
		return &amount
	}

	cases := []struct {
		amount   decimal.Decimal
		expected *MultiAssets
	}{
		{
			amount:   decimal.NewFromInt(100),
			expected: &MultiAssets{V2: []V2MultiAssets{{Id: AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}}, Fun: AssetsFun{Fungible: toPtr(decimal.New(100, 0))}}}},
		},
		{
			amount:   decimal.NewFromInt(20000),
			expected: &MultiAssets{V2: []V2MultiAssets{{Id: AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}}, Fun: AssetsFun{Fungible: toPtr(decimal.New(20000, 0))}}}},
		},
	}
	for _, v := range cases {
		assert.EqualValues(t, v.expected, SimplifyMultiAssets(v.amount))
	}
}
