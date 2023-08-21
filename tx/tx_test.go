package tx

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Client(t *testing.T) {
	client := NewClient("wss://rpc.polkadot.io")

	assert.NotNil(t, client.getModule("Balances"))
	assert.Nil(t, client.getModule("xtokens"))

	assert.NotNil(t, client.getCallByName("Balances", "transfer"))
	assert.NotNil(t, client.getCallByName("XcmPallet", "reserve_transfer_assets"))
	assert.EqualValues(t, client.XcmVersion, uint(3))
	client.Close()

	// will raise panic because this network doesn't support XCM
	assert.Panics(t, func() {
		NewClient("wss://mainnet-node.dock.io")
	})
}

func TestXcmTransfer(t *testing.T) {
	client := initClient("wss://rococo-asset-hub-rpc.polkadot.io")
	defer client.Close()

	t.Run("Test_XCM_Ump_Transfer", func(t *testing.T) {
		txHash, err := client.SendUmpTransfer("0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", decimal.New(1, 0))
		assert.NoError(t, err)
		assert.Len(t, txHash, 66)
	})

	t.Run("Test_XCM_HRMP_Send", func(t *testing.T) {
		txHash, err := client.SendHrmpTransfer(2087, "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", decimal.New(1, 0))
		assert.NoError(t, err)
		assert.Len(t, txHash, 66)
	})
}

func TestDmpTransfer(t *testing.T) {
	client := initClient("wss://rococo-rpc.polkadot.io")
	defer client.Close()

	t.Run("Test_XCM_Dmp_Transfer", func(t *testing.T) {
		txHash, err := client.SendDmpTransfer(1000, "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", decimal.New(1, 0))
		assert.NoError(t, err)
		assert.Len(t, txHash, 66)
	})
}
