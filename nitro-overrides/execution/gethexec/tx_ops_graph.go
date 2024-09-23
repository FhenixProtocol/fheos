package gethexec

import (
	"errors"
	"sync"
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
	mu sync.Mutex

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
	bg.mu.Lock()
	defer bg.mu.Unlock()

	tx := bg.addTransaction(gethtx, txOptions)
	decrypt := bg.addDecryption(decryptKey)

	tx.pendingDecryptions[decryptKey] = decrypt
	decrypt.transactions[gethtx.Hash()] = tx
}

func (bg *TxOpsGraph) ResolveTransaction(txHash common.Hash) {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	bg.removeTransaction(txHash)
}

func (bg *TxOpsGraph) ResolveDecryption(decryptKey fheos.PendingDecryption) []*transaction {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	decrypt, exists := bg.decryptions[decryptKey]
	if !exists {
		return nil
	}

	pushedTransactions := []*transaction{}

	for _, tx := range decrypt.transactions {
		delete(tx.pendingDecryptions, decryptKey)
		tx.resolvedDecryptions[decryptKey] = decrypt

		if len(tx.pendingDecryptions) == 0 {
			if tx.Retries > 5 /* TODO this number is just a placeholder */ {
				bg.rejectTransaction(tx.Tx.Hash(), errors.New("transaction rejected after 5 retries"))
			} else {
				pushedTransactions = append(pushedTransactions, tx)
			}
		}
	}

	return pushedTransactions
}

func (bg *TxOpsGraph) SetTxQueueItem(queueItem txQueueItem) {
	bg.mu.Lock()
	defer bg.mu.Unlock()

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

func (bg *TxOpsGraph) rejectTransaction(txHash common.Hash, customErr error) error {
	tx, exists := bg.transactions[txHash]
	if !exists {
		return nil // Transaction not found, nothing to reject
	}

	// Remove the transaction from the graph
	bg.removeTransaction(txHash)

	// Check if the queueItem is still valid
	if tx.QueueItem != nil && tx.QueueItem.ctx.Err() == nil {
		tx.QueueItem.returnResult(customErr)
	}

	return nil
}
