package tracker

import (
	"context"
	"errors"
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/gmajor-encrypt/xcm-tools/util"
	"log"
	"strings"
)

type TrackBridgeMessageOptions struct {
	// ethereum => polkadot message transaction hash, like 0x270b9592600015788b279df9eab62670349d006cf3ffaf3185f84163b37b9154
	Tx string
	// ChainId is the chain id of the ethereum network
	ChainId uint

	// polkadot => ethereum message extrinsic index
	ExtrinsicIndex string
	// bridge hub websocket endpoint
	BridgeHubEndpoint string
	// origin chain websocket endpoint
	OriginEndpoint string
	// relay chain websocket endpoint
	RelayEndpoint string

	// bridge hub name
	BridgeHubName string
}

var (
	bridgeHubRococo   = "bridgehub-rococo"
	bridgeHubWestend  = "bridgehub-westend"
	bridgeHubPolkadot = "bridgehub-polkadot"
	bridgeHubKusama   = "bridgehub-kusama"
)

// Subscan api endpoint name, for example: bridgehub-rococo,bridgehub-polkadot, bridgehub-westend, bridgehub-kusama

func (t *TrackBridgeMessageOptions) bridgeHubName() string {
	var defaultBridgeHubName = bridgeHubRococo
	if t.BridgeHubName != "" {
		if !util.InSlice(t.BridgeHubName, []string{bridgeHubRococo, bridgeHubWestend, bridgeHubPolkadot, bridgeHubKusama}) {
			panic(fmt.Sprintf("bridgeHubName %s is not supported", t.BridgeHubName))
		}
		return t.BridgeHubName
	}
	return defaultBridgeHubName
}

var (
	bridgeContract = map[string]string{
		bridgeHubPolkadot: "0x27ca963C279c93801941e1eB8799c23f407d68e7", // polkadot bridge contract address
		bridgeHubRococo:   "0x5B4909cE6Ca82d2CE23BD46738953c7959E710Cd", // rococo bridge contract address
	}

	extrinsicIndexEmptyError    = fmt.Errorf("extrinsicIndex is empty")
	bridgeHubEndpointEmptyError = fmt.Errorf("bridgeHubEndpoint is empty")
	originEndpointEmptyError    = fmt.Errorf("originEndpoint is empty")
	relayEndpointEmptyError     = fmt.Errorf("relayEndpoint is empty")
	notFindBridgeMessageIdError = fmt.Errorf("not found message id")
)

// TrackEthBridgeMessage ethereum <=> polkadot
func TrackEthBridgeMessage(ctx context.Context, opt *TrackBridgeMessageOptions) (*Event, error) {
	// ethereum -> polkadot
	// if tx is not empty, will track ethereum => polkadot message
	isPolkadot := opt.bridgeHubName() == bridgeHubPolkadot

	if opt.Tx != "" {
		// Get transaction receipt
		receipt, err := util.EthGetTransactionReceipt(ctx, isPolkadot, opt.Tx)
		if err != nil {
			return nil, err
		}
		blockNum := util.HexToUint64(receipt.BlockNumber)

		// Get block timestamp
		block, err := util.EthGetBlockByNum(ctx, isPolkadot, blockNum)
		if err != nil {
			return nil, err
		}
		timestamp := util.HexToUint64(block.Timestamp)

		var messageId string
		for _, l := range receipt.Logs {
			if len(l.Topics) == 0 {
				continue
			}
			// OutboundMessageAccepted (index_topic_1 bytes32 channelID, uint64 nonce, index_topic_2 bytes32 messageID, bytes payload)
			// messageId is unique id
			if l.Topics[0] == OutboundMessageAcceptedTopic {
				messageId = l.Topics[2]
				break
			}
		}
		if messageId == "" {
			return nil, notFindBridgeMessageIdError
		}

		log.Println("Get ethereum message Id", messageId, "timestamp", timestamp, "blockNum", blockNum)

		polkadotBlock, err := util.SubscanGetBlockByTime(ctx, opt.bridgeHubName(), uint(timestamp))
		if err != nil {
			log.Printf("SubscanGetBlockByTime get err %v\n", err)
			return nil, err
		}
		startCrawlNum := polkadotBlock.BlockNum

		log.Println("Start crawl block num", startCrawlNum)
		// check one day block events
		// 	ethereuminboundqueue (MessageReceived)
		eventReqParams := util.SubscanEventRequestParams{
			Row:     100,
			Page:    0,
			Module:  "ethereumInboundQueue",
			EventId: "MessageReceived",
			Order:   "asc",
		}
		const maxEndBlockNum = 99999999
		for {
			eventReqParams.BlockRange = fmt.Sprintf("%d-%d", startCrawlNum, maxEndBlockNum)
			events, err := util.SubscanGetEvents(ctx, opt.bridgeHubName(), &eventReqParams)
			if len(events) == 0 {
				return nil, errors.New("not found message")
			}
			if err != nil {
				log.Printf("SubscanGetEvents get err %v\n", err)
				return nil, err
			}
			for _, e := range events {
				if e.Params[2].Value.(string) == messageId {
					log.Printf("Find bridge message has process in extrinsic index %s,event index %s \n", e.ExtrinsicIndex, e.EventIndex)
					extrinsicIndexArr := strings.Split(e.ExtrinsicIndex, "-")
					event := Event{
						ExtrinsicIdx: util.ToInt(extrinsicIndexArr[1]),
						BlockTime:    int64(e.BlockTimestamp),
						BlockNum:     util.ToInt(extrinsicIndexArr[0]),
					}
					return &event, nil
				}
				startCrawlNum = int(findOutBlockByExtrinsicIndex(e.ExtrinsicIndex).BlockNum)
			}
		}
	}

	// polkadot -> ethereum

	if opt.ExtrinsicIndex == "" {
		return nil, extrinsicIndexEmptyError
	}
	if opt.BridgeHubEndpoint == "" {
		return nil, bridgeHubEndpointEmptyError
	}
	if opt.OriginEndpoint == "" {
		return nil, originEndpointEmptyError
	}
	if opt.RelayEndpoint == "" {
		return nil, relayEndpointEmptyError
	}

	h := Hrmp{extrinsicIndex: opt.ExtrinsicIndex, originEndpoint: opt.OriginEndpoint, destEndpoint: opt.BridgeHubEndpoint, relayChainEndpoint: opt.RelayEndpoint}

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

	// track assetHub => bridgeHub HRMP message
	event, err := h.Track(ctx)

	if err != nil {
		return nil, err
	}
	etherStartBlockNum, err := util.EtherscanGetBlockByTime(ctx, isPolkadot, event.BlockTime)
	if err != nil {
		return nil, err
	}
	log.Println("Get etherStartBlockNum", etherStartBlockNum)

	logs, err := util.EtherscanGetLogs(ctx, isPolkadot, uint64(etherStartBlockNum), bridgeContract[opt.bridgeHubName()], InboundMessageDispatchedTopic, 1, 1000)
	if err != nil {
		return nil, err
	}
	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}
		// InboundMessageDispatchedTopic InboundMessageDispatched (index_topic_1 bytes32 channelID, uint64 nonce, index_topic_2 bytes32 messageID, bool success)
		// check message id is equal
		if l.Topics[2] == h.messageId {
			log.Println("Find bridge message have process in", util.HexToUint64(l.BlockNumber), l.TransactionHash)
			event.TxHash = l.TransactionHash
			return event, nil
		}
	}
	return nil, nil
}
