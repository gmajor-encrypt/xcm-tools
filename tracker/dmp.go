package tracker

import (
	"context"
	"encoding/json"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"log"
)

type Dmp struct {
	extrinsicIndex string
	originEndpoint string
	destEndpoint   string
}

func (d *Dmp) Track(ctx context.Context) (*Event, error) {
	extrinsic := findOutBlockByExtrinsicIndex(d.extrinsicIndex)
	if extrinsic == nil {
		return nil, NotfoundXcmMessageErr
	}

	client, metadataInstant, closeClient := CreateSnapshotClient(d.originEndpoint)
	blockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum))

	if err != nil {
		return nil, err
	}

	metadataStruct := types.MetadataStruct(*metadataInstant)
	events, err := getEvents(ctx, client, blockHash, &metadataStruct)
	if err != nil {
		return nil, err
	}
	event := findEventByEventId(events, extrinsic.Index, []string{"Attempted"})
	if event == nil {
		return nil, NotfoundXcmMessageErr
	}
	// check xcm send success
	_, err = ParseAttempted(*event, 0)
	if err != nil {
		return nil, err
	}

	// get dest para id
	extrinsics, err := getExtrinsics(ctx, client, blockHash)
	if err != nil {
		return nil, err
	}
	extrinsicData, err := getExtrinsicByIndex(ctx, extrinsics, int(extrinsic.Index), &metadataStruct)
	if err != nil {
		return nil, err
	}

	var multiLocation tx.VersionedMultiLocation
	bytes, _ := json.Marshal(extrinsicData.Params[0].Value)
	err = json.Unmarshal(bytes, &multiLocation)
	if err != nil {
		return nil, err
	}
	destParaId := multiLocation.GetParaId()
	if destParaId == 0 {
		return nil, InvalidDestParaId
	}
	log.Println("destParaId", destParaId)

	downwardMessageQueues, err := DownwardMessageQueues(destParaId, blockHash)
	if err != nil {
		return nil, err
	}

	var messageHash string
	for _, message := range downwardMessageQueues {
		if message.SentAt == extrinsic.BlockNum {
			messageHash = hash(message.Msg)
			break
		}
	}

	log.Println("messageHash", messageHash)

	nextBlockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum+2))
	if err != nil {
		return nil, err
	}

	log.Println("nextBlockHash", nextBlockHash)
	pendingAvailability, err := PendingAvailability(destParaId, nextBlockHash)
	if err != nil {
		return nil, err
	}

	paraHead := pendingAvailability.Descriptor.ParaHead
	log.Println("get para block hash", paraHead)
	closeClient()
	types.Clean()

	client, _, closeClient = CreateSnapshotClient(d.destEndpoint)
	defer closeClient()
	raw, err := rpc.GetMetadataByHash(nil, paraHead)
	if err != nil {
		return nil, err
	}
	metadataInstant = metadata.RegNewMetadataType(0, raw)
	metadataStruct = types.MetadataStruct(*metadataInstant)
	events, err = getEvents(ctx, client, paraHead, &metadataStruct)
	if err != nil {
		return nil, err
	}
	event = findEventByEventId(events, 0, []string{"ExecutedDownward"})
	if event != nil && event.Params[0].Value.(string) == messageHash {
		result, _ := ParseAttempted(*event, 1)
		log.Printf("find DMP messageHash %s, result %t", event.Params[0].Value.(string), result)
		return event, nil
	}
	return nil, nil
}
