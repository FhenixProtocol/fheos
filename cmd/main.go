package main

import (
	"fmt"
	"github.com/fhenixprotocol/fheos/precompiles"
	"github.com/fhenixprotocol/go-tfhe"
	"github.com/spf13/cobra"
	"math/big"
	"os"
)

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

func initDbOnly() error {
	err := precompiles.InitFheos(&tfhe.ConfigDefault)
	if err != nil {
		return err
	}

	return nil
}

func initFheos() (*precompiles.TxParams, error) {
	err := generateKeys()
	if err != nil {
		return nil, err
	}

	if os.Getenv("FHEOS_DB_PATH") == "" {
		err := os.Setenv("FHEOS_DB_PATH", "./fheosdb")
		if err != nil {
			return nil, err
		}
	}

	err = precompiles.InitFheos(&tfhe.ConfigDefault)
	if err != nil {
		return nil, err
	}

	err = os.Setenv("FHEOS_DB_PATH", "")
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

type operationFunc func(t byte, lhs, rhs []byte, txParams *precompiles.TxParams) ([]byte, uint64, error)

func setupOperationCommand(use, short string, op operationFunc) *cobra.Command {
	var lhs, rhs uint32
	var t uint8

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
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

			result, _, err := op(t, elhs, erhs, txParams)
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

			action := "operated" // You can choose a better word based on your operations
			fmt.Printf("%s (%+v aka %d) with (%+v aka %d) and the result was (%+v aka %d)\n", action, elhs, dlhs, erhs, drhs, result, decrypted)
			return nil
		},
	}

	cmd.Flags().Uint32VarP(&lhs, "lhs", "l", 0, "lhs")
	cmd.Flags().Uint32VarP(&rhs, "rhs", "r", 0, "rhs")
	cmd.Flags().Uint8VarP(&t, "utype", "t", 0, "uint type(0-uint8, 1-uint16, 2-uint32)")

	return cmd
}

func main() {
	var rootCmd = &cobra.Command{Use: "fheos"}

	var initState = &cobra.Command{
		Use:   "init-state",
		Short: "Initialize fheos state",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := initFheos()
			return err
		},
	}

	var initDb = &cobra.Command{
		Use:   "init-db",
		Short: "Initialize fheos db only (no keys)",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := initDbOnly()
			return err
		},
	}

	rootCmd.AddCommand(initDb, initState)

	var add = setupOperationCommand("add", "add two numbers", precompiles.Add)
	var sub = setupOperationCommand("sub", "subtract two numbers", precompiles.Sub)
	var lte = setupOperationCommand("lte", "perform less-than-or-equal comparison between two numbers", precompiles.Lte)
	var mul = setupOperationCommand("mul", "multiply two numbers", precompiles.Mul)
	var lt = setupOperationCommand("lt", "perform less-than comparison between two numbers", precompiles.Lt)
	var div = setupOperationCommand("div", "divide two numbers", precompiles.Div)
	var gt = setupOperationCommand("gt", "perform greater-than comparison between two numbers", precompiles.Gt)
	var gte = setupOperationCommand("gte", "perform greater-than-or-equal comparison between two numbers", precompiles.Gte)
	var rem = setupOperationCommand("rem", "get the remainder of a division between two numbers", precompiles.Rem)
	var and = setupOperationCommand("and", "bitwise-and two numbers", precompiles.And)
	var or = setupOperationCommand("or", "bitwise-or two numbers", precompiles.Or)
	var xor = setupOperationCommand("xor", "bitwise-xor two numbers", precompiles.Xor)
	var eq = setupOperationCommand("eq", "perform equality comparison between two numbers", precompiles.Eq)
	var ne = setupOperationCommand("ne", "perform inequality comparison between two numbers", precompiles.Ne)
	var min = setupOperationCommand("min", "return the smaller of two numbers", precompiles.Min)
	var max = setupOperationCommand("max", "return the greater of two numbers", precompiles.Max)
	var shl = setupOperationCommand("shl", "shift-left a number", precompiles.Shl)
	var shr = setupOperationCommand("shr", "shift-right a number", precompiles.Shr)

	rootCmd.AddCommand(add, sub, lte, sub, mul, lt, div, gt, gte, rem, and, or, xor, eq, ne, min, max, shl, shr)
	rootCmd.AddCommand(versionCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
