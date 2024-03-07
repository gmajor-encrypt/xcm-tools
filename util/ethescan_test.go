package util

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_EthGetTransactionByHash(t *testing.T) {
	type args struct {
		ctx  context.Context
		hash string
	}
	tests := []struct {
		args args
		want *EtherscanTransaction
		err  error
	}{
		{
			args: args{ctx: context.Background(), hash: "0x270b9592600015788b279df9eab62670349d006cf3ffaf3185f84163b37b9154"},
			want: &EtherscanTransaction{
				Hash:        "0x270b9592600015788b279df9eab62670349d006cf3ffaf3185f84163b37b9154",
				BlockNumber: "0x4d9dd9",
				From:        "0x62418364c0d732274759a64fd319e3f0ebf2f931",
			},
		},
		{
			args: args{ctx: context.Background(), hash: "0xaf7dfd2ae2f1d8bce96f935aa4b79a08b1040cf0c02fbd6b87d0957dba38651d"},
			want: &EtherscanTransaction{
				Hash:        "0xaf7dfd2ae2f1d8bce96f935aa4b79a08b1040cf0c02fbd6b87d0957dba38651d",
				BlockNumber: "0x51933a",
				From:        "0xbfe9debe43f8083bec64258d71113fad8632650a",
			},
		},
	}
	for _, tt := range tests {
		got, err := EthGetTransactionByHash(tt.args.ctx, tt.args.hash)
		if !errors.Is(err, tt.err) {
			t.Errorf("EthGetTransactionByHash() error = %v, wantErr %v", err, tt.err)
			return
		}
		if got != nil {
			assert.Equal(t, got.Hash, tt.want.Hash)
			assert.Equal(t, got.BlockNumber, tt.want.BlockNumber)
			assert.Equal(t, got.From, tt.want.From)
		}
	}
}

func Test_EthGetTransactionReceipt(t *testing.T) {
	type args struct {
		ctx  context.Context
		hash string
	}
	tests := []struct {
		args args
		want *EtherscanTransactionReceipt
		err  error
	}{
		{
			args: args{ctx: context.Background(), hash: "0x270b9592600015788b279df9eab62670349d006cf3ffaf3185f84163b37b9154"},
			want: &EtherscanTransactionReceipt{
				BlockNumber: "0x4d9dd9",
				From:        "0x62418364c0d732274759a64fd319e3f0ebf2f931",
				Status:      "0x1",
				Logs:        []EthReceiptLog{{Address: "0xfff9976782d46cc05630d1f6ebab18b2324d6b14"}},
			},
		},
		{
			args: args{ctx: context.Background(), hash: "0xaf7dfd2ae2f1d8bce96f935aa4b79a08b1040cf0c02fbd6b87d0957dba38651d"},
			want: &EtherscanTransactionReceipt{
				BlockNumber: "0x51933a",
				From:        "0xbfe9debe43f8083bec64258d71113fad8632650a",
				Status:      "0x1",
				Logs:        []EthReceiptLog{{Address: "0x679863b64072b7562c0fc7d8d831a6047681986a"}},
			},
		},
	}
	for _, tt := range tests {
		got, err := EthGetTransactionReceipt(tt.args.ctx, tt.args.hash)
		if !errors.Is(err, tt.err) {
			t.Errorf("EthGetTransactionByHash() error = %v, wantErr %v", err, tt.err)
			return
		}
		if got != nil {
			assert.Equal(t, got.BlockNumber, tt.want.BlockNumber)
			assert.Equal(t, got.From, tt.want.From)
			assert.Equal(t, got.Status, tt.want.Status)
			assert.Equal(t, got.Logs[0].Address, tt.want.Logs[0].Address)
		}
	}
}

func Test_EthGetBlockByNum(t *testing.T) {
	type args struct {
		ctx      context.Context
		blockNum uint64
	}
	tests := []struct {
		args args
		want *EtherscanBlock
		err  error
	}{
		{
			args: args{ctx: context.Background(), blockNum: 5364196},
			want: &EtherscanBlock{
				Hash:       "0x46c97d74e1a002361cb1930fbc5a20b9cde13308306543cdfcba51f161de70b5",
				Number:     "0x51d9e4",
				ParentHash: "0x5970a9ea94a2464bad21861fd5b6d48b79b21c25678d2cddd6f382a046c1e56c",
				Timestamp:  "0x65dbf70c",
			},
		},
		{
			args: args{ctx: context.Background(), blockNum: 4536419},
			want: &EtherscanBlock{
				Hash:       "0xd303268a2d5907c453c37fa6ea6c9ab40920c7e010191aad7222cddf92f8f6b0",
				Number:     "0x453863",
				ParentHash: "0xdc011e2787fcbd6acb4d6a2007a5e199dc3bd35ab1292828bef5a2a0933d8ea9",
				Timestamp:  "0x65346624",
			},
		},
	}
	for _, tt := range tests {
		time.Sleep(1 * time.Second)
		got, err := EthGetBlockByNum(tt.args.ctx, tt.args.blockNum)
		if !errors.Is(err, tt.err) {
			t.Errorf("Test_EthGetBlockByNum() error = %v, wantErr %v", err, tt.err)
			return
		}
		if got != nil {
			assert.Equal(t, got.Hash, tt.want.Hash)
			assert.Equal(t, got.Number, tt.want.Number)
			assert.Equal(t, got.ParentHash, tt.want.ParentHash)
			assert.Equal(t, got.Timestamp, tt.want.Timestamp)
		}
	}
}
