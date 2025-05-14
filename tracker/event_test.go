package tracker

import (
	"context"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getEventsFromChain(t *testing.T) {
	client, closeClient := CreateSnapshotClient("wss://moonbeam-rpc.dwellir.com")
	RegDefaultMetadata()
	defer closeClient()
	ctx := context.TODO()
	cases := []struct {
		blockHash string
		EventRaw  string
	}{
		{
			blockHash: "0x3bc4c37edbaffc7bb0cf19e314fbf74ea18bad3779d977b003809b688e35c18c",
			EventRaw:  "0x14025008d7ce4754605fc2350141040000000000abb08e8f3f033a360141040000000000000000000000000002c0cd1715e1020100000100000000004216c63db90f02000000020000000000c2beb359c9a202010000030000000000028e957cf592020100",
		},
		{
			blockHash: "0xae45308edf1fc25d09ebb4cff57658fc72129dd0efd875394114c9c451dc871d",
			EventRaw:  "0x14025008716fab4e31e98a45024104000000000089c22d646fb088460241040000000000000000000000000002c0cd176dea020100000100000000004216c63db90f02000000020000000000c2beb359c9a202010000030000000000028e957cf592020100",
		},
	}
	for _, v := range cases {
		raw, err := getEventsFromChain(ctx, client.Conn, v.blockHash)
		assert.NoError(t, err)
		assert.Equal(t, raw, v.EventRaw)
	}
}

func Test_getEvents(t *testing.T) {
	client, closeClient := CreateSnapshotClient("wss://moonbeam-rpc.dwellir.com")
	defer closeClient()
	ctx := context.TODO()
	cases := []struct {
		blockHash   string
		EventsCount int
	}{
		{
			blockHash:   "0x3bc4c37edbaffc7bb0cf19e314fbf74ea18bad3779d977b003809b688e35c18c",
			EventsCount: 5,
		}, {
			blockHash:   "0x28471ff0e1b727c71d368530d944bbe71d2da978800b8bff71762ca9dc3c072c",
			EventsCount: 53,
		},
	}

	for _, v := range cases {
		metadataRaw, err := rpc.GetMetadataByHash(nil, v.blockHash)
		if err != nil {
			panic(err)
		}
		metadataInstant := metadata.RegNewMetadataType(0, metadataRaw)
		metadataStruct := types.MetadataStruct(*metadataInstant)
		raw, err := getEvents(ctx, client, v.blockHash, &metadataStruct)
		assert.NoError(t, err)

		assert.Equal(t, v.EventsCount, len(raw))
	}
}

func Test_findEventByEventId(t *testing.T) {
	events := []Event{
		{
			ExtrinsicIdx: 0,
			EventId:      "set_validation_data",
		},
		{
			ExtrinsicIdx: 1,
			EventId:      "transact",
		},
		{
			ExtrinsicIdx: 2,
			EventId:      "set_babe_randomness_results",
		},
	}

	assert.NotNil(t, findEventByEventId(events, 0, []string{"set_validation_data"}))
	assert.NotNil(t, findEventByEventId(events, 1, []string{"transact"}))
	assert.Nil(t, findEventByEventId(events, 2, []string{"transact"}))
}
