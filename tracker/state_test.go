package tracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHRMPWatermark(t *testing.T) {
	_, closeF := CreateSnapshotClient("wss://moonbeam-rpc.dwellir.com")
	_ = RegDefaultMetadata()
	blockNum, err := HRMPWatermark("0xebaa16b5cc4c53595acb541ea46fdf9d6625f00ea335e91c53391d13cb599f36")
	assert.NoError(t, err)
	assert.Equal(t, 17040495, blockNum)
	closeF()
}

func TestParasInfo(t *testing.T) {
	_, closeF := CreateSnapshotClient("wss://rococo-bridge-hub-rpc.polkadot.io")
	_ = RegDefaultMetadata()
	paraInfo, err := ParasInfo("bridgeWestendParachains", 1002, "0x332184c37060ffaa3986cd142354bc8c34dc06146ee65608d4533169b389e37d")
	assert.NoError(t, err)
	assert.Equal(t, paraInfo.BestHeadHash.HeadHash, "0x315c4f176ead8620140d90af90991092fd401644509a899290fdd240d08cfd54")
	closeF()
}
