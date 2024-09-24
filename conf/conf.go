package conf

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	flag "github.com/spf13/pflag"
)

func FheosConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.String(prefix+".fheos-db-path", ConfigDefault.FheosDbPath, "Path for FheOs database")
	f.Int32(prefix+".security-zones", ConfigDefault.SecurityZones, "Amount of security Zones to load")
}

type FheosConfig struct {
	FheosDbPath   string `koanf:"fheos-db-path"`
	SecurityZones int32  `koanf:"security-zones"`
}

var ConfigDefault = FheosConfig{
	FheosDbPath:   "/tmp/fheos",
	SecurityZones: 1,
}

var loadedConfig FheosConfig

func GetConfig() *FheosConfig {
	return &loadedConfig
}

func SetConfig(config FheosConfig) {
	log.Info(fmt.Sprintf("Setting Fheos config: %+v", config))
	loadedConfig = config
}
