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

type Hrmp struct {
	extrinsicIndex     string
	originEndpoint     string
	destEndpoint       string
	relayChainEndpoint string
}

func (h *Hrmp) Track(ctx context.Context) (*Event, error) {
	// origin parachain
	client, metadataInstant, closeClient := CreateSnapshotClient(h.originEndpoint)

	extrinsic := findOutBlockByExtrinsicIndex(h.extrinsicIndex)
	if extrinsic == nil {
		return nil, NotfoundXcmMessageErr
	}
	blockHash, err := rpc.GetChainGetBlockHash(client.Conn, int(extrinsic.BlockNum))
	if err != nil {
		return nil, err
	}

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
	log.Println("messageHash", messageHash)

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

	// relay chain
	client, _, closeClient = CreateSnapshotClient(h.relayChainEndpoint)
	nextBlockHash, err = rpc.GetChainGetBlockHash(client.Conn, relayChainBlockNum+2)
	if err != nil {
		return nil, err
	}

	log.Println("relaychain blockHash", blockHash)
	raw, err := rpc.GetMetadataByHash(nil, nextBlockHash)
	if err != nil {
		return nil, err
	}

	metadataInstant = metadata.RegNewMetadataType(0, raw)
	metadataStruct = types.MetadataStruct(*metadataInstant)
	log.Println("nextBlockHash", nextBlockHash)
	pendingAvailability, err := PendingAvailability(destParaId, nextBlockHash)
	if err != nil {
		return nil, err
	}

	paraHead := pendingAvailability.Descriptor.ParaHead

	log.Println("get para block hash", paraHead)
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
	events, err = getEvents(ctx, client, paraHead, &metadataStruct)
	if err != nil {
		return nil, err
	}

	event = findEventByEventId(events, 1, []string{"Success", "Failed"})
	if event != nil && event.Params[0].Value.(string) == messageHash {
		log.Printf("find HRMP messageHash %s, result %s", event.Params[0].Value.(string), event.EventId)
		return event, nil
	}
	return nil, nil
}
