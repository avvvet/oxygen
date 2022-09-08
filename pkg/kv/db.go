package kv

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Ledger struct {
	Db *leveldb.DB
}

func NewLedger() (*Ledger, error) {
	db, err := leveldb.OpenFile("./ledger/store", nil)
	if err != nil {
		return nil, err
	}
	return &Ledger{db}, nil
}
