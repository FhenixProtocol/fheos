package main

import (
	"fmt"
	"github.com/fhenixprotocol/fheos/precompiles"
	"github.com/fhenixprotocol/go-tfhe"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"path/filepath"
)

type MockGasBurner struct{}

func (b MockGasBurner) Burn(_ uint64) error {
	return nil
}

func (b MockGasBurner) Burned() uint64 {
	return 0
}

func initFheosState() error {
	if os.Getenv("FHEOS_DB_PATH") == "" {
		err := os.Setenv("FHEOS_DB_PATH", "./fheosdb")
		if err != nil {
			return err
		}
	}

	err := precompiles.InitializeFheosState(MockGasBurner{})
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(os.Getenv("FHEOS_DB_PATH"), "LOCK"))
	if err != nil {
		return err
	}

	return nil
}

func removeDb() error {
	env := os.Getenv("FHEOS_DB_PATH")
	if env == "" {
		return nil
	}

	err := os.Remove(env)
	if err != nil {
		return err
	}

	return os.Setenv("FHEOS_DB_PATH", "")
}

func initConfig() error {
	return precompiles.InitTfheConfig(&tfhe.ConfigDefault)
}

func initLogger() {
	precompiles.InitLogger()
}

func generateKeys() error {
	if _, err := os.Stat("./keys/"); os.IsNotExist(err) {
		err := os.Mkdir("./keys/", 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("./keys/tfhe/"); os.IsNotExist(err) {
		err := os.Mkdir("./keys/tfhe/", 0755)
		if err != nil {
			return err
		}
	}

	err :=
		tfhe.GenerateFheKeys("./keys/tfhe/", "./sks", "./cks", "./pks")
	if err != nil {
		return fmt.Errorf("error from tfhe GenerateFheKeys: %s", err)
	}
	return nil
}

func initFheos() (*precompiles.TxParams, error) {
	initLogger()

	err := generateKeys()
	if err != nil {
		return nil, err
	}

	err = initConfig()
	if err != nil {
		return nil, err
	}

	err = initFheosState()
	if err != nil {
		return nil, err
	}

	tp := precompiles.TxParams{false, false, true}

	return &tp, err
}

func getValue(a []byte) *big.Int {
	var value big.Int
	value.SetBytes(a[:32])

	return &value
}

func encrypt(val uint32, t uint8, tp *precompiles.TxParams) ([]byte, error) {
	bval := new(big.Int).SetUint64(uint64(val))

	valBz := make([]byte, 32)
	bval.FillBytes(valBz)

	result, _, err := precompiles.TrivialEncrypt(valBz, t, tp)
	return result, err
}

func decrypt(utype byte, val []byte, tp *precompiles.TxParams) (uint64, error) {
	decrypted, _, err := precompiles.Decrypt(utype, val, tp)
	if err != nil {
		return 0, err
	}

	return decrypted.Uint64(), nil
}

func serialize2Params(lhs *big.Int, rhs *big.Int) []byte {
	val := make([]byte, 32)
	copy(val, lhs.Bytes())
	serialized := append(val, rhs.Bytes()...)
	return serialized
}

func main() {

	var rootCmd = &cobra.Command{Use: "fheos"}

	var initState = &cobra.Command{
		Use:   "init-state",
		Short: "Initialize fheos state",
		RunE: func(cmd *cobra.Command, args []string) error {
			initLogger()
			err := initFheosState()
			return err
		},
	}

	var lhs uint32
	var rhs uint32
	var t uint8
	var add = &cobra.Command{
		Use:   "add",
		Short: "add two numbers",
		RunE: func(cmd *cobra.Command, args []string) error {
			txParams, err := initFheos()
			defer removeDb()
			if err != nil {
				return err
			}
			elhs, err := encrypt(lhs, t, txParams)
			if err != nil {
				return err
			}
			erhs, err := encrypt(rhs, t, txParams)
			if err != nil {
				return err
			}

			result, _, err := precompiles.Add(t, elhs, erhs, txParams)
			if err != nil {
				return err
			}

			decrypted, err := decrypt(t, result, txParams)
			if err != nil {
				return err
			}

			dlhs, err := decrypt(t, elhs, txParams)
			if err != nil {
				return err
			}

			drhs, err := decrypt(t, erhs, txParams)
			if err != nil {
				return err
			}

			fmt.Printf("Added %d with %d and the result was %d\n", dlhs, drhs, decrypted)
			return nil
		},
	}
	add.Flags().Uint32VarP(&lhs, "lhs", "l", 0, "lhs")
	add.Flags().Uint32VarP(&rhs, "rhs", "r", 0, "rhs")
	add.Flags().Uint8VarP(&t, "utype", "t", 0, "uint type(0-uint8, 1-uint16, 2-uint32)")

	rootCmd.AddCommand(initState, add)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
