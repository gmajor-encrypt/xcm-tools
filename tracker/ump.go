package tracker

import (
	"context"
	"errors"
	"github.com/gmajor-encrypt/xcm-tools/parse"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"log"
)

type Ump struct {
	// Sent UMP message extrinsic index
	ExtrinsicIndex string
	// origin chain websocket endpoint
	OriginEndpoint string
	// destination chain websocket endpoint
	DestEndpoint string
}

func (u *Ump) Track(ctx context.Context) (*Event, error) {
	extrinsic := findOutBlockByExtrinsicIndex(u.ExtrinsicIndex)
	if extrinsic == nil {
		return nil, InvalidExtrinsic
	}

	// origin para chain
	client, metadataInstant, closeClient := CreateSnapshotClient(u.OriginEndpoint)
	blockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum))
	if err != nil {
		return nil, err
	}

	metadataStruct := types.MetadataStruct(*metadataInstant)
	events, err := getEvents(ctx, client, blockHash, &metadataStruct)
	if err != nil {
		return nil, err
	}
	// find UpwardMessageSent event get MessageHash
	event := findEventByEventId(events, extrinsic.Index, []string{"UpwardMessageSent"})
	if event == nil {
		return nil, NotfoundXcmMessageErr
	}

	messageHash := event.Params[0].Value.(string)
	log.Println("Find messageHash", messageHash)

	// Get messageRaw by messageHash
	messageRaw, err := getUmpMessageByMessageHash(messageHash, blockHash)
	if err != nil {
		return nil, err
	}
	log.Println("Find messageRaw", messageRaw)

	// pickout topic id
	parseClient := parse.New(metadataInstant)
	instruction, err := parseClient.ParseXcmMessageInstruction(messageRaw)
	if err != nil {
		return nil, err
	}
	messageId := instruction.PickoutTopicId()
	log.Println("Find messageId", messageId)

	nextBlockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum+1))
	if err != nil {
		return nil, err
	}
	log.Println("Find nextBlock Hash", nextBlockHash)

	relayChainBlockNum, err := HRMPWatermark(nextBlockHash)
	if err != nil {
		return nil, err
	}

	log.Println("Find relayChainBlockNum", relayChainBlockNum)
	closeClient()
	types.Clean()

	client, _, closeClient = CreateSnapshotClient(u.DestEndpoint)
	defer closeClient()
	blockHash, err = rpc.GetChainGetBlockHash(client.Conn, relayChainBlockNum)
	if err != nil {
		return nil, err
	}

	log.Println("Find relaychain blockHash", blockHash)
	raw, err := rpc.GetMetadataByHash(nil, blockHash)
	if err != nil {
		return nil, err
	}

	metadataInstant = metadata.RegNewMetadataType(0, raw)
	metadataStruct = types.MetadataStruct(*metadataInstant)
	var retryTime int
	for {
		if retryTime > 5 {
			break
		}
		if blockHash == "" {
			relayChainBlockNum++
			blockHash, err = rpc.GetChainGetBlockHash(client.Conn, relayChainBlockNum)
			if err != nil {
				return nil, err
			}
		}
		log.Println("Get relaychain events with blockHash", blockHash)
		events, err = getEvents(ctx, client, blockHash, &metadataStruct)
		if err != nil {
			return nil, err
		}
		event = findEventByEventId(events, 0, []string{"Processed"})

		// Processed event [id, origin, weight, result]
		if event != nil {
			value := event.Params[0].Value.(string)
			// Find by messageId
			if value == messageId {
				log.Printf("Find UMP messageHash %s,messageId %s result %t", messageHash, messageId, event.Params[3].Value)
				return event, nil
			}
			// Find by messageHash
			if value == messageHash {
				log.Printf("Find UMP messageHash %s, result %t", event.Params[0].Value.(string), event.Params[3].Value)
				return event, nil
			}
		}
		retryTime++
		blockHash = ""
	}
	return nil, errors.New("not found xcm exec result")
}

// getUmpMessageByMessageHash get message by query parachainSystem.UpwardMessages
// filter by messageHash
func getUmpMessageByMessageHash(messageHash, blockHash string) (string, error) {
	messages, err := UpwardMessages(blockHash)
	if err != nil {
		return "", err
	}
	for _, message := range messages {
		if hash(message) == messageHash {
			return message, nil
		}
	}
	return "", errors.New("not found message")
}
