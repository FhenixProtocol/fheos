package precompiles

import (
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	"math/big"
)

type FheOps struct {
	Address addr // 0x80
}

func (con FheOps) Add(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Add(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) And(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.And(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Cast(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Cast(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Decrypt(c ctx, evm mech, input []byte) (*big.Int, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Decrypt(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Div(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Div(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Eq(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Eq(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) GetNetworkPublicKey(c ctx, evm mech) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.GetNetworkPublicKey(&tp, c.FheosState)
}

func (con FheOps) Gt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Gt(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Gte(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Gte(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Lt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Lt(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Lte(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Lte(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Max(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Max(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Min(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Min(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Mul(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Mul(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Ne(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Ne(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Not(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Not(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Or(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Or(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Rem(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Rem(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Req(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Req(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) SealOutput(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.SealOutput(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Select(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Select(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Shl(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Shl(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Shr(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Shr(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Sub(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Sub(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) TrivialEncrypt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.TrivialEncrypt(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Verify(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Verify(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Xor(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Xor(input, &tp, c.FheosState)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}
