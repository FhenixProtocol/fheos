package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fhenixprotocol/fheos/chains/arbitrum"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("bad path")
	}
	parent := filepath.Dir(filename)

	if os.Args[1] == "1" {
		arbitrum.CreateTemplate(parent)
		return
	}

	arbitrum.Gen(parent, os.Args[2])
}
