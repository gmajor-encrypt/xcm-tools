package tx

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Client(t *testing.T) {
	client := NewClient("wss://rpc.polkadot.io")

	assert.NotNil(t, GetModule("Balances", client.m))
	assert.Nil(t, GetModule("xtokens", client.m))

	assert.NotNil(t, GetCallByName("Balances", "transfer_keep_alive", client.m))
	assert.NotNil(t, GetCallByName("XcmPallet", "reserve_transfer_assets", client.m))
	client.Close()
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

func TestPolkadotToEthereum(t *testing.T) {
	client := initClient("wss://rococo-asset-hub-rpc.polkadot.io")
	client.SetKeyRing(AliceSeed2)
	defer client.Close()
	destH160 := strings.ToLower("0x6EB228b7ab726b8B44892e8e273ACF3dcC9C0492")
	t.Run("Test_XCM_To_Ethereum", func(t *testing.T) {
		txHash, err := client.SendTokenToEthereum(
			destH160,
			"0xfff9976782d46cc05630d1f6ebab18b2324d6b14",
			decimal.New(5000000000, 0),
			11155111)
		assert.NoError(t, err)
		assert.Len(t, txHash, 66)
	})
}
