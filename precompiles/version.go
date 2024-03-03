package precompiles

import "github.com/fhenixprotocol/go-tfhe"

const VERSION = "0.0.5"

func GetFullVersionString() string {
	return VERSION + "\nTfhe-rs: " + tfhe.Version()
}
