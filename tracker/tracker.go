package tracker

import (
	"context"
	"errors"
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/gmajor-encrypt/xcm-tools/util"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/hasher"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"log"
	"strings"
)

type ITracker interface {
	Track(ctx context.Context) (*Event, error)
}

var (
	NotfoundXcmMessageErr = errors.New("not found xcm message")
	InvalidExtrinsic      = errors.New("invalid extrinsic")
	InvalidDestParaId     = errors.New("invalid dest para id")
	InvalidParaHead       = errors.New("invalid para head")
)

func CreateSnapshotClient(endpoint string) (*websocket.PoolConn, *metadata.Instant, func()) {
	websocket.SetEndpoint(endpoint)
	// websocket init and set metadata
	p, err := websocket.Init()
	if err != nil {
		panic(err)
	}
	raw, err := rpc.GetMetadataByHash(nil)
	if err != nil {
		panic(err)
	}
	metadataInstant := metadata.RegNewMetadataType(0, raw)
	return p, metadataInstant, func() {
		websocket.Close()
	}
}

const (
	Parachain  = "parachain"
	Relaychain = "relaychain"
	Solo       = "solo"
)

// checkChain check chain is parachain/relaychain or not
func checkChain(m *metadata.Instant) string {
	switch {
	case tx.GetModule("XcmPallet", m) != nil:
		return Relaychain
	case tx.GetModule("PolkadotXcm", m) != nil:
		return Parachain
	default:
		return Solo
	}
}

type ExtrinsicIndex struct {
	BlockNum uint
	Index    uint
}

func findOutBlockByExtrinsicIndex(extrinsicIndex string) *ExtrinsicIndex {
	if sliceIndex := strings.Split(extrinsicIndex, "-"); len(sliceIndex) == 2 {
		return &ExtrinsicIndex{BlockNum: util.ToUint(sliceIndex[0]), Index: util.ToUint(sliceIndex[1])}
	}
	return nil
}

func hash(hex string) string {
	return utiles.AddHex(utiles.BytesToHex(hasher.HashByCryptoName(utiles.HexToBytes(hex), "Blake2_256")))
}

func TrackXcmMessage(extrinsicIndex string, protocol tx.Protocol, originEndpoint, destEndpoint, relayEndpoint string) (*Event, error) {
	log.Println("Start track xcm message with ExtrinsicIndex:", extrinsicIndex,
		"Protocol:", protocol,
		"OriginEndpoint:", originEndpoint,
		"DestEndpoint:", destEndpoint,
		"RelayEndpoint:", relayEndpoint)

	_, metadataInstant, cancel := CreateSnapshotClient(originEndpoint)
	chain := checkChain(metadataInstant)
	if chain == Solo {
		return nil, errors.New("originEndpoint not parachain or relaychain")
	}
	ctx := context.Background()
	cancel()
	defer func() {
		types.Clean()
	}()
	switch protocol {
	case tx.UMP:
		if chain != Parachain {
			return nil, errors.New("originEndpoint not parachain")
		}
		u := Ump{ExtrinsicIndex: extrinsicIndex, OriginEndpoint: originEndpoint, DestEndpoint: destEndpoint}
		return u.Track(ctx)
	case tx.HRMP:
		if chain != Parachain {
			return nil, errors.New("originEndpoint not parachain")
		}
		h := Hrmp{extrinsicIndex: extrinsicIndex, originEndpoint: originEndpoint, destEndpoint: destEndpoint, relayChainEndpoint: relayEndpoint, filterCallBack: hrmpFilter}
		return h.Track(ctx)
	case tx.DMP:
		if chain != Relaychain {
			return nil, errors.New("originEndpoint not relaychain")
		}
		d := Dmp{extrinsicIndex: extrinsicIndex, originEndpoint: originEndpoint, destEndpoint: destEndpoint}
		return d.Track(ctx)
	default:
		return nil, errors.New("not support solo chain")
	}
}

// TrackBridgeMessage ethereum <=> polkadot
type TrackBridgeMessageOptions struct {
	Tx                string
	ChainId           uint
	BridgeHubEndpoint string // bridge hub websocket endpoint
	extrinsicIndex    string
	originEndpoint    string
	relayEndpoint     string
}

