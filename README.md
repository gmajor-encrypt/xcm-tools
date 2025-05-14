# XCM Tools

---
XCM Tools is a tool and an SDK. This library is written in golang.
It provides the following functions: sending xcm messages, parsing xcm message instructions and getting the execution
result of the execution after sending xcm.

## Features

- [x] Send VMP(UMP & HRMP) message
- [x] Send HRMP message
- [x] Parse xcm message
- [x] Tracer xcm message result
- [x] Cli Support
- [x] XCM V2,V3,V4 Support
- [x] Ethereum <=> Polkadot SnowBridge support
- [x] Polkadot <=> Kusama bridge

## Get Start

### Requirement

1. Go 1.23+
2. docker(optional)

### Build

#### Build Docker Image

```bash 
docker build -f Dockerfile-build -t xcm-tools .
docker run -it xcm-tools -h
```

#### Build Binary

```bash
cd cmd && go build -o xcm-tools .
```

## Usage

### Installation

```bash 
go install github.com/gmajor-encrypt/xcm-tools/cmd@latest 
```

You can find binary file(cmd) in $GOPATH/bin

### CLI Usage

XCM tools also support sending xcm messages, parsing messages, and tracking xcm transaction results through cli
commands.

```bash
go run .  -h
```

### Commands

```
NAME:
   Xcm tools - Xcm tools

USAGE:
   cmd [global options] command [command options] [arguments...]

COMMANDS:
   send           send xcm message
   parse          parse xcm message
   tracker        tracker xcm message transaction
   trackerEthBridge  tracker snowBridge message transaction
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```

### Args

#### Send XCM Message

| Name     | Description                                                                   | Suitable            |
|----------|-------------------------------------------------------------------------------|---------------------|
| endpoint | Set substrate endpoint, only support websocket protocol, like ws:// or wss:// | ALL                 |
| dest     | Dest address                                                                  | SendXCM             |
| amount   | Send xcm transfer amount                                                      | SendXCM             |
| keyring  | Set sr25519 secret key                                                        | SendXCM             |
| paraId   | Dest para id                                                                  | SendXCM(DMP & hRMP) |

#### Send Token(SnowBridge)

| Name     | Description                                                                   | Suitable  |
|----------|-------------------------------------------------------------------------------|-----------|
| endpoint | Set substrate endpoint, only support websocket protocol, like ws:// or wss:// | ALL       |
| dest     | Dest address                                                                  | EthBridge |
| amount   | Send xcm transfer amount                                                      | EthBridge |
| keyring  | Set sr25519 secret key                                                        | EthBridge |
| paraId   | Send xcm transfer amount                                                      | EthBridge |
| chainId  | Ethereum chain id                                                             | EthBridge |
| contract | Erc20 contract address                                                        | EthBridge |

#### Send S2S Message(polkadot <=>kusama)

| Name      | Description                                                                   | Suitable |
|-----------|-------------------------------------------------------------------------------|----------|
| endpoint  | Set substrate endpoint, only support websocket protocol, like ws:// or wss:// | ALL      |
| dest      | Dest address                                                                  | s2s      |
| amount    | Send xcm transfer amount                                                      | s2s      |
| keyring   | Set sr25519 secret key                                                        | s2s      |
| paraId    | Send xcm transfer amount                                                      | s2s      |
| chainId   | Ethereum chain id                                                             | s2s      |
| destChain | dest chain (support polkadot, kusama, westend )                               | s2s      |

#### Parse Message

| Name     | Description                                                                   | Suitable |
|----------|-------------------------------------------------------------------------------|----------|
| endpoint | Set substrate endpoint, only support websocket protocol, like ws:// or wss:// | ALL      |
| message  | Parsed xcm message raw data                                                   | Parse    |

#### Tracker(HRMP,UMP,DMP)

| Name               | Description                                                                   | Suitable |
|--------------------|-------------------------------------------------------------------------------|----------|
| endpoint           | Set substrate endpoint, only support websocket protocol, like ws:// or wss:// | ALL      |
| protocol           | Xcm protocol, such as UMP,HRMP,DMP                                            | Tracker  |
| destEndpoint       | Dest chain endpoint, only support websocket protocol, like ws:// or wss://    | Tracker  |
| extrinsicIndex     | Xcm message extrinsicIndex                                                    | Tracker  |
| relaychainEndpoint | Relay chain endpoint, only support websocket protocol, like ws:// or wss://   | Tracker  |

#### Tracker SnowBridge

