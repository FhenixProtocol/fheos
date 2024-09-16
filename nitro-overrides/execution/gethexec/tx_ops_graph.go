package gethexec

import (
	"time"

	"github.com/ethereum/go-ethereum/arbitrum_types"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	fheos "github.com/fhenixprotocol/fheos/precompiles/types"
)

type transaction struct {
	Tx        *geth.Transaction
	TxOptions *arbitrum_types.ConditionalOptions
	QueueItem *txQueueItem
	Retries   uint
	LastRun   time.Time

	pendingDecryptions  map[fheos.PendingDecryption]*decryption
	resolvedDecryptions map[fheos.PendingDecryption]*decryption
}

type decryption struct {
	transactions map[common.Hash]*transaction
}

type TxOpsGraph struct {
	transactions map[common.Hash]*transaction
	decryptions  map[fheos.PendingDecryption]*decryption
}

func NewTxOpsGraph() *TxOpsGraph {
	return &TxOpsGraph{
		transactions: make(map[common.Hash]*transaction),
		decryptions:  make(map[fheos.PendingDecryption]*decryption),
	}
}

func (bg *TxOpsGraph) AddEdge(gethtx *geth.Transaction, txOptions *arbitrum_types.ConditionalOptions, decryptKey fheos.PendingDecryption) {
	tx := bg.addTransaction(gethtx, txOptions)
	decrypt := bg.addDecryption(decryptKey)

	tx.pendingDecryptions[decryptKey] = decrypt
	decrypt.transactions[gethtx.Hash()] = tx
}

func (bg *TxOpsGraph) ResolveTransaction(txHash common.Hash) {
	bg.removeTransaction(txHash)
}

func (bg *TxOpsGraph) ResolveDecryption(decryptKey fheos.PendingDecryption) []*transaction {
	decrypt, exists := bg.decryptions[decryptKey]
	if !exists {
		return nil
	}

	pushedTransactions := []*transaction{}

	for _, tx := range decrypt.transactions {
		delete(tx.pendingDecryptions, decryptKey)
		tx.resolvedDecryptions[decryptKey] = decrypt

		if len(tx.pendingDecryptions) == 0 {
			pushedTransactions = append(pushedTransactions, tx)
		}
	}

	return pushedTransactions
}

func (bg *TxOpsGraph) SetTxQueueItem(queueItem txQueueItem) {
	txHash := queueItem.tx.Hash()
	tx, exists := bg.transactions[txHash]
	if !exists {
		return
	}
	tx.QueueItem = &queueItem
}

// Private helper methods

func (bg *TxOpsGraph) addTransaction(gethtx *geth.Transaction, txOptions *arbitrum_types.ConditionalOptions) *transaction {
	txHash := gethtx.Hash()
	if tx, exists := bg.transactions[txHash]; exists {
		return tx
	}
	tx := &transaction{
		Tx:                  gethtx,
		TxOptions:           txOptions,
		pendingDecryptions:  make(map[fheos.PendingDecryption]*decryption),
		resolvedDecryptions: make(map[fheos.PendingDecryption]*decryption),
	}
	bg.transactions[txHash] = tx
	return tx
}

func (bg *TxOpsGraph) addDecryption(decryptKey fheos.PendingDecryption) *decryption {
	if decrypt, exists := bg.decryptions[decryptKey]; exists {
		return decrypt
	}
	decrypt := &decryption{
		transactions: make(map[common.Hash]*transaction),
	}
	bg.decryptions[decryptKey] = decrypt
	return decrypt
}

func (bg *TxOpsGraph) removeTransaction(txHash common.Hash) {
	tx, exists := bg.transactions[txHash]
	if !exists {
		return
	}

	for decryptKey, decrypt := range tx.pendingDecryptions {
		delete(decrypt.transactions, txHash)
		if len(decrypt.transactions) == 0 {
			delete(bg.decryptions, decryptKey)
		}
	}
	for decryptKey, decrypt := range tx.resolvedDecryptions {
		delete(decrypt.transactions, txHash)
		if len(decrypt.transactions) == 0 {
			delete(bg.decryptions, decryptKey)
		}
	}

	delete(bg.transactions, txHash)
}
