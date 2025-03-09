package main

import (
	"math/big"

	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type FheOperationRequest struct {
	UType        byte                      `json:"uType"`
	Inputs       []fhedriver.CiphertextKey `json:"inputs"`
	RequesterUrl string                    `json:"requesterUrl"`
}

type TrivialEncryptRequest struct {
	Value        *big.Int `json:"value"`
	ToType       byte     `json:"toType"`
	SecurityZone int32    `json:"securityZone"`
	RequesterUrl string   `json:"requesterUrl"`
}

type CastRequest struct {
	UType        byte                    `json:"utype"`
	Input        fhedriver.CiphertextKey `json:"input"`
	ToType       string                  `json:"toType"`
	RequesterUrl string                  `json:"requesterUrl"`
}

type DecryptRequest struct {
	UType           byte                    `json:"utype"`
	Key             fhedriver.CiphertextKey `json:"key"`
	RequesterUrl    string                  `json:"requesterUrl"`
	TransactionHash string                  `json:"transactionHash"`
	ChainId         uint64                  `json:"chainId"`
}

type StoreCtsEntry struct {
	UType        byte   `json:"utype"`
	Value        string `json:"value"`
	SecurityZone byte   `json:"securityZone"`
}

type StoreCtsRequest struct {
	Cts       []StoreCtsEntry `json:"cts"`
	Signature string          `json:"signature"`
}

type GetNetworkPublicKeyRequest struct {
	SecurityZone byte `json:"securityZone"`
}

type GetCrsRequest struct {
	SecurityZone byte `json:"securityZone"`
}

type GetCTRequest struct {
	Hash string `json:"hash"`
}
