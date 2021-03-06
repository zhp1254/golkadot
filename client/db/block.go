package clientdb

import (
	"github.com/opennetsys/golkadot/common/db"
	types "github.com/opennetsys/golkadot/types"
)

func createU8a(dbs db.BaseDB, fn types.StorageFunction) StorageMethodU8a {
	return NewStorageMethodU8a(dbs, fn)
}

func createBn(dbs db.BaseDB, createKey types.StorageFunction, bitLen int) StorageMethodBn {

	return NewStorageMethodBn(dbs, createKey, bitLen)
}

// NewBlockDB ...
func NewBlockDB(dbs db.BaseDB) *BlockDB {
	return &BlockDB{
		DB:         dbs,
		BestHash:   createU8a(dbs, KeyBestHash()),
		BestNumber: createBn(dbs, KeyBestNumber(), 64),
		BlockData:  createU8a(dbs, KeyBlockByHash()),
		Header:     createU8a(dbs, KeyHeaderByHash()),
		Hash:       createU8a(dbs, KeyHashByNumber()),
	}
}
