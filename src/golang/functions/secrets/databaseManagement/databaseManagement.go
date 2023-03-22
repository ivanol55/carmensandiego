// Sets the package name for the function package
package databaseManagement

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/helpers/errorManagement"

	"github.com/dgraph-io/badger/v3"
)

func InitializeDatabase() *badger.DB {
	// Initialize default Badger options
	var badgerOptions badger.Options = badger.DefaultOptions("")
	// Set in-memory storage for Badger
	badgerOptions = badgerOptions.WithInMemory(true)
	badgerOptions = badgerOptions.WithLoggingLevel(badger.WARNING)
	var database *badger.DB
	var err error
	database, err = badger.Open(badgerOptions)
	errorManagement.CheckError(err)
	return database
}

func WriteEntry(db *badger.DB, key string, value string) error {
	// Abstraction to write the requested entry into the selected database
	var err error
	var keyBytes []byte
	var valueBytes []byte
	keyBytes = []byte(key)
	valueBytes = []byte(value)
	err = db.Update(func(txn *badger.Txn) error {
		var dbEntry *badger.Entry = badger.NewEntry(keyBytes, valueBytes)
		var err error
		err = txn.SetEntry(dbEntry)
		return err
	})
	return err
}

func ReadEntry(db *badger.DB, key string) []byte {
	// Abstraction to read the requested entry from the selected database
	var err error
	var keyBytes []byte
	keyBytes = []byte(key)
	var fetchedValueCopy []byte
	err = db.View(func(txn *badger.Txn) error {
		var item *badger.Item
		var err error
		item, err = txn.Get(keyBytes)
		errorManagement.CheckError(err)
		err = item.Value(func(fetchedValue []byte) error {
			fetchedValueCopy = append(fetchedValueCopy, fetchedValue...)
			return nil
		})
		errorManagement.CheckError(err)
		return nil
	})
	errorManagement.CheckError(err)
	return fetchedValueCopy
}

func DeleteEntry(db *badger.DB, key string) error {
	// Abstraction to delete the requested entry from the selected database
	var err error
	var keyBytes []byte
	keyBytes = []byte(key)
	err = db.Update(func(txn *badger.Txn) error {
		var err error
		err = txn.Delete(keyBytes)
		return err
	})
	return err
}
