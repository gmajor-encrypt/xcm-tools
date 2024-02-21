package tx

import (
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/shopspring/decimal"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 	 5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY
const AliceSeed = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"
const DestAccountId = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"

func initClient(endpoint string) *Client {
	client := NewClient(endpoint)
	client.Conn.SetKeyRing(keyring.New(keyring.Sr25519Type, AliceSeed))
	return client
}

var transferAmount = decimal.New(1, 10)

func TestXcmSend(t *testing.T) {
	client := initClient("wss://rococo-asset-hub-rpc.polkadot.io")
	defer client.Close()

	t.Run("Test_XCM_Ump_Send", func(t *testing.T) {
		dest := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}}
		beneficiary := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{
			X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: Enum{"Any": "NULL"}, Id: DestAccountId}}},
			Parents: 0,
		}}
		amount := transferAmount
		assets := MultiAssets{V2: []V2MultiAssets{{
			Id:  AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}},
			Fun: AssetsFun{Fungible: &amount},
		}}}
		var unlimited string
		weight := Weight{Unlimited: &unlimited}
		callName, args := client.Ump.LimitedTeleportAssets(&dest, &beneficiary, &assets, 0, &weight)

		signed, err := client.Conn.SignTransaction(client.Ump.GetModuleName(), callName, args...)
		assert.NoError(t, err)
		_, err = client.Conn.SendAuthorSubmitAndWatchExtrinsic(signed)
		assert.NoError(t, err)
	})

	t.Run("Test_XCM_HRMP_Send", func(t *testing.T) {
		var destParaId uint32 = 2087
		dest := VersionedMultiLocation{V3: &V1MultiLocation{Interior: V0MultiLocation{X1: &XCMJunction{Parachain: &destParaId}}, Parents: 1}}

		beneficiary := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{
			X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: Enum{"Any": "NULL"}, Id: DestAccountId}},
		}, Parents: 0}}
		amount := transferAmount
		assets := MultiAssets{V2: []V2MultiAssets{{
			Id:  AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}},
			Fun: AssetsFun{Fungible: &amount},
		}}}
		var unlimited string
		weight := Weight{Unlimited: &unlimited}
		callName, args := client.Hrmp.LimitedReserveTransferAssets(&dest, &beneficiary, &assets, 0, &weight)

		signed, err := client.Conn.SignTransaction(client.Hrmp.GetModuleName(), callName, args...)
		assert.NoError(t, err)
		_, err = client.Conn.SendAuthorSubmitAndWatchExtrinsic(signed)
		assert.NoError(t, err)
	})
}

func Test_XCM_Dmp_Send(t *testing.T) {
	client := initClient("wss://rococo-rpc.polkadot.io")
	defer client.Close()
	var destParaId uint32 = 1000
	dest := VersionedMultiLocation{V3: &V1MultiLocation{Interior: V0MultiLocation{X1: &XCMJunction{Parachain: &destParaId}}, Parents: 0}}

	beneficiary := VersionedMultiLocation{V3: &V1MultiLocation{
		Interior: V0MultiLocation{X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: map[string]string{"Rococo": ""}, Id: DestAccountId}}},
		Parents:  0,
	}}
	amount := transferAmount
	assets := MultiAssets{V2: []V2MultiAssets{{Id: AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 0}}, Fun: AssetsFun{Fungible: &amount}}}}
	var unlimited string
	weight := Weight{Unlimited: &unlimited}
	callName, args := client.Dmp.LimitedTeleportAssets(&dest, &beneficiary, &assets, 0, &weight)
	signed, err := client.Conn.SignTransaction(client.Dmp.GetModuleName(), callName, args...)
	assert.NoError(t, err)
	_, err = client.Conn.SendAuthorSubmitAndWatchExtrinsic(signed)
	assert.NoError(t, err)
}
