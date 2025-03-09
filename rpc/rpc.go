package rpc

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type FhenixAPI struct {
}

// Get the public key of the chain
func (s *FhenixAPI) GetNetworkPublicKey() (string, error) {
	expectedPk, err := fhe.PublicKey(0)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(expectedPk), nil
}

// Get the CRS of the chain
func (s *FhenixAPI) GetCrs() (string, error) {
	crs, err := fhe.GetCrs(0)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(crs), nil
}

func GetRpcApis() rpc.API {
	return rpc.API{
		Namespace: "fhenix",
		Version:   "1.0",
		Service:   &FhenixAPI{},
		Public:    true,
	}
}