| Name               | Description                                                                       | Suitable         |
|--------------------|-----------------------------------------------------------------------------------|------------------|
| endpoint           | Set substrate endpoint, only support websocket protocol, like ws:// or wss://     | ALL              |
| extrinsicIndex     | Xcm message extrinsicIndex                                                        | trackerEthBridge |
| relaychainEndpoint | Relay chain endpoint, only support websocket protocol, like ws:// or wss://       | trackerEthBridge |
| hash               | Ethereum send token to polkadot transaction hash                                  | trackerEthBridge |
| bridgeHubEndpoint  | BridgeHubEndpoint endpoint, only support websocket protocol, like ws:// or wss:// | trackerEthBridge |

#### Tracker S2S Message(polkadot <=>kusama)

| Name                   | Description                                                                       | Suitable         |
|------------------------|-----------------------------------------------------------------------------------|------------------|
| endpoint               | Set substrate endpoint, only support websocket protocol, like ws:// or wss://     | ALL              |
| extrinsicIndex         | Xcm message extrinsicIndex                                                        | trackerS2SBridge |
| relaychainEndpoint     | Relay chain endpoint, only support websocket protocol, like ws:// or wss://       | trackerS2SBridge |
| bridgeHubEndpoint      | BridgeHubEndpoint endpoint, only support websocket protocol, like ws:// or wss:// | trackerS2SBridge |
| destEndpoint           | dest endpoint, only support websocket protocol, like ws:// or wss://              | trackerS2SBridge |
| destBridgeHubEndpoint  | destBridgeHub endpoint, only support websocket protocol, like ws:// or wss://     | trackerS2SBridge |
| destRelaychainEndpoint | destRelaychain endpoint, only support websocket protocol, like ws:// or wss://    | trackerS2SBridge |

#### Example

```bash
# UMP
go run . send UMP --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a --endpoint wss://rococo-asset-hub-rpc.polkadot.io
# DMP
go run . send DMP --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a --endpoint wss://rococo-rpc.polkadot.io --paraId 1000
# HRMP
go run . send HRMP --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a --endpoint wss://rococo-asset-hub-rpc.polkadot.io --paraId 2087
# Send bridge message(polkadot to ethereum)
go run . send EthBridge --dest 0x6EB228b7ab726b8B44892e8e273ACF3dcC9C0492 --amount 10  --keyring 0xc0417c253312107d808921fb1dd3b740b64e99794dca74bcc550179f7c42a255 --endpoint wss://rococo-asset-hub-rpc.polkadot.io --contract 0xfff9976782d46cc05630d1f6ebab18b2324d6b14 --chainId 11155111
# Send S2S message(polkadot to kusama)
go run . send S2SBridge --paraId 1000 --destChain westend --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xc0417c253312107d808921fb1dd3b740b64e99794dca74bcc550179f7c42a255 --endpoint wss://rococo-asset-hub-rpc.polkadot.io
```

### Parse Xcm Message

We provide a function to parse xcm transaction instructions and deserialize the encoded raw message into readable JSON.
Support XCM V0,V1,V2,V3,V4

```go
package example

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/parse"
	"github.com/gmajor-encrypt/xcm-tools/tx"
)

func ParseMessage() {
	client := tx.NewClient("wss://rococo-rpc.polkadot.io")
	defer client.Close()
	p := parse.New(client.Metadata())
	instruction, err := p.ParseXcmMessageInstruction("0x031000040000000007f5c1998d2a0a130000000007f5c1998d2a000d01020400010100ea294590dbcfac4dda7acd6256078be26183d079e2739dd1e8b1ba55d94c957a")
	fmt.Println(instruction, err)
}

```

### Tracker Xcm Message

We provide a function to track xcm transaction results. Support protocol UMP,HRMP,DMP.

```go
package example

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tracker"
	"github.com/gmajor-encrypt/xcm-tools/tx"
)

// TrackerMessage Tracker UMP Message with extrinsic_index 4310901-13
func TrackerMessage() {
	event, err := tracker.TrackXcmMessage("4310901-13", tx.UMP, "wss://moonbeam-rpc.dwellir.com", "wss://polkadot-rpc.dwellir.com", "")
	fmt.Println(event, err)
}
```

### Xcm Client

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

