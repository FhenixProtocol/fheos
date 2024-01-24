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

func (con FheOps) Add(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Add")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Add", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Add(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Add/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Add/success/total/"
		if err != nil {
			metricPath = "/Add/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) And(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "And")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "And", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.And(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/And/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/And/success/total/"
		if err != nil {
			metricPath = "/And/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Cast(c ctx, evm mech, utype byte, input []byte, toType byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Cast")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Cast", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Cast(utype, input, toType, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Cast/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Cast/success/total/"
		if err != nil {
			metricPath = "/Cast/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Decrypt(c ctx, evm mech, utype byte, input []byte) (*big.Int, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Decrypt")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Decrypt", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Decrypt(utype, input, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Decrypt/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Decrypt/success/total/"
		if err != nil {
			metricPath = "/Decrypt/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Div(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Div")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Div", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Div(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Div/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Div/success/total/"
		if err != nil {
			metricPath = "/Div/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Eq(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Eq")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Eq", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Eq(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Eq/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Eq/success/total/"
		if err != nil {
			metricPath = "/Eq/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) GetNetworkPublicKey(c ctx, evm mech) ([]byte, error) {

	tp := fheos.TxParamsFromEVM(evm)
	return fheos.GetNetworkPublicKey(&tp)
}

func (con FheOps) Gt(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Gt")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Gt", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Gt(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Gt/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Gt/success/total/"
		if err != nil {
			metricPath = "/Gt/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Gte(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Gte")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Gte", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Gte(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Gte/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Gte/success/total/"
		if err != nil {
			metricPath = "/Gte/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Lt(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Lt")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Lt", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Lt(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Lt/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Lt/success/total/"
		if err != nil {
			metricPath = "/Lt/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Lte(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Lte")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Lte", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Lte(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Lte/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Lte/success/total/"
		if err != nil {
			metricPath = "/Lte/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Max(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Max")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Max", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Max(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Max/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Max/success/total/"
		if err != nil {
			metricPath = "/Max/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Min(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Min")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Min", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Min(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Min/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Min/success/total/"
		if err != nil {
			metricPath = "/Min/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Mul(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Mul")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Mul", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Mul(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Mul/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Mul/success/total/"
		if err != nil {
			metricPath = "/Mul/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Ne(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Ne")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Ne", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Ne(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Ne/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Ne/success/total/"
		if err != nil {
			metricPath = "/Ne/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Not(c ctx, evm mech, utype byte, value []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Not")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Not", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Not(utype, value, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Not/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Not/success/total/"
		if err != nil {
			metricPath = "/Not/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Or(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Or")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Or", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Or(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Or/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Or/success/total/"
		if err != nil {
			metricPath = "/Or/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Rem(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Rem")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Rem", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Rem(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Rem/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Rem/success/total/"
		if err != nil {
			metricPath = "/Rem/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Req(c ctx, evm mech, utype byte, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Req")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Req", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Req(utype, input, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Req/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Req/success/total/"
		if err != nil {
			metricPath = "/Req/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) SealOutput(c ctx, evm mech, utype byte, ctHash []byte, pk []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "SealOutput")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "SealOutput", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.SealOutput(utype, ctHash, pk, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/SealOutput/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/SealOutput/success/total/"
		if err != nil {
			metricPath = "/SealOutput/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Select(c ctx, evm mech, utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Select")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Select", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Select(utype, controlHash, ifTrueHash, ifFalseHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Select/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Select/success/total/"
		if err != nil {
			metricPath = "/Select/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Shl(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Shl")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Shl", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Shl(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Shl/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Shl/success/total/"
		if err != nil {
			metricPath = "/Shl/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Shr(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Shr")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Shr", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Shr(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Shr/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Shr/success/total/"
		if err != nil {
			metricPath = "/Shr/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Sub(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Sub")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Sub", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Sub(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Sub/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Sub/success/total/"
		if err != nil {
			metricPath = "/Sub/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) TrivialEncrypt(c ctx, evm mech, input []byte, toType byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "TrivialEncrypt")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "TrivialEncrypt", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.TrivialEncrypt(input, toType, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/TrivialEncrypt/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/TrivialEncrypt/success/total/"
		if err != nil {
			metricPath = "/TrivialEncrypt/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Verify(c ctx, evm mech, utype byte, input []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Verify")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Verify", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Verify(utype, input, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Verify/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Verify/success/total/"
		if err != nil {
			metricPath = "/Verify/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}

func (con FheOps) Xor(c ctx, evm mech, utype byte, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	tp := fheos.TxParamsFromEVM(evm)

	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s", "fheos", "Xor")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
			// fmt.Printf("FHEOS: %s took %d\n", "Xor", time.Since(start).Milliseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.Xor(utype, lhsHash, rhsHash, &tp)

	if err != nil {
		if metrics.Enabled {
			metrics.GetOrRegisterCounter("fheos"+"/Xor/error/fhe_failure/", nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := "/Xor/success/total/"
		if err != nil {
			metricPath = "/Xor/error/fhe_gas_failure/"
		}

		metrics.GetOrRegisterCounter("fheos"+metricPath, nil).Inc(1)
	}

	return ret, err
}
