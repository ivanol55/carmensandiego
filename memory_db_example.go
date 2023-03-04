// Sets the package name for the main script
package main

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

func ifErrorPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main2() {
	// Initialize default Badger options
	var badgerOptions badger.Options = badger.DefaultOptions("")
	// Set in-memory storage for Badger
	badgerOptions = badgerOptions.WithInMemory(true)
	var err error
	db, err := badger.Open(badgerOptions)
	ifErrorPanic(err)
	defer db.Close()
	err = setBadgerEntry(db, "answer", "42")
	ifErrorPanic(err)
	var fetchedValue []byte
	fetchedValue = getBadgerEntry(db, "answer")
	fmt.Println(string(fetchedValue))
}

func setBadgerEntry(db *badger.DB, key string, value string) error {
	var err error
	var keyBytes []byte
	var valueBytes []byte
	keyBytes = []byte(key)
	valueBytes = []byte(value)

	err = db.Update(func(txn *badger.Txn) error {
		var dbEntry *badger.Entry = badger.NewEntry(keyBytes, valueBytes)
		err := txn.SetEntry(dbEntry)
		return err
	})
	return err
}
func getBadgerEntry(db *badger.DB, key string) []byte {
	var err error
	var keyBytes []byte
	keyBytes = []byte(key)
	var fetchedValueCopy []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(keyBytes)
		ifErrorPanic(err)
		err = item.Value(func(fetchedValue []byte) error {
			// Copying or parsing val is valid.
			fetchedValueCopy = append(fetchedValueCopy, fetchedValue...)
			return nil
		})
		ifErrorPanic(err)
		return nil
	})
	ifErrorPanic(err)
	return fetchedValueCopy
}
