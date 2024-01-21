package precompiles

import (
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	"math/big"
)

type FheOps struct {
	Address addr // 0x80
}

func (con FheOps) Add(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Add(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) And(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.And(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Cast(c ctx, evm mech, utype byte, input []byte, toType byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Cast(utype, input, toType, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Decrypt(c ctx, evm mech, utype byte, input []byte) (*big.Int, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Decrypt(utype, input, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Div(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Div(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Eq(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Eq(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) GetNetworkPublicKey(c ctx, evm mech) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.GetNetworkPublicKey(&tp)
}

func (con FheOps) Gt(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Gt(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Gte(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Gte(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Lt(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Lt(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Lte(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Lte(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Max(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Max(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Min(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Min(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Mul(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Mul(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Ne(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Ne(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Not(c ctx, evm mech, utype byte, value []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Not(utype, value, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Or(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Or(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Rem(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Rem(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Req(c ctx, evm mech, utype byte, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Req(utype, input, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) SealOutput(c ctx, evm mech, utype byte, ctHash []byte, pk []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.SealOutput(utype, ctHash, pk, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Select(c ctx, evm mech, utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Select(utype, controlHash, ifTrueHash, ifFalseHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Shl(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Shl(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Shr(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Shr(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Sub(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Sub(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) TrivialEncrypt(c ctx, evm mech, input []byte, toType byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.TrivialEncrypt(input, toType, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Verify(c ctx, evm mech, utype byte, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Verify(utype, input, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}

func (con FheOps) Xor(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	ret, gas, err := fheos.Xor(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		return ret, err
	}

	err = c.Burn(gas)
	return ret, err
}
