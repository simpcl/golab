package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

func main() {
	opts := badger.DefaultOptions("./data")
	opts.Dir = "./data"
	opts.ValueDir = "./data"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// set
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("answer"), []byte("42"))
		return err
	})
	if err != nil {
		panic(err)
	}
	// get
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			return err
		}
		var val []byte = nil
		err = item.Value(func(v []byte) error {
			val = v
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Printf("The answer is: %s\n", val)
		return nil
	})
	if err != nil {
		panic(err)
	}
	// iterate
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			var val []byte
			err := item.Value(func(v []byte) error {
				val = v
				return nil
			})
			if err != nil {
				return err
			}
			fmt.Printf("key=%s, value=%s\n", k, val)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// Prefix scans
	db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("ans")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			var val []byte
			err := item.Value(func(v []byte) error {
				val = v
				return nil
			})
			if err != nil {
				return err
			}
			fmt.Printf("key=%s, value=%s\n", k, val)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// iterate keys
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Printf("key=%s\n", k)
		}
		return nil
	})
	// delete
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("answer"))
	})
	if err != nil {
		panic(err)
	}
}
