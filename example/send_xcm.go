package main

import (
	. "github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/hasher"
	"github.com/shopspring/decimal"
	"log"
)

func main() {
	endpoint := "endpoint" // you need set an endpoint

	client := NewClient(endpoint)
	client.SetKeyRing(".....") // you need set a sr25519 secret key

	// dest is a relay chain
	dest := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}}

	// beneficiary is a relay chain account id
	beneficiary := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{
		X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: Enum{"Any": "NULL"}, Id: "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"}}},
		Parents: 0,
	}}

	// transfer amount
	amount := decimal.New(1, 0)

	// multi assets
	assets := MultiAssets{V2: []V2MultiAssets{{
		Id:  AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}},
		Fun: AssetsFun{Fungible: &amount},
	}}}

	// weight set
	weight := Weight{Limited: &WeightLimited{ProofSize: 0, RefTime: 4000000000}}

	// send an ump message use limited_teleport_assets
	callName, args := client.Ump.LimitedTeleportAssets(&dest, &beneficiary, &assets, 0, &weight)

	// sign the extrinsic
	signed, err := client.Conn.SignTransaction(client.Ump.GetModuleName(), callName, args...)

	if err != nil {
		log.Panic(err)
	}
	_, err = client.Conn.SendAuthorSubmitAndWatchExtrinsic(signed)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("This Extrinsic has successfully sent with hash %s", utiles.AddHex(utiles.BytesToHex(hasher.HashByCryptoName(utiles.HexToBytes(signed), "Blake2_256"))))
}
