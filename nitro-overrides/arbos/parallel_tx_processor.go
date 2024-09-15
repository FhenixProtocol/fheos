package arbos

import (
	"github.com/ethereum/go-ethereum/core/vm"
	fheos "github.com/fhenixprotocol/fheos/precompiles/types"
)

// PendingDecryption is an alias for fheos.PendingDecryption to simplify the code
type PendingDecryption = fheos.PendingDecryption

// Ensure ParallelTxProcessor implements ParallelTxProcessingHook
var _ fheos.ParallelTxProcessingHook = (*ParallelTxProcessor)(nil)

// Update ParallelTxProcessor to implement ParallelTxProcessingHook and hold TxProcessingHook
type ParallelTxProcessor struct {
	vm.TxProcessingHook
	notifyCt         func(*PendingDecryption) error
	notifyDecryptRes func(*PendingDecryption) error
}

// NewParallelTxProcessor creates a new ParallelTxProcessor with the provided TxProcessingHook and notification functions
func NewParallelTxProcessor(
	txProcessingHook vm.TxProcessingHook,
	notifyCt func(*PendingDecryption) error,
	notifyDecryptRes func(*PendingDecryption) error,
) *ParallelTxProcessor {
	return &ParallelTxProcessor{
		TxProcessingHook: txProcessingHook,
		notifyCt:         notifyCt,
		notifyDecryptRes: notifyDecryptRes,
	}
}

func (p *ParallelTxProcessor) NotifyCt(data *PendingDecryption) error {
	if p.notifyCt != nil {
		return p.notifyCt(data)
	}
	return nil
}

func (p *ParallelTxProcessor) NotifyDecryptRes(data *PendingDecryption) error {
	if p.notifyDecryptRes != nil {
		return p.notifyDecryptRes(data)
	}
	return nil
}
