package precompiles

import (
	"fmt"
	"github.com/ethereum/go-ethereum/metrics"
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	"math/big"
	"time"
)

type FheOps struct {
	Address addr // 0x80
}

func (con FheOps) Log(c ctx, evm mech, s string) error {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	gas, err := fheos.Log(s, &tp)

	if err != nil {
		return err
	}

	return c.Burn(gas)
}

func (con FheOps) Add(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Add", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Add(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Add", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Add", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Add", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) And(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "And", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.And(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "And", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "And", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "And", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Cast(c ctx, evm mech, utype byte, input []byte, toType byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Cast", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Cast(utype, input, toType, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Cast", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Cast", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Cast", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Decrypt(c ctx, evm mech, utype byte, input []byte) (*big.Int, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Decrypt", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Decrypt(utype, input, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Decrypt", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Decrypt", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Decrypt", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Div(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Div", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Div(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Div", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Div", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Div", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Eq(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Eq", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Eq(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Eq", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Eq", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Eq", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) GetNetworkPublicKey(c ctx, evm mech, securityZone int32) ([]byte, error) {

	tp := fheos.TxParamsFromEVM(evm, c.caller)
	return fheos.GetNetworkPublicKey(securityZone, &tp)
}

func (con FheOps) Gt(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Gt", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Gt(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Gt", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Gt", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Gt", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Gte(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Gte", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Gte(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Gte", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Gte", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Gte", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Lt(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Lt", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Lt(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Lt", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Lt", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Lt", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Lte(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Lte", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Lte(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Lte", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Lte", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Lte", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Max(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Max", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Max(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Max", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Max", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Max", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Min(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Min", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Min(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Min", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Min", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Min", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Mul(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Mul", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Mul(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Mul", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Mul", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Mul", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Ne(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Ne", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Ne(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Ne", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Ne", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Ne", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Not(c ctx, evm mech, utype byte, value []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Not", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Not(utype, value, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Not", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Not", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Not", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Or(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Or", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Or(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Or", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Or", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Or", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Random(c ctx, evm mech, utype byte, seed uint64) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Random", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Random(utype, seed, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Random", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Random", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Random", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Rem(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Rem", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Rem(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Rem", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Rem", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Rem", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Req(c ctx, evm mech, utype byte, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Req", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Req(utype, input, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Req", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Req", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Req", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) SealOutput(c ctx, evm mech, utype byte, ctHash []byte, pk []byte) (string, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "SealOutput", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.SealOutput(utype, ctHash, pk, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "SealOutput", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "SealOutput", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "SealOutput", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Select(c ctx, evm mech, utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Select", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Select(utype, controlHash, ifTrueHash, ifFalseHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Select", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Select", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Select", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Shl(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Shl", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Shl(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Shl", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Shl", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Shl", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Shr(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Shr", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Shr(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Shr", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Shr", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Shr", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Sub(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Sub", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Sub(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Sub", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Sub", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Sub", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) TrivialEncrypt(c ctx, evm mech, input []byte, toType byte, securityZone int32) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "TrivialEncrypt", fheos.UtypeToString(toType))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.TrivialEncrypt(input, toType, securityZone, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "TrivialEncrypt", fheos.UtypeToString(toType), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "TrivialEncrypt", fheos.UtypeToString(toType), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "TrivialEncrypt", fheos.UtypeToString(toType), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Verify(c ctx, evm mech, utype byte, input []byte, securityZone int32) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Verify", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Verify(utype, input, securityZone, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Verify", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Verify", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Verify", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Xor(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "Xor", fheos.UtypeToString(utype))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Xor(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Xor", fheos.UtypeToString(utype), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "Xor", fheos.UtypeToString(utype), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "Xor", fheos.UtypeToString(utype), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}
