# XCM Tools

---
XCM Tools is a tool and an SDK. This library is written in golang.
It provides the following functions: sending xcm messages, parsing xcm message instructions and getting the execution
result of the execution after sending xcm.

## Features

- [x] Send VMP(UMP & HRMP) message
- [x] Send HRMP message
- [ ] Parse xcm message
- [ ] Tracer xcm message result

## Get Start

### Requirement

1. Go 1.18+
2. docker(optional)

### Build

```bash 
docker build -t xcm-tools .
```

### Usage

#### Installation

```bash 
go get -u github.com/gmajor-encrypt/xcm-tools  
```

#### Xcm Client

```go
package example

import (
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/substrate-api-rpc/keyring"
)

func New() *tx.Client {
	endpoint := "endpoint" // need set endpoint
	client := tx.NewClient(endpoint)
	client.SetKeyRing("...") //  need set secret key
	return client
}
```

#### Send Xcm Transfer(simplified)

We provide the following methods to simplify the parameters of sending xcm transfer, so that there is no need to
construct complex **multiLocation** and **multiAssets**.

```go
package example

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/shopspring/decimal"
)

// SendUmpMessage
// Send ump message from asset-hub to rococo relay chain
func SendUmpMessage() {
	client := New()
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendUmpTransfer(beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

// SendHrmpMessage
// Send hrmp message from rococo asset-hub to picasso-rococo testnet
func SendHrmpMessage() {
	client := New()
	destParaId := 2087 // destination parachain id
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendHrmpTransfer(destParaId, beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

// SendDmpMessage 
// Send dmp message from rococo to asset-hub 
func SendDmpMessage() {
	client := New()
	destParaId := 1000 // destination parachain id, rococo asset hub
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendDmpTransfer(destParaId, beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

```

### Send Xcm Message

#### Support Methods

* `LimitedReserveTransferAssets`
* `LimitedTeleportAssets`
* `TeleportAssets`
* `ReserveTransferAssets`
* `Send`

#### Example

limited_reserve_transfer_assets, transfer assets from asset-hub to relay chain(rococo)

```go
package main

import (
	. "github.com/gmajor-encrypt/xcm-tools/tx"
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

	// send an ump message use limited_reserve_transfer_assets
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
}

```

More examples can be found in the [example](./example) or [xcm_test](./tx/xcm_test.go) directory.


### Test

```bash
go test -v ./...
```

docker

```bash
docker run -it --rm xcm-tools
```

## License

The package is available as open source under the terms of
the [Apache LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)