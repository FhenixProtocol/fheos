package main

import (
	"fmt"
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	"github.com/spf13/cobra"
)

var versionCommand *cobra.Command

func init() {
	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("FheOs: " + fheos.GetFullVersionString())
			return nil
		},
	}
}
