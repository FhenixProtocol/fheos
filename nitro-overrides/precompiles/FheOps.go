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
	v, g, e := fheos.Add(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) And(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.And(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Cast(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Cast(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Decrypt(c ctx, evm mech, input []byte) (*big.Int, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Decrypt(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Div(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Div(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Eq(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Eq(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) GetNetworkPublicKey(c ctx, evm mech) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.GetNetworkPublicKey(&tp, c.FheosState)
}

func (con FheOps) Gt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Gt(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Gte(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Gte(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Lt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Lt(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Lte(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Lte(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Max(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Max(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Min(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Min(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Mul(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Mul(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Ne(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Ne(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Not(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Not(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Or(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Or(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Rem(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Rem(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Req(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Req(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) SealOutput(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.SealOutput(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Select(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Select(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Shl(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Shl(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Shr(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Shr(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Sub(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Sub(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) TrivialEncrypt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.TrivialEncrypt(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Verify(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Verify(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}

func (con FheOps) Xor(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	v, g, e := fheos.Xor(input, &tp, c.FheosState)

	if e != nil {
		return v, e
	}

	e = c.Burn(g)
	return v, e
}
