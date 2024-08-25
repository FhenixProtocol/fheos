package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/fheos/precompiles"
	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
	"github.com/spf13/cobra"
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

func getenvInt(key string, defaultValue int) (int, error) {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func generateKeys(securityZones int32) error {
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

	for i := int32(0); i < securityZones; i++ {
		err := fhedriver.GenerateFheKeys(i)
		if err != nil {
			return fmt.Errorf("error generating FheKeys for securityZone %d: %s", i, err)
		}
	}
	return nil
}

func initDbOnly() error {
	err := precompiles.InitFheos(&fhedriver.ConfigDefault)
	if err != nil {
		return err
	}

	return nil
}

func initFheos() (*precompiles.TxParams, error) {
	if os.Getenv("FHEOS_DB_PATH") == "" {
		err := os.Setenv("FHEOS_DB_PATH", "./fheosdb")
		if err != nil {
			return nil, err
		}
	}

	err := precompiles.InitFheConfig(&fhedriver.ConfigDefault)
	if err != nil {
		return nil, err
	}

	securityZones, err := getenvInt("FHEOS_SECURITY_ZONES", 1)
	if err != nil {
		return nil, err
	}
	err = generateKeys(int32(securityZones))
	if err != nil {
		return nil, err
	}

	err = precompiles.InitializeFheosState()
	if err != nil {
		return nil, err
	}

	err = os.Setenv("FHEOS_DB_PATH", "")
	if err != nil {
		return nil, err
	}

	tp := precompiles.TxParams{
		Commit:          false,
		GasEstimation:   false,
		EthCall:         true,
		CiphertextDb:    nil,
		ContractAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GetBlockHash:    vm.GetHashFunc(nil),
		BlockNumber:     nil,
		ParallelTxHooks: nil,
	}

	return &tp, err
}

func getValue(a []byte) *big.Int {
	var value big.Int
	value.SetBytes(a[:32])

	return &value
}

func encrypt(val uint32, t uint8, securityZone int32, tp *precompiles.TxParams) ([]byte, error) {
	bval := new(big.Int).SetUint64(uint64(val))

	valBz := make([]byte, 32)
	bval.FillBytes(valBz)

	result, _, err := precompiles.TrivialEncrypt(valBz, t, securityZone, tp)
	return result, err
}

func decrypt(utype byte, val []byte, tp *precompiles.TxParams) (uint64, error) {
	decrypted, _, err := precompiles.Decrypt(utype, val, big.NewInt(0), tp)
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
	var securityZone int32
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

			elhs, err := encrypt(lhs, t, securityZone, txParams)
			if err != nil {
				return err
			}

			erhs, err := encrypt(rhs, t, securityZone, txParams)
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
	cmd.Flags().Int32VarP(&securityZone, "security-zone", "z", 0, "security zone")

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

	var add = setupOperationCommand("add", "add two numbers", precompiles.Add)
	var sub = setupOperationCommand("sub", "subtract two numbers", precompiles.Sub)
	var lte = setupOperationCommand("lte", "lte two numbers", precompiles.Lte)
	var mul = setupOperationCommand("mul", "mul two numbers", precompiles.Mul)
	var lt = setupOperationCommand("lt", "lt two numbers", precompiles.Lt)
	var div = setupOperationCommand("div", "div two numbers", precompiles.Div)
	var gt = setupOperationCommand("gt", "gt two numbers", precompiles.Gt)
	var gte = setupOperationCommand("gte", "gte two numbers", precompiles.Gte)
	var rem = setupOperationCommand("rem", "rem two numbers", precompiles.Rem)
	var and = setupOperationCommand("and", "and two numbers", precompiles.And)
	var or = setupOperationCommand("or", "or two numbers", precompiles.Or)
	var xor = setupOperationCommand("xor", "xor two numbers", precompiles.Xor)
	var eq = setupOperationCommand("eq", "eq two numbers", precompiles.Eq)
	var ne = setupOperationCommand("ne", "ne two numbers", precompiles.Ne)
	var min = setupOperationCommand("min", "min two numbers", precompiles.Min)
	var max = setupOperationCommand("max", "max two numbers", precompiles.Max)
	var shl = setupOperationCommand("shl", "shl two numbers", precompiles.Shl)
	var shr = setupOperationCommand("shr", "shr two numbers", precompiles.Shr)

	rootCmd.AddCommand(initDb, initState, add, sub, lte, sub, mul, lt, div, gt, gte, rem, and, or, xor, eq, ne, min, max, shl, shr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
