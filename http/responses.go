package main

type FheOperationResponse struct {
	HadError   bool   `json:"hadError"`
	TempKey    []byte `json:"tempKey"`
	ActualHash []byte `json:"actualHash"`
}

type DecryptResponse struct {
	HadError        bool   `json:"hadError"`
	CtHash          []byte `json:"ctHash"`
	Plaintext       string `json:"plaintext"`
	TransactionHash string `json:"transactionHash"`
}

type StoreCtsResponse struct {
	Hashes []string `json:"hashes"`
}
type GetNetworkPublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

type GetCrsResponse struct {
	Crs string `json:"crs"`
}

type GetCTResponse struct {
	Data         string `json:"data"`
	SecurityZone uint8  `json:"security_zone"`
	UintType     uint8  `json:"uint_type"`
	Compact      bool   `json:"compact"`
	Gzipped      bool   `json:"gzipped"`
}
