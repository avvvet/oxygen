package kv

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Ledger struct {
	Db *leveldb.DB
}

func NewLedger(path string) (*Ledger, error) {
	// o := &opt.Options{
	// 	Filter: filter.NewBloomFilter(10),
	// }

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Ledger{db}, nil
}
