package kv

import "errors"

func (l *Ledger) Upsert(key []byte, value []byte) error {
	if l.Db != nil {
		err := l.Db.Put(key, value, nil)
		return err
	}

	return errors.New("not allowed ! db not initialzed")
}

func (l *Ledger) Get(key []byte) ([]byte, error) {
	if l.Db != nil {
		data, err := l.Db.Get(key, nil)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("not allowed ! db not initialzed")
}