func TrackBridgeMessage(ctx context.Context, opt *TrackBridgeMessageOptions) (*Event, error) {
	// ethereum -> polkadot
	// 1. 根据 tx hash 获取 receipt, 以及block header https://sepolia.etherscan.io/tx/0x799f01445e2be3103a1a751e33b395c4b894529ce3b320d2fd94c22d4e3d6e01
	// 2. 根据时间戳获取 bridge hub 的 block_num
	// 3. 根据block_num 获取 events, 直到找到 bridge message
	if opt.Tx != "" {
		receipt, err := util.EthGetTransactionReceipt(ctx, opt.Tx)
		if err != nil {
			return nil, err
		}
		blockNum := util.HexToUint64(receipt.BlockNumber)

		// Get block timestamp
		block, err := util.EthGetBlockByNum(ctx, blockNum)
		if err != nil {
			return nil, err
		}
		timestamp := util.HexToUint64(block.Timestamp)

		var messageId string
		for _, l := range receipt.Logs {
			if len(l.Topics) == 0 {
				continue
			}
			if l.Topics[0] == OutboundMessageAcceptedTopic {
				messageId = l.Topics[2]
				break
			}
		}
		if messageId == "" {
			return nil, errors.New("not found message id")
		}

		polkadotBlock, err := util.SubscanGetBlockByTime(ctx, "polkadot", uint(timestamp))
		if err != nil {
			return nil, err
		}
		startCrawlNum := polkadotBlock.BlockNum
		// check one day block events
		// 	ethereuminboundqueue (MessageReceived)
		eventReqParams := util.SubscanEventRequestParams{
			Row:     100,
			Page:    0,
			Module:  "ethereumInboundQueue",
			EventId: "MessageReceived",
			Order:   "asc",
		}
		for {
			eventReqParams.BlockRange = fmt.Sprintf("%d-0", startCrawlNum)
			events, err := util.SubscanGetEvents(ctx, "bridgehub-rococo", &eventReqParams)
			if err != nil {
				return nil, err
			}
			for _, e := range events {
				if e.Params[2].Value.(string) == messageId {
					log.Println("Find bridge message", e.EventIndex, e.ExtrinsicIndex)
					return nil, nil
				}
				startCrawlNum = int(findOutBlockByExtrinsicIndex(e.ExtrinsicIndex).BlockNum)
			}
		}
	}

	// polkadot -> ethereum
	// 同 hrmp 一样，assetHub -> bridgeHub，获取到messageId
	// 根据时间戳去获取 ethereum 的 block_num，然后获取 events，直到找到 bridge message https://sepolia.etherscan.io/tx/0x00ec2debb2c1fbca53ec2c72388edafde4279fa22a832f56c3c186244433b4d9#eventlog
	h := Hrmp{extrinsicIndex: opt.extrinsicIndex, originEndpoint: opt.originEndpoint, destEndpoint: opt.BridgeHubEndpoint, relayChainEndpoint: opt.relayEndpoint}
	filterCall := func(events []Event, i *tx.VersionedXcm, _, blockHash string) (*Event, string, error) {
		// MessageQueued(H256)
		messageId := i.PickoutExportMessageTopic()
		if event := findEventByEventId(events, 0, []string{"MessageQueued"}); event != nil {
			if len(event.Params) == 1 {
				if event.Params[0].Value.(string) == messageId {
					log.Printf("Find Message messageid %s", messageId)
					blockTimestamp, err := timestampNow(blockHash)
					if err != nil {
						return nil, "", errors.New("not found block timestamp")
					}
					event.BlockTime = blockTimestamp
					return event, messageId, nil
				}
			}
		}
		return nil, "", nil
	}
	h.filterCallBack = filterCall

	event, err := h.Track(ctx)
	if err != nil {
		return nil, err
	}
	etherStartBlockNum, err := util.EtherscanGetBlockByTime(ctx, event.BlockTime)
	if err != nil {
		return nil, err
	}
	log.Println("Get etherStartBlockNum", etherStartBlockNum)

	logs, err := util.EtherscanGetLogs(ctx, uint64(etherStartBlockNum), bridgeContract, InboundMessageDispatchedTopic, 1, 1000)
	if err != nil {
		return nil, err
	}
	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}
		if l.Topics[2] == h.messageId {
			log.Println("Find bridge message have process in", util.HexToUint64(l.BlockNumber), l.TransactionHash)
		}
	}
	return nil, nil
}