### Send Xcm Transfer(simplified)

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
	client := tx.NewClient("endpoint")
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendUmpTransfer(beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

// SendHrmpMessage
// Send hrmp message from rococo asset-hub to picasso-rococo testnet
func SendHrmpMessage() {
	client := tx.NewClient("endpoint")
	destParaId := 2087 // destination parachain id
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendHrmpTransfer(uint32(destParaId), beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

// SendDmpMessage 
// Send dmp message from rococo to asset-hub 
func SendDmpMessage() {
	client := tx.NewClient("endpoint")
	destParaId := 1000 // destination parachain id, rococo asset hub
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendDmpTransfer(uint32(destParaId), beneficiary, transferAmount)
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
* `transferAssets`

#### Example

limited_reserve_transfer_assets, transfer assets from asset-hub to relay chain(rococo)

```go
package main

import (
	. "github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/shopspring/decimal"
	"github.com/itering/substrate-api-rpc/hasher"
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

```

More examples can be found in the [example](./example) or [xcm_test](./tx/xcm_test.go) directory.

### SnowBridge Support

##### Tracker SnowBridge Message

We provide a function to track snowBridge transaction results(Ethereum <=> Polkadot).

```go
package main

import (
	"context"
	"github.com/gmajor-encrypt/xcm-tools/tracker"
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/shopspring/decimal"
	"log"
)

func TrackerMessage() {
	var err error
	ctx := context.Background()

	// Track Eth => Polkadot
	_, err = tracker.TrackEthBridgeMessage(ctx, &tracker.TrackBridgeMessageOptions{Tx: "0x799f01445e2be3103a1a751e33b395c4b894529ce3b320d2fd94c22d4e3d6e01"})
	if err != nil {
		log.Fatal(err)
	}

	// Track Polkadot => ETH
	_, err = tracker.TrackEthBridgeMessage(ctx, &tracker.TrackBridgeMessageOptions{
		ExtrinsicIndex:    "3879712-2",
		BridgeHubEndpoint: "wss://rococo-bridge-hub-rpc.polkadot.io",
		OriginEndpoint:    "wss://rococo-rockmine-rpc.polkadot.io",
		RelayEndpoint:     "wss://rococo-rpc.polkadot.io",
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

#### Send Erc20 Token

We provide a function to send Erc20 token by snowBridge(Polkadot => Ethereum).

```go
package main

import (
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/substrate-api-rpc/keyring"
	"strings"
	"github.com/shopspring/decimal"
)

func SendEthErc20Token() {
	client := NewClient(endpoint)
	client.Conn.SetKeyRing(keyring.New(keyring.Sr25519Type, AliceSeed))
	defer client.Close()

	destH160 := strings.ToLower("0x6EB228b7ab726b8B44892e8e273ACF3dcC9C0492") // Send to Ethereum address
	Erc20Contract := "0xfff9976782d46cc05630d1f6ebab18b2324d6b14"             // Erc20 contract address
	amount := decimal.New(1, 0)                                               // Transfer amount

	_, _ = client.SendTokenToEthereum(
		destH160,
		Erc20Contract,
		amount,
		11155111) // Ethereum chain id

}
```

### Polkadot <=> Kusama Bridge

#### Send S2S Message

```go
package main

import (
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/substrate-api-rpc/keyring"
	"log"
	"strings"
	"github.com/shopspring/decimal"
)

// Demo Send S2S Message, this demo is sending a message from rococo to westend
func Demo() {
	client := tx.NewClient("")
	client.SetKeyRing("")
	defer client.Close()

	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(1, 0)
	txHash, err := client.SendDotKsmChainToken(beneficiary, uint32(1000), tx.ConvertToGlobalConsensusNetworkId("westend"), transferAmount)
	log.Println(txHash, err)
}

```

#### Tracker S2S Message

```go
package main

import (
	"context"
	"github.com/gmajor-encrypt/xcm-tools/tracker"
)

// TrackerMessage Tracker S2S Message(rococo=>westend) with extrinsic_index 5816546-2
// https://assethub-rococo.subscan.io/extrinsic/5816546-2
func Demo() {
	ctx := context.Background()
	_, err := tracker.TrackS2sBridgeMessage(ctx, &tracker.S2STrackBridgeMessageOptions{
		ExtrinsicIndex:      "5816546-2",
		
		BridgeHubEndpoint:   "wss://rococo-bridge-hub-rpc.polkadot.io",
		OriginEndpoint:      "wss://rococo-rockmine-rpc.polkadot.io",
		OriginRelayEndpoint: "wss://rococo-rpc.polkadot.io",

		DestinationEndpoint:          "wss://westend-asset-hub-rpc.polkadot.io",
		DestinationBridgeHubEndpoint: "wss://bridge-hub-westend-rpc.dwellir.com",
		DestinationRelayEndpoint:     "wss://westend-rpc.polkadot.io",
	})
}


```

### Test

```bash
# Ensure you have in code root directory
go test -v ./... 
```

docker

```bash
docker build -t xcm-tools-test .
docker run -it --rm xcm-tools-test
```

## License

The package is available as open source under the terms of
the [Apache LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)