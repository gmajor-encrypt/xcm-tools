package util

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SubscanGetBlockByTime(t *testing.T) {
	cases := []struct {
		network   string
		blockTime uint
		expected  *SubscanBlock
	}{
		{network: "polkadot", blockTime: 1610000000, expected: &SubscanBlock{BlockNum: 3234503, BlockTimestamp: 1609999998, Hash: "0x1cd9f421528d0cd219d736a48b3788a774f58995bbaaded76c979d34b5505aef"}},
		{network: "kusama", blockTime: 1610000000, expected: &SubscanBlock{BlockNum: 5666943, BlockTimestamp: 1609999998, Hash: "0x952470bc3dc9a82b9f60c1d45b7e01b736d639b20f60831f39bab05198456d59"}},
	}
	for _, v := range cases {
		block, err := SubscanGetBlockByTime(context.TODO(), v.network, v.blockTime)
		if err != nil {
			t.Errorf("SubscanGetBlockByTime(%s, %d) != %v", v.network, v.blockTime, v.expected)
		}
		assert.Equal(t, block.Hash, v.expected.Hash)
		assert.Equal(t, block.BlockNum, v.expected.BlockNum)
		assert.Equal(t, block.BlockTimestamp, v.expected.BlockTimestamp)
	}
}

func TestSubscanGetEvents(t *testing.T) {
	cases := []struct {
		network  string
		params   SubscanEventRequestParams
		expected *SubscanEvent
	}{
		{network: "polkadot", params: SubscanEventRequestParams{Row: 1, BlockRange: "19647148-19647148"}, expected: &SubscanEvent{
			Events: []SubscanEventWithParams{{EventIndex: "19647148-52", EventId: "ExtrinsicSuccess", ModuleId: "system"}},
		}},
		{network: "kusama", params: SubscanEventRequestParams{Row: 1, BlockRange: "19647101-19647101"}, expected: &SubscanEvent{
			Events: []SubscanEventWithParams{{EventIndex: "19647101-38", EventId: "ExtrinsicSuccess", ModuleId: "system"}},
		}}}
	for _, v := range cases {
		events, err := SubscanGetEvents(context.TODO(), v.network, &v.params)
		assert.NoError(t, err)
		assert.Equal(t, len(events), len(v.expected.Events))
		assert.Equal(t, events[0].EventIndex, v.expected.Events[0].EventIndex)
		assert.Equal(t, events[0].EventId, v.expected.Events[0].EventId)
		assert.Equal(t, events[0].ModuleId, v.expected.Events[0].ModuleId)
		assert.Greater(t, len(events[0].Params), 0) // params is not empty
	}
}
