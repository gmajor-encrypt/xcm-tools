package tx

import (
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/shopspring/decimal"
	"strings"

	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
)

type Client struct {
	Conn               *rpc.Client
	Ump                IXCMP
	Dmp                IXCMP
	Hrmp               IXCMP
	m                  *metadata.Instant
	XcmVersionTypeName string
}

// NewClient create a new XCM client
func NewClient(endpoint string) *Client {
	websocket.SetEndpoint(endpoint)
	// websocket init and set metadata
	_, err := websocket.Init()
	if err != nil {
		panic(err)
	}

	raw, err := rpc.GetMetadataByHash(nil)
	if err != nil {
		panic(err)
	}
	c := &rpc.Client{}
	metadataInstant := metadata.RegNewMetadataType(0, raw)
	c.SetMetadata(metadataInstant)
	client := &Client{
		Conn: c,
		Ump:  NewMessage(UMP),
		Dmp:  NewMessage(DMP),
		Hrmp: NewMessage(HRMP),
		m:    metadataInstant,
	}
	return client
}

func (c *Client) SetKeyRing(sk string) {
	c.Conn.SetKeyRing(keyring.New(keyring.Sr25519Type, sk))
}

func (c *Client) Close() {
	websocket.Close()
}

func (c *Client) Metadata() *metadata.Instant {
	return c.m
}

// GetCallByName get call by module name and call name
func GetCallByName(moduleName, callName string, m *metadata.Instant) *types.MetadataCalls {
	module := GetModule(moduleName, m)
	if module == nil {
		return nil
	}

	for i, call := range module.Calls {
		if strings.EqualFold(call.Name, callName) {
			return &module.Calls[i]
		}
	}
	return nil
}

// GetModule get module by name
func GetModule(moduleName string, m *metadata.Instant) *types.MetadataModules {
	for i, v := range m.Metadata.Modules {
		if strings.EqualFold(v.Name, moduleName) {
			return &m.Metadata.Modules[i]
		}
	}
	return nil
}

// SendUmpTransfer send XCM UMP message
// accountId: the account id of the beneficiary
// amount: the amount of the asset to be transferred
func (c *Client) SendUmpTransfer(accountId string, amount decimal.Decimal) (string, error) {
	callName, args := c.Ump.LimitedTeleportAssets(
		SimplifyMultiLocationRelayChain(),
		SimplifyMultiLocationAccountId32(accountId),
		SimplifyMultiAssets(amount),
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Ump.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}

// SendDmpTransfer send XCM DMP message
// paraId: the para id of the parachain
// accountId: the account id of the beneficiary
// amount: the amount of the asset to be transferred
func (c *Client) SendDmpTransfer(paraId uint32, accountId string, amount decimal.Decimal) (string, error) {
	callName, args := c.Dmp.LimitedTeleportAssets(
		SimplifyMultiLocationParaId(paraId),
		SimplifyMultiLocationAccountId32(accountId),
		SimplifyMultiAssets(amount),
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Dmp.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}

// SendHrmpTransfer send XCM HRMP message
// paraId: the para id of the parachain
// accountId: the account id of the beneficiary
// amount: the amount of the asset to be transferred
func (c *Client) SendHrmpTransfer(paraId uint32, accountId string, amount decimal.Decimal) (string, error) {
	callName, args := c.Hrmp.LimitedReserveTransferAssets(
		SimplifyMultiLocationParaId(paraId),
		SimplifyMultiLocationAccountId32(accountId),
		SimplifyMultiAssets(amount),
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Hrmp.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}

func (c *Client) SendTokenToEthereum(h160, tokenContract string, amount decimal.Decimal, chainId uint64) (string, error) {
	callName, args := c.Hrmp.TransferAssets(
		&VersionedMultiLocation{V3: &V3MultiLocation{
			Interior: V3MultiLocationJunctions{X1: &XCMJunctionV3{GlobalConsensus: &GlobalConsensusNetworkId{Ethereum: &chainId}}}, Parents: 2,
		}},
		&VersionedMultiLocation{V3: &V3MultiLocation{
			Interior: V3MultiLocationJunctions{X1: &XCMJunctionV3{AccountKey20: &XCMJunctionV3AccountKey20{Key: h160}}}, Parents: 0,
		}},
		&MultiAssets{V3: []V3MultiAssets{SimplifyEthereumAssets(chainId, tokenContract, amount)}},
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Hrmp.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}
