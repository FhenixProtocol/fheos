package arbitrum

import (
	"flag"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

func FhenixConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.Bool(prefix+".is-oracle", fhe.ConfigDefault.IsOracle, "Is FHE oracle node?")
	f.String(prefix+".oracle-type", fhe.ConfigDefault.OracleType, "FHE oracle type")
	f.String(prefix+".oracle-db-path", fhe.ConfigDefault.OracleDbPath, "FHE oracle DB path")
	f.String(prefix+".oracle-address", fhe.ConfigDefault.OracleAddress, "FHE oracle address")
	f.String(prefix+".server-key-path", fhe.ConfigDefault.ServerKeyPath, "FHE server key path")
	f.String(prefix+".client-key-path", fhe.ConfigDefault.ClientKeyPath, "FHE client key path")
	f.String(prefix+".public-key-path", fhe.ConfigDefault.PublicKeyPath, "FHE public key path")
	f.String(prefix+".oracle-private-key-path", fhe.ConfigDefault.OraclePrivateKeyPath, "FHE oracle private key path")
	f.String(prefix+".oracle-public-key-path", fhe.ConfigDefault.OraclePublicKeyPath, "FHE oracle public key path")
	f.String(prefix+".home-dir", fhe.ConfigDefault.HomeDir, "FHE home directory")
}

type FhenixConfig = fhe.Config
