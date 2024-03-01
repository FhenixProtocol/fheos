package main

import (
	"fmt"
	"github.com/fhenixprotocol/go-tfhe"
	"github.com/spf13/cobra"
)

const VERSION = "0.0.5"

var versionCommand *cobra.Command

func init() {
	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("FheOs: " + GetFullVersionString())
			return nil
		},
	}
}

func GetFullVersionString() string {
	return VERSION + "\nTfhe-rs: " + tfhe.Version()
}
