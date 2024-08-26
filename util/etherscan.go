package util

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	EtherscanAPIKey = "46K2QRJJ5KBAI7BBPK8CPGSYKEWS7ZAJEU"
)

func etherscanAPIBaseURL(isMainnet bool) string {
	if isMainnet {
		return "https://api.etherscan.io/api"
	}
	return "https://api-sepolia.etherscan.io/api"
}

type EtherscanProxyRes[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  T      `json:"result"`
}

type EtherscanTransaction struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	AccessList           []interface{} `json:"accessList"`
	ChainId              string        `json:"chainId"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	YParity              string        `json:"yParity"`
}

// EthGetTransactionByHash returns the transaction details by hash
// eth_getTransactionByHash
func EthGetTransactionByHash(ctx context.Context, isMainnet bool, hash string) (*EtherscanTransaction, error) {
	var endpoint = etherscanAPIBaseURL(isMainnet) + "?module=proxy&action=eth_getTransactionByHash&txhash=" + hash + "&apikey=" + EtherscanAPIKey
	body, err := HttpGet(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	var txn EtherscanProxyRes[EtherscanTransaction]
	if err = json.Unmarshal(body, &txn); err != nil {
		return nil, err
	}
	return &txn.Result, nil
}

type EtherscanTransactionReceipt struct {
	BlockHash         string          `json:"blockHash"`
	BlockNumber       string          `json:"blockNumber"`
	ContractAddress   interface{}     `json:"contractAddress"`
	CumulativeGasUsed string          `json:"cumulativeGasUsed"`
	EffectiveGasPrice string          `json:"effectiveGasPrice"`
	From              string          `json:"from"`
	GasUsed           string          `json:"gasUsed"`
	Logs              []EthReceiptLog `json:"logs"`
	LogsBloom         string          `json:"logsBloom"`
	Status            string          `json:"status"`
	To                string          `json:"to"`
	TransactionHash   string          `json:"transactionHash"`
	TransactionIndex  string          `json:"transactionIndex"`
	Type              string          `json:"type"`
}

type EthReceiptLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

// EthGetTransactionReceipt returns the transaction receipt by hash
// eth_getTransactionReceipt
func EthGetTransactionReceipt(ctx context.Context, isMainnet bool, hash string) (*EtherscanTransactionReceipt, error) {
	var endpoint = etherscanAPIBaseURL(isMainnet) + "?module=proxy&action=eth_getTransactionReceipt&txhash=" + hash + "&apikey=" + EtherscanAPIKey
	body, err := HttpGet(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	var txn EtherscanProxyRes[EtherscanTransactionReceipt]
	if err = json.Unmarshal(body, &txn); err != nil {
		return nil, err
	}
	return &txn.Result, nil
}

type EtherscanBlock struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
}

func EthGetBlockByNum(ctx context.Context, isMainnet bool, blockNum uint64) (*EtherscanBlock, error) {
	// https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x10d4f&boolean=false&apikey=YourApiKeyToken
	var endpoint = etherscanAPIBaseURL(isMainnet) + "?module=proxy&action=eth_getBlockByNumber&tag=0x" + fmt.Sprintf("%x", blockNum) + "&boolean=false&apikey=" + EtherscanAPIKey
	body, err := HttpGet(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	var block EtherscanProxyRes[EtherscanBlock]
	if err = json.Unmarshal(body, &block); err != nil {
		return nil, err
	}
	return &block.Result, nil
}

type EtherscanRes[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

func EtherscanGetBlockByTime(ctx context.Context, isMainnet bool, timestamp int64) (uint, error) {
	// https://api-sepolia.etherscan.io/api?module=block&action=getblocknobytime&timestamp=1708921620&closest=after&apikey=5D91FZ48V3XPNV58ZSNAYBBWS14K3GGSZC
	var endpoint = etherscanAPIBaseURL(isMainnet) + "?module=block&action=getblocknobytime&timestamp=" + fmt.Sprintf("%d", timestamp) + "&closest=after&apikey=" + EtherscanAPIKey
	body, err := HttpGet(ctx, endpoint)
	if err != nil {
		return 0, err
	}
	var block EtherscanRes[string]
	if err = json.Unmarshal(body, &block); err != nil {
		return 0, err
	}

	return ToUint(block.Result), nil
}

type EtherscanLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	TimeStamp        string   `json:"timeStamp"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

func EtherscanGetLogs(ctx context.Context, isMainnet bool, fromBlock uint64, address, topic0 string, page uint64, offset uint64) ([]EtherscanLog, error) {
	// https://api.etherscan.io/api?module=logs&action=getLogs&fromBlock=12878196&topic0=0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef&page=1&offset=1&apikey=YourApiKeyToken
	var endpoint = etherscanAPIBaseURL(isMainnet) + "?module=logs&action=getLogs&fromBlock=" + fmt.Sprintf("%d", fromBlock) + "&address=" + address + "&topic0=" + topic0 + "&page=" + fmt.Sprintf("%d", page) + "&offset=" + fmt.Sprintf("%d", offset) + "&apikey=" + EtherscanAPIKey
	body, err := HttpGet(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	var logs EtherscanRes[[]EtherscanLog]
	if err = json.Unmarshal(body, &logs); err != nil {
		return nil, err
	}
	return logs.Result, nil
}
