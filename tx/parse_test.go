package tx

import (
	"github.com/itering/scale.go/utiles"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseXcmMessageInstruction(t *testing.T) {
	client := initClient("wss://rococo-rpc.polkadot.io")
	defer client.Close()
	cases := []struct {
		Raw         string
		Instruction string
	}{
		{
			Raw:         "0x0000040a000780b8cfac090807010000000000000000005ed0b20000000000000104010102004e119e3b472c5c219da0ae9a131f3e2c30332b05b249e8bbdb9996c1dcd61843",
			Instruction: "{\"V0\":{\"WithdrawAsset\":{\"assets\":[{\"ConcreteFungible\":{\"amount\":\"41554000000\",\"id\":{\"Here\":null}}}],\"effects\":[{\"BuyExecution\":{\"debt\":3000000000,\"fees\":{\"All\":null},\"haltOnError\":false,\"weight\":0,\"xcm\":null}},{\"DepositAsset\":{\"assets\":[{\"All\":null}],\"dest\":{\"X1\":{\"AccountId32\":{\"id\":\"0x4e119e3b472c5c219da0ae9a131f3e2c30332b05b249e8bbdb9996c1dcd61843\",\"network\":{\"Any\":null}}}}}}]}}}",
		},
		{
			Raw:         "0x01010400010100591f00170000d01309468e1501080700010100591f00170000d01309468e1501000000000000000000ca9a3b0000000001000101000100000000010300b926e36d439106090be1151347cfb916e44afe00",
			Instruction: "{\"V1\":{\"ReserveAssetDeposit\":{\"assets\":[{\"fun\":{\"Fungible\":\"20000000000000000000\"},\"id\":{\"Concrete\":{\"interior\":{\"X1\":{\"Parachain\":2006}},\"parents\":1}}}],\"effects\":[{\"BuyExecution\":{\"debt\":1000000000,\"fees\":{\"fun\":{\"Fungible\":\"20000000000000000000\"},\"id\":{\"Concrete\":{\"interior\":{\"X1\":{\"Parachain\":2006}},\"parents\":1}}},\"halt_on_error\":true,\"instructions\":null,\"weight\":0}},{\"DepositAsset\":{\"assets\":{\"Wild\":{\"All\":null}},\"beneficiary\":{\"interior\":{\"X1\":{\"AccountKey20\":{\"key\":\"0xb926e36d439106090be1151347cfb916e44afe00\",\"network\":{\"Any\":null}}}},\"parents\":0},\"max_assets\":1}}]}}}",
		},
		{
			Raw:         "0x0210000400000106080002000b4512385180010a1300000106080002000b451238518001010300286bee0d01000400010100d817ce752e152f7e3f92212959bcafeee4f7208c25da3a15cb9b6e8fa4e00b1f",
			Instruction: "{\"V2\":[{\"WithdrawAsset\":[{\"fun\":{\"Fungible\":\"1650630070853\"},\"id\":{\"Concrete\":{\"interior\":{\"X1\":{\"GeneralKey\":\"0x0002\"}},\"parents\":0}}}]},{\"ClearOrigin\":\"NULL\"},{\"BuyExecution\":{\"fees\":{\"id\":{\"Concrete\":{\"interior\":{\"X1\":{\"GeneralKey\":\"0x0002\"}},\"parents\":0}},\"fun\":{\"Fungible\":\"1650630070853\"}},\"weight_limit\":{\"Limited\":4000000000}}},{\"DepositAsset\":{\"assets\":{\"Wild\":{\"All\":\"NULL\"}},\"max_assets\":1,\"beneficiary\":{\"interior\":{\"X1\":{\"AccountId32\":{\"network\":{\"Any\":\"NULL\"},\"id\":\"0xd817ce752e152f7e3f92212959bcafeee4f7208c25da3a15cb9b6e8fa4e00b1f\"}}},\"parents\":0}}}]}",
		},
		{
			Raw:         "0x0210000400000106080002000b00203d88792d0a1300000106080002000b00203d88792d010300286bee0d01000400010100a013f022969e34b78a40b1b1bc5557c369e8fbf03236cdb805643851b2221d26",
			Instruction: "{\"V2\":[{\"WithdrawAsset\":[{\"fun\":{\"Fungible\":\"50000000000000\"},\"id\":{\"Concrete\":{\"interior\":{\"X1\":{\"GeneralKey\":\"0x0002\"}},\"parents\":0}}}]},{\"ClearOrigin\":\"NULL\"},{\"BuyExecution\":{\"fees\":{\"id\":{\"Concrete\":{\"interior\":{\"X1\":{\"GeneralKey\":\"0x0002\"}},\"parents\":0}},\"fun\":{\"Fungible\":\"50000000000000\"}},\"weight_limit\":{\"Limited\":4000000000}}},{\"DepositAsset\":{\"assets\":{\"Wild\":{\"All\":\"NULL\"}},\"max_assets\":1,\"beneficiary\":{\"interior\":{\"X1\":{\"AccountId32\":{\"network\":{\"Any\":\"NULL\"},\"id\":\"0xa013f022969e34b78a40b1b1bc5557c369e8fbf03236cdb805643851b2221d26\"}}},\"parents\":0}}}]}",
		},
		{
			Raw:         "0x0310000400000000079e2d3d47280a1300000000079e2d3d4728010300286bee020010000d010204000101003628971d6b91628910aceeed80f922a1c539fa6bb201733d464b883acdd81b33",
			Instruction: "{\"V3\":[{\"WithdrawAsset\":[{\"fun\":{\"Fungible\":\"172993883550\"},\"id\":{\"Concrete\":{\"interior\":{\"Here\":\"NULL\"},\"parents\":0}}}]},{\"ClearOrigin\":\"NULL\"},{\"BuyExecution\":{\"fees\":{\"id\":{\"Concrete\":{\"interior\":{\"Here\":\"NULL\"},\"parents\":0}},\"fun\":{\"Fungible\":\"172993883550\"}},\"weight_limit\":{\"Limited\":{\"proof_size\":262144,\"ref_time\":4000000000}}}},{\"DepositAsset\":{\"assets\":{\"Wild\":{\"AllCounted\":1}},\"max_assets\":0,\"beneficiary\":{\"interior\":{\"X1\":{\"AccountId32\":{\"network\":null,\"id\":\"0x3628971d6b91628910aceeed80f922a1c539fa6bb201733d464b883acdd81b33\"}}},\"parents\":0}}}]}",
		},
		{
			Raw:         "0x031000040000000007f5c1998d2a0a130000000007f5c1998d2a000d01020400010100ea294590dbcfac4dda7acd6256078be26183d079e2739dd1e8b1ba55d94c957a",
			Instruction: "{\"V3\":[{\"WithdrawAsset\":[{\"fun\":{\"Fungible\":\"182764290549\"},\"id\":{\"Concrete\":{\"interior\":{\"Here\":\"NULL\"},\"parents\":0}}}]},{\"ClearOrigin\":\"NULL\"},{\"BuyExecution\":{\"fees\":{\"id\":{\"Concrete\":{\"interior\":{\"Here\":\"NULL\"},\"parents\":0}},\"fun\":{\"Fungible\":\"182764290549\"}},\"weight_limit\":{\"Unlimited\":\"NULL\"}}},{\"DepositAsset\":{\"assets\":{\"Wild\":{\"AllCounted\":1}},\"max_assets\":0,\"beneficiary\":{\"interior\":{\"X1\":{\"AccountId32\":{\"network\":null,\"id\":\"0xea294590dbcfac4dda7acd6256078be26183d079e2739dd1e8b1ba55d94c957a\"}}},\"parents\":0}}}]}",
		},
		{
			Raw:         "0x041402040100000700e40b54020a130100000700e40b5402000d0102040001010105d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d2c54d4fdcbe0c2034c399c945db38c2284bbc7eae8adbc03b753bc777c9a463d6d",
			Instruction: "{\"V4\":[{\"ReceiveTeleportedAsset\":[{\"fun\":{\"Fungible\":\"10000000000\"},\"id\":{\"interior\":{\"Here\":\"NULL\"},\"parents\":1}}]},{\"ClearOrigin\":\"NULL\"},{\"BuyExecution\":{\"fees\":{\"fun\":{\"Fungible\":\"10000000000\"},\"id\":{\"interior\":{\"Here\":\"NULL\"},\"parents\":1}},\"weight_limit\":{\"Unlimited\":\"NULL\"}}},{\"DepositAsset\":{\"assets\":{\"Wild\":{\"AllCounted\":1}},\"beneficiary\":{\"interior\":{\"X1\":[{\"AccountId32\":{\"id\":\"0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d\",\"network\":{\"Rococo\":\"NULL\"}}}]},\"parents\":0}}},{\"SetTopic\":\"0x54d4fdcbe0c2034c399c945db38c2284bbc7eae8adbc03b753bc777c9a463d6d\"}]}",
		},
	}

	for _, v := range cases {
		instruction, err := client.ParseXcmMessageInstruction(v.Raw)
		assert.NoError(t, err)
		assert.Equal(t, v.Instruction, utiles.ToString(instruction))
	}
	// Will raise MessageRawIsEmptyErr error
	_, err := client.ParseXcmMessageInstruction("")
	assert.ErrorIs(t, err, MessageRawIsEmptyErr)

}
