package tracker

import (
	"context"
	"errors"
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/gmajor-encrypt/xcm-tools/util"
	"log"
)

type TrackBridgeMessageOptions struct {
	Tx                string
	ChainId           uint
	ExtrinsicIndex    string
	BridgeHubEndpoint string // bridge hub websocket endpoint
	OriginEndpoint    string
	RelayEndpoint     string
}

var (
	bridgeHubName  = "bridgehub-rococo"                           // Subscan api endpoint name, for example: bridgehub-rococo,bridgehub-polkadot
	bridgeContract = "0x5B4909cE6Ca82d2CE23BD46738953c7959E710Cd" // bridge contract address
)

// TrackBridgeMessage ethereum <=> polkadot
func TrackBridgeMessage(ctx context.Context, opt *TrackBridgeMessageOptions) (*Event, error) {
	// ethereum -> polkadot
	// if tx is not empty, will track ethereum => polkadot message
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

		log.Println("Get ethereum message Id", messageId, "timestamp", timestamp, "blockNum", blockNum)

		polkadotBlock, err := util.SubscanGetBlockByTime(ctx, bridgeHubName, uint(timestamp))
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
			events, err := util.SubscanGetEvents(ctx, bridgeHubName, &eventReqParams)
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
					return nil, nil
				}
				startCrawlNum = int(findOutBlockByExtrinsicIndex(e.ExtrinsicIndex).BlockNum)
			}
		}
	}

	// polkadot -> ethereum

	if opt.ExtrinsicIndex == "" {
		return nil, errors.New("extrinsicIndex is empty")
	}
	if opt.BridgeHubEndpoint == "" {
		return nil, errors.New("bridgeHubEndpoint is empty")
	}
	if opt.OriginEndpoint == "" {
		return nil, errors.New("originEndpoint is empty")
	}
	if opt.RelayEndpoint == "" {
		return nil, errors.New("relayEndpoint is empty")
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
