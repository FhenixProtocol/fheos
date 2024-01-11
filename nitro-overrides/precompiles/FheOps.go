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
	return fheos.Add(input, &tp, c.FheosState)
}

func (con FheOps) And(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.And(input, &tp, c.FheosState)
}

func (con FheOps) Cast(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Cast(input, &tp, c.FheosState)
}

func (con FheOps) Decrypt(c ctx, evm mech, input []byte) (*big.Int, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Decrypt(input, &tp, c.FheosState)
}

func (con FheOps) Div(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Div(input, &tp, c.FheosState)
}

func (con FheOps) Eq(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Eq(input, &tp, c.FheosState)
}

func (con FheOps) GetNetworkPublicKey(c ctx, evm mech) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.GetNetworkPublicKey(&tp, c.FheosState)
}

func (con FheOps) Gt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Gt(input, &tp, c.FheosState)
}

func (con FheOps) Gte(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Gte(input, &tp, c.FheosState)
}

func (con FheOps) Lt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Lt(input, &tp, c.FheosState)
}

func (con FheOps) Lte(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Lte(input, &tp, c.FheosState)
}

func (con FheOps) Max(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Max(input, &tp, c.FheosState)
}

func (con FheOps) Min(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Min(input, &tp, c.FheosState)
}

func (con FheOps) Mul(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Mul(input, &tp, c.FheosState)
}

func (con FheOps) Ne(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Ne(input, &tp, c.FheosState)
}

func (con FheOps) Not(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Not(input, &tp, c.FheosState)
}

func (con FheOps) Or(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Or(input, &tp, c.FheosState)
}

func (con FheOps) Rem(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Rem(input, &tp, c.FheosState)
}

func (con FheOps) Req(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Req(input, &tp, c.FheosState)
}

func (con FheOps) SealOutput(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.SealOutput(input, &tp, c.FheosState)
}

func (con FheOps) Select(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Select(input, &tp, c.FheosState)
}

func (con FheOps) Shl(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Shl(input, &tp, c.FheosState)
}

func (con FheOps) Shr(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Shr(input, &tp, c.FheosState)
}

func (con FheOps) Sub(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Sub(input, &tp, c.FheosState)
}

func (con FheOps) TrivialEncrypt(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.TrivialEncrypt(input, &tp, c.FheosState)
}

func (con FheOps) Verify(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Verify(input, &tp, c.FheosState)
}

func (con FheOps) Xor(c ctx, evm mech, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)
	return fheos.Xor(input, &tp, c.FheosState)
}
