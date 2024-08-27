package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/gmajor-encrypt/xcm-tools/util"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"log"
	"strings"
)

// S2STrackBridgeMessageOptions is the options for TrackS2sBridgeMessage
type S2STrackBridgeMessageOptions struct {
	// substrate => substrate message extrinsic index
	ExtrinsicIndex string

	// origin chain websocket endpoint
	OriginEndpoint string
	// origin bridge hub websocket endpoint
	BridgeHubEndpoint string
	// origin relay chain websocket endpoint
	OriginRelayEndpoint string

	// destination chain websocket endpoint
	DestinationEndpoint string
	// destination bridge hub websocket endpoint
	DestinationBridgeHubEndpoint string
	// destination relay chain websocket endpoint
	DestinationRelayEndpoint string
}

type ReceivedMessagesType struct {
	Lane           string `json:"lane"`
	ReceiveResults []struct {
		Col1 int `json:"col1"`
	} `json:"receive_results"`
}

// bridgeParaIds is the map of bridge module name and paraId
var bridgeParaIds = map[string]uint{
	// rococo
	"BridgeRococoParachains": 1013,
	// westend
	"bridgeWestendParachains": 1002,
	// kusama
	"bridgeKusamaParachains": 1002,
	// polkadot
	"bridgePolkadotParachains": 1002,
}

// TrackS2sBridgeMessage tracks the s2s bridge message
// flow: parachain => relay chain => bridge hub => dest bridge hub => dest relay chain => dest parachain
func TrackS2sBridgeMessage(ctx context.Context, opt *S2STrackBridgeMessageOptions) (*Event, error) {

	// checkBridgeModule checks the bridge module in metadata
	var checkBridgeModule = func(metadataInstant *metadata.Instant) (string, uint) {
		for moduleName, paraId := range bridgeParaIds {
			for _, module := range metadataInstant.Metadata.Modules {
				if strings.EqualFold(module.Name, moduleName) {
					return moduleName, paraId
				}
			}
		}
		panic("not found bridge module")
	}

	origin := Hrmp{
		extrinsicIndex:     opt.ExtrinsicIndex,
		originEndpoint:     opt.OriginEndpoint,
		destEndpoint:       opt.BridgeHubEndpoint,
		relayChainEndpoint: opt.OriginRelayEndpoint,
	}

	type ResourceData struct {
		LaneId    string
		Nonce     int
		BlockHash string
	}

	var filter = func(events []Event, i *tx.VersionedXcm, messageHash, blockHash string) (*Event, string, error) {
		messageId := i.PickoutTopicId()
		// messageQueue (Processed)
		event := findEventByEventId(events, 0, []string{"Processed"})
		if event == nil {
			return nil, "", nil
		}
		if event.Params[0].Value.(string) != messageId {
			return nil, "", nil
		}

		// get bridge lane_id and nonce
		// MessageAccepted(lane_id,nonce)
		messageAcceptedEvent := findEventByEventId(events, 0, []string{"MessageAccepted"})
		if messageAcceptedEvent == nil {
			return nil, "", fmt.Errorf("not found MessageAccepted event")
		}
		event.ResourceData = ResourceData{
			LaneId:    messageAcceptedEvent.Params[0].Value.(string),
			Nonce:     util.AnyToInt(messageAcceptedEvent.Params[1].Value),
			BlockHash: blockHash,
		}
		return event, messageId, nil
	}
	origin.filterCallBack = filter
	// common hrmp filter
	event, err := origin.Track(ctx)
	if err != nil {
		return nil, err
	}
	resourceData := event.ResourceData.(ResourceData)

	moduleName, paraId := checkBridgeModule(metadata.Latest(nil))

	paraInfo, err := ParasInfo(moduleName, paraId, resourceData.BlockHash)
	if err != nil {
		return nil, err
	}
	types.Clean()

	destBridgeHeaderHash := paraInfo.BestHeadHash.HeadHash
	log.Println("found destBridgeHeaderHash", destBridgeHeaderHash)

	websocket.Close()
	client, closeClient := CreateSnapshotClient(opt.DestinationBridgeHubEndpoint)

	raw, err := rpc.GetMetadataByHash(nil, destBridgeHeaderHash)
	if err != nil {
		return nil, err
	}
	metadataInstant := metadata.RegNewMetadataType(0, raw)
	metadataStruct := types.MetadataStruct(*metadataInstant)
	paraHead := destBridgeHeaderHash
	nextHash, err := rpc.GetChainGetBlockHash(client.Conn, 1)
	if err != nil {
		return nil, err
	}

	fmt.Println("nextHash", nextHash)

	var destChainExtrinsicIndex string
	var currenBlock int64
	var retry int

OuterLoop:
	// loop to find destChainExtrinsicIndex
	for {
		log.Println("start check block", currenBlock, paraHead)
		events, err := getEvents(ctx, client, paraHead, &metadataStruct)
		if err != nil {
			return nil, err
		}
		messagesReceivedEvent := findEventByEventId(events, -1, []string{"MessagesReceived"})
		if messagesReceivedEvent != nil {
			var receivedMessagesType []ReceivedMessagesType
			bytes, _ := json.Marshal(messagesReceivedEvent.Params[0].Value)
			_ = json.Unmarshal(bytes, &receivedMessagesType)
			for _, receivedMessages := range receivedMessagesType {
				lane := receivedMessages.Lane
				for _, receiveResult := range receivedMessages.ReceiveResults {
					nonce := receiveResult.Col1
					// match nonce and lane
					if nonce == resourceData.Nonce && lane == resourceData.LaneId {
						destChainExtrinsicIndex = fmt.Sprintf("%d-%d", currenBlock, messagesReceivedEvent.ExtrinsicIdx)
						break OuterLoop
					}
				}
			}

		}
		// check next block
		currenBlock, err = SystemNumber(paraHead)
		if err != nil {
			return nil, err
		}
		currenBlock++
		nextHash, err := rpc.GetChainGetBlockHash(client.Conn, int(currenBlock))
		if err != nil {
			return nil, err
		}
		paraHead = nextHash
		retry++
		if retry > 30 {
			break
		}
	}

	if destChainExtrinsicIndex == "" {
		return nil, fmt.Errorf("not found destChainExtrinsicIndex")
	}
	log.Printf("found destChainExtrinsicIndex: %s", destChainExtrinsicIndex)

	closeClient()

	destination := Hrmp{
		extrinsicIndex:     destChainExtrinsicIndex,
		originEndpoint:     opt.DestinationBridgeHubEndpoint,
		destEndpoint:       opt.DestinationEndpoint,
		relayChainEndpoint: opt.DestinationRelayEndpoint,
	}
	destination.filterCallBack = hrmpFilter
	// common hrmp filter
	return destination.Track(ctx)
}
