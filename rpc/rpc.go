package rpc

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type ChainPublicKeyAPI struct {
}

// Get the public key of the chain
func (s *ChainPublicKeyAPI) GetNetworkPublicKey() (string, error) {
	expectedPk, err := fhe.PublicKey(0)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(expectedPk), nil
}

func GetRpcApis() []rpc.API {
	var apis []rpc.API
	apis = append(apis, rpc.API{
		Namespace: "chainpublickey",
		Version:   "1.0",
		Service:   &ChainPublicKeyAPI{},
		Public:    true,
	})

	return apis
}
