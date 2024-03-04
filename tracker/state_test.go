package tracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHRMPWatermark(t *testing.T) {
	_, _, closeF := CreateSnapshotClient("wss://moonbeam-rpc.dwellir.com")
	blockNum, err := HRMPWatermark("0xebaa16b5cc4c53595acb541ea46fdf9d6625f00ea335e91c53391d13cb599f36")
	assert.NoError(t, err)
	assert.Equal(t, 17040495, blockNum)
	closeF()
}
