// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package precompiles

// FHENIX: this file can't be moved to "Arbitrum" package as it will cause import cycle

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type GasBurner interface {
	Burn(amount uint64) error
	Burned() uint64
}

func uintToHash(val uint64) common.Hash {
	return common.BigToHash(new(big.Int).SetUint64(val))
}

// FHENIX: consider changing those values to be lower as for every 32bytes set it costs "SstoreSetGasEIP2200"
func writeCost(value common.Hash) uint64 {
	if value == (common.Hash{}) {
		return params.SstoreResetGasEIP2200
	}
	return params.SstoreSetGasEIP2200
}

// Storage allows Fheos to store data persistently in the Ethereum-compatible stateDB. This is represented in
// the stateDB as the storage of a fictional Ethereum account at address 0xA4B06FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF.
type Storage struct {
	account    common.Address
	db         vm.StateDB
	storageKey []byte
	burner     GasBurner
}

func NewStorage(statedb vm.StateDB, burner GasBurner) *Storage {
	account := common.HexToAddress("0xA4B06FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	statedb.SetNonce(account, 1) // setting the nonce ensures Geth won't treat the storage as empty
	return &Storage{
		account: account,
		db:      statedb,
		burner:  burner,
	}
}

// Arbitrum map addresses using "pages" of 256 storage slots. They hash over the page number but not the offset within
// a page, to preserve contiguity within a page. This will reduce cost if/when Ethereum switches to storage
// representations that reward contiguity.
// Because page numbers are 248 bits, this gives us 124-bit security against collision attacks, which is good enough.
func mapAddress(storageKey []byte, key common.Hash) common.Hash {
	keyBytes := key.Bytes()
	boundary := common.HashLength - 1
	return common.BytesToHash(
		append(
			crypto.Keccak256(storageKey, keyBytes[:boundary])[:boundary],
			keyBytes[boundary],
		),
	)
}

func (store *Storage) OpenSubStorage(id []byte) *Storage {
	return &Storage{
		store.account,
		store.db,
		crypto.Keccak256(store.storageKey, id),
		store.burner,
	}
}

func (store *Storage) Set(key common.Hash, value common.Hash) error {
	err := store.burner.Burn(params.SstoreSetGasEIP2200)
	if err != nil {
		return err
	}

	store.db.SetState(store.account, mapAddress(store.storageKey, key), value)
	return nil
}

func (store *Storage) SetByUint64(key uint64, value common.Hash) error {
	return store.Set(uintToHash(key), value)
}

func (store *Storage) SetUint64ByUint64(key uint64, value uint64) error {
	return store.Set(uintToHash(key), uintToHash(value))
}

func (store *Storage) Get(key common.Hash) (common.Hash, error) {
	err := store.burner.Burn(params.SloadGasEIP2200)
	if err != nil {
		return common.Hash{}, err
	}

	return store.db.GetState(store.account, mapAddress(store.storageKey, key)), nil
}

func (store *Storage) GetUint64(key common.Hash) (uint64, error) {
	value, err := store.Get(key)
	return value.Big().Uint64(), err
}

func (store *Storage) GetByUint64(key uint64) (common.Hash, error) {
	return store.Get(uintToHash(key))
}

func (store *Storage) GetUint64ByUint64(key uint64) (uint64, error) {
	return store.GetUint64(uintToHash(key))
}

func (store *Storage) GetBytes() ([]byte, error) {
	bytesLeft, err := store.GetUint64ByUint64(0)
	if err != nil {
		return nil, err
	}
	ret := []byte{}
	offset := uint64(1)
	for bytesLeft >= 32 {
		next, err := store.GetByUint64(offset)
		if err != nil {
			return nil, err
		}
		ret = append(ret, next.Bytes()...)
		bytesLeft -= 32
		offset++
	}
	next, err := store.GetByUint64(offset)
	if err != nil {
		return nil, err
	}
	ret = append(ret, next.Bytes()[32-bytesLeft:]...)
	return ret, nil
}

func (store *Storage) GetBytesSize() (uint64, error) {
	return store.GetUint64ByUint64(0)
}

func (store *Storage) ClearByUint64(key uint64) error {
	return store.Set(uintToHash(key), common.Hash{})
}

func (store *Storage) ClearBytes() error {
	bytesLeft, err := store.GetUint64ByUint64(0)
	if err != nil {
		return err
	}
	offset := uint64(1)
	for bytesLeft > 0 {
		err := store.ClearByUint64(offset)
		if err != nil {
			return err
		}
		offset++
		if bytesLeft < 32 {
			bytesLeft = 0
		} else {
			bytesLeft -= 32
		}
	}
	return store.ClearByUint64(0)
}

func (store *Storage) SetBytes(b []byte) error {
	err := store.ClearBytes()
	if err != nil {
		return err
	}
	err = store.SetUint64ByUint64(0, uint64(len(b)))
	if err != nil {
		return err
	}
	offset := uint64(1)
	for len(b) >= 32 {
		err = store.SetByUint64(offset, common.BytesToHash(b[:32]))
		if err != nil {
			return err
		}
		b = b[32:]
		offset++
	}
	return store.SetByUint64(offset, common.BytesToHash(b))
}

type BytesStorage struct {
	Storage
}

func (store *Storage) OpenBytesStorage(id []byte) BytesStorage {
	return BytesStorage{
		*store.OpenSubStorage(id),
	}
}

func (bs *BytesStorage) Get() ([]byte, error) {
	return bs.Storage.GetBytes()
}

func (bs *BytesStorage) Set(val []byte) error {
	return bs.Storage.SetBytes(val)
}
func (bs *BytesStorage) Size() (uint64, error) {
	return bs.Storage.GetBytesSize()
}

type SingleSlotStorage struct {
	account common.Address
	db      vm.StateDB
	slot    common.Hash
	burner  GasBurner
}

func (store *Storage) NewSlot(offset uint64) SingleSlotStorage {
	return SingleSlotStorage{store.account, store.db, mapAddress(store.storageKey, uintToHash(offset)), store.burner}
}

func (sss *SingleSlotStorage) Get() (common.Hash, error) {
	err := sss.burner.Burn(params.SloadGasEIP2200)
	if err != nil {
		return common.Hash{}, err
	}

	return sss.db.GetState(sss.account, sss.slot), nil
}

func (sss *SingleSlotStorage) Set(value common.Hash) error {
	err := sss.burner.Burn(writeCost(value))
	if err != nil {
		return err
	}

	sss.db.SetState(sss.account, sss.slot, value)
	return nil
}

type UintStorage struct {
	SingleSlotStorage
}

func (store *Storage) OpenUintStorage(offset uint64) UintStorage {
	return UintStorage{store.NewSlot(offset)}
}

func (us *UintStorage) Get() (uint64, error) {
	raw, err := us.SingleSlotStorage.Get()
	if !raw.Big().IsUint64() {
		panic("expected uint64 compatible value in storage")
	}
	return raw.Big().Uint64(), err
}

func (us *UintStorage) Set(value uint64) error {
	bigValue := new(big.Int).SetUint64(value)
	return us.SingleSlotStorage.Set(common.BigToHash(bigValue))
}
