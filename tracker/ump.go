package tracker

import (
	"context"
	"errors"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"log"
)

type Ump struct {
	ExtrinsicIndex string
	OriginEndpoint string
	DestEndpoint   string
}

func (u *Ump) Track(ctx context.Context) (*Event, error) {
	extrinsic := findOutBlockByExtrinsicIndex(u.ExtrinsicIndex)
	if extrinsic == nil {
		return nil, NotfoundXcmMessageErr
	}
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
	event := findEventByEventId(events, extrinsic.Index, []string{"UpwardMessageSent"})
	if event == nil {
		return nil, NotfoundXcmMessageErr
	}

	messageHash := event.Params[0].Value.(string)
	log.Println("messageHash", messageHash)

	nextBlockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum+1))
	if err != nil {
		return nil, err
	}
	log.Println("nextBlockHash", nextBlockHash)

	relayChainBlockNum, err := HRMPWatermark(nextBlockHash)
	if err != nil {
		return nil, err
	}
	log.Println("relayChainBlockNum", relayChainBlockNum)
	closeClient()
	types.Clean()

	client, _, closeClient = CreateSnapshotClient(u.DestEndpoint)
	defer closeClient()
	blockHash, err = rpc.GetChainGetBlockHash(client.Conn, relayChainBlockNum)
	if err != nil {
		return nil, err
	}

	log.Println("relaychain blockHash", blockHash)
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
		log.Println("get relaychain events with blockHash", blockHash)
		events, err = getEvents(ctx, client, blockHash, &metadataStruct)
		if err != nil {
			return nil, err
		}
		event = findEventByEventId(events, 0, []string{"Processed"})
		if event != nil && event.Params[0].Value.(string) == messageHash {
			log.Printf("find UMP messageHash %s, result %t", event.Params[0].Value.(string), event.Params[3].Value)
			return event, nil
		}
		retryTime++
		blockHash = ""
	}
	return nil, errors.New("not found xcm exec result")
}
