package arbos

import (
	"github.com/ethereum/go-ethereum/core/vm"
	fheos "github.com/fhenixprotocol/fheos/precompiles/types"
)

// PendingDecryption is an alias for fheos.PendingDecryption to simplify the code
type PendingDecryption = fheos.PendingDecryption

// Ensure ParallelTxProcessor implements ParallelTxProcessingHook
var _ fheos.ParallelTxProcessingHook = (*ParallelTxProcessor)(nil)

// Update ParallelTxProcessor to implement ParallelTxProcessingHook
type ParallelTxProcessor struct {
	vm.DefaultTxProcessor
	notifyCt         func(*PendingDecryption)
	notifyDecryptRes func(*PendingDecryption) error
}

// NewParallelTxProcessor creates a new ParallelTxProcessor with the provided notification functions
func NewParallelTxProcessor(
	notifyCt func(*PendingDecryption),
	notifyDecryptRes func(*PendingDecryption) error,
) *ParallelTxProcessor {
	return &ParallelTxProcessor{
		DefaultTxProcessor: vm.DefaultTxProcessor{},
		notifyCt:           notifyCt,
		notifyDecryptRes:   notifyDecryptRes,
	}
}

func (p *ParallelTxProcessor) NotifyCt(data *PendingDecryption) {
	if p.notifyCt != nil {
		p.notifyCt(data)
	}
}

func (p *ParallelTxProcessor) NotifyDecryptRes(data *PendingDecryption) error {
	if p.notifyDecryptRes != nil {
		return p.notifyDecryptRes(data)
	}
	return nil
}
