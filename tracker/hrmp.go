package tracker

import (
	"context"
	"github.com/gmajor-encrypt/xcm-tools/parse"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	"log"
)

type Hrmp struct {
	// sent HRMP message extrinsic index
	extrinsicIndex string
	// origin parachain endpoint
	originEndpoint string
	// dest parachain endpoint
	destEndpoint string
	// relay chain endpoint
	relayChainEndpoint string
	filterCallBack     func(events []Event, i *tx.VersionedXcm, messageHash, blockHash string) (*Event, string, error)
	// hrmp message id
	messageId string
}

func (h *Hrmp) Track(ctx context.Context) (*Event, error) {
	// origin parachain
	client, metadataInstant, closeClient := CreateSnapshotClient(h.originEndpoint)

	extrinsic := findOutBlockByExtrinsicIndex(h.extrinsicIndex)
	if extrinsic == nil {
		return nil, InvalidExtrinsic
	}

	blockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum))
	if err != nil {
		return nil, err
	}

	// get spec runtime version
	raw, err := rpc.GetMetadataByHash(nil, blockHash)
	if err != nil {
		return nil, err
	}

	metadataInstant = metadata.RegNewMetadataType(0, raw)

	metadataStruct := types.MetadataStruct(*metadataInstant)
	events, err := getEvents(ctx, client, blockHash, &metadataStruct)
	if err != nil {
		return nil, err
	}

	event := findEventByEventId(events, extrinsic.Index, []string{"XcmpMessageSent"})
	if event == nil {
		return nil, NotfoundXcmMessageErr
	}

	messageHash := event.Params[0].Value.(string)
	log.Println("Find messageHash", messageHash)

	hrmpOutboundMessages, err := HrmpOutboundMessages(blockHash)
	if err != nil {
		return nil, err
	}
	var (
		messageRaw string
		destParaId uint
	)
	parseClient := parse.New(metadataInstant)

	// pick the message by messageHash
	for _, message := range hrmpOutboundMessages {
		for _, raw := range parseClient.DecodeFixedMessage(util.TrimHex(message.Data)[2:]) {
			if hash(raw) == messageHash {
				messageRaw = util.TrimHex(raw)
				destParaId = message.Recipient
			}
		}
	}
	log.Println("Find messageRaw", messageRaw)

	instruction, err := parseClient.ParseXcmMessageInstruction(messageRaw)
	if err != nil {
		return nil, err
	}

	nextBlockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum+1))
	if err != nil {
		return nil, err
	}
	log.Println("Get NextBlockHash", nextBlockHash)

	relayChainBlockNum, err := HRMPWatermark(nextBlockHash)
	if err != nil {
		return nil, err
	}
	log.Println("Get RelayChainBlockNum", relayChainBlockNum)
	closeClient()
	types.Clean()

	// relay chain
	client, _, closeClient = CreateSnapshotClient(h.relayChainEndpoint)
	nextBlockHash, err = rpc.GetChainGetBlockHash(client.Conn, relayChainBlockNum)
	if err != nil {
		return nil, err
	}

	raw, err = rpc.GetMetadataByHash(nil, nextBlockHash)
	if err != nil {
		return nil, err
	}

	metadataInstant = metadata.RegNewMetadataType(0, raw)
	metadataStruct = types.MetadataStruct(*metadataInstant)
	log.Println("Get NextBlockHash", nextBlockHash)

	paraHeadBlockNum := relayChainBlockNum + 2
	var paraHead string
	var retry int
	for {
		nextBlockHash, err = rpc.GetChainGetBlockHash(client.Conn, paraHeadBlockNum)
		if err != nil {
			return nil, err
		}
		log.Println("Find nextBlockHash,start fetch PendingAvailability", nextBlockHash)
		pendingAvailability, err := PendingAvailability(destParaId, nextBlockHash)
		if err != nil {
			return nil, err
		}
		paraHead = pendingAvailability.Descriptor.ParaHead
		retry++
		paraHeadBlockNum++
		if paraHead != "" || retry > 5 {
			break
		}
	}

	log.Println("Get para block hash", paraHead)
	closeClient()
	types.Clean()

	// dest parachain
	client, _, closeClient = CreateSnapshotClient(h.destEndpoint)
	defer closeClient()
	raw, err = rpc.GetMetadataByHash(nil, paraHead)
	if err != nil {
		return nil, err
	}

	metadataInstant = metadata.RegNewMetadataType(0, raw)
	metadataStruct = types.MetadataStruct(*metadataInstant)
	retry = 0
	var messageId string
	for {
		events, err = getEvents(ctx, client, paraHead, &metadataStruct)
		if err != nil {
			return nil, err
		}
		event, messageId, err = h.filterCallBack(events, instruction, messageHash, paraHead)
		if err != nil {
			return nil, err
		}
		h.messageId = messageId
		if event != nil {
			return event, nil
		}
		// check next block
		currenBlock, err := SystemNumber(paraHead)
		if err != nil {
			return nil, err
		}
		nextHash, err := rpc.GetChainGetBlockHash(client.Conn, int(currenBlock+1))
		if err != nil {
			return nil, err
		}
		paraHead = nextHash
		retry++
		if retry > 5 {
			break
		}
	}
	return nil, nil
}

func hrmpFilter(events []Event, i *tx.VersionedXcm, messageHash, _ string) (*Event, string, error) {
	// 	xcmpQueue (Success) [messageHash, messageId, result]
	// 	xcmpQueue (Failed)  [messageHash, messageId, result]
	messageId := i.PickoutTopicId()
	event := findEventByEventId(events, 1, []string{"Success", "Failed"})
	if event != nil {
		// find messageHash
		if len(event.Params) == 2 {
			if event.Params[0].Value.(string) == messageHash {
				log.Printf("Find HRMP messageHash %s, result %s", messageHash, event.EventId)
				return event, messageId, nil
			}
			// find messageId
		} else if len(event.Params) == 3 {
			if event.Params[1].Value.(string) == messageId {
				log.Printf("Find HRMP messageHash %s,messageId %s result %s", messageHash, messageId, event.EventId)
				return event, messageId, nil
			}
		}
	}
	return nil, "", nil
}
