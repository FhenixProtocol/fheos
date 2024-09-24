package conf

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	flag "github.com/spf13/pflag"
)

func FheosConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.String(prefix+".fheos-db-path", ConfigDefault.FheosDbPath, "Path for FheOs database")
}

type FheosConfig struct {
	FheosDbPath string `koanf:"fheos-db-path"`
}

var ConfigDefault = FheosConfig{
	FheosDbPath: "/tmp/fheos",
}

var loadedConfig FheosConfig

func GetConfig() *FheosConfig {
	return &loadedConfig
}

func SetConfig(config FheosConfig) {
	log.Info(fmt.Sprintf("Setting Fheos config: %+v", config))
	loadedConfig = config
}
