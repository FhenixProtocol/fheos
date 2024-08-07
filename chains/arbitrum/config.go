package arbitrum

import (
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	flag "github.com/spf13/pflag"
)

func FhenixConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.String(prefix+".oracle-type", fhe.ConfigDefault.OracleType, "FHE oracle type")
	f.String(prefix+".oracle-address", fhe.ConfigDefault.OracleAddress, "FHE oracle address")
	f.String(prefix+".fhe-engine-address", fhe.ConfigDefault.FheEngineAddress, "FHE engine address")
	f.String(prefix+".fallback-fhe-engine-address", fhe.ConfigDefault.FallbackFheEngineAddress, "FHE engine fallback address")
	f.String(prefix+".home-dir", fhe.ConfigDefault.HomeDir, "FHE home directory")
	f.Int(prefix+".log-level", fhe.ConfigDefault.LogLevel, "Minimum log level to display (0-5)")
}

type FhenixConfig = fhe.Config

var ConfigDefault = fhe.ConfigDefault
