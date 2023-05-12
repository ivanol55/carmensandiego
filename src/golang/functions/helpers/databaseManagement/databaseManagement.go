// Sets the package name for the function package
package databaseManagement

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/helpers/errorManagement"

	"github.com/dgraph-io/badger/v3"
)

// Creates a key/value store struct for results fetched with FetchPrefixedKeysAndValues()
type KeyValueEntry struct {
	Key   string
	Value string
}

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

func FetchAllKeys(db *badger.DB) []string {
	// Abstraction to read all keys in the selected database
	var keysArray []string
	keysArray = []string{}
	var err error
	err = db.View(func(txn *badger.Txn) error {
		var iteratorOptions badger.IteratorOptions
		iteratorOptions = badger.DefaultIteratorOptions
		iteratorOptions.PrefetchValues = false
		var iterator *badger.Iterator
		iterator = txn.NewIterator(iteratorOptions)
		defer iterator.Close()
		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			var item *badger.Item
			item = iterator.Item()
			var keyBytes []byte
			keyBytes = item.Key()
			var keyString string
			keyString = string(keyBytes)
			keysArray = append(keysArray, keyString)
		}
		return nil
	})
	errorManagement.CheckError(err)
	return keysArray
}

func FetchAllKeysAndValues(db *badger.DB) []KeyValueEntry {
	var keyValueResults []KeyValueEntry
	var keyValueEntry KeyValueEntry
	keyValueEntry = KeyValueEntry{}
	var err error
	// Abstraction to read the requested keys and values from the database
	db.View(func(txn *badger.Txn) error {
		var iterator *badger.Iterator
		iterator = txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()
		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			var item *badger.Item
			item = iterator.Item()
			var keyBytes []byte
			keyBytes = item.Key()
			var keyString string
			keyString = string(keyBytes)
			err = item.Value(func(valueBytes []byte) error {
				var valueString string
				valueString = string(valueBytes)
				keyValueEntry = KeyValueEntry{keyString, valueString}
				keyValueResults = append(keyValueResults, keyValueEntry)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	errorManagement.CheckError(err)
	return keyValueResults
}

func FetchPrefixedKeys(db *badger.DB, prefix string) []string {
	// Abstraction to read the requested keys based on prefix from the database
	var keyBytesArray []string
	keyBytesArray = []string{}
	var err error
	db.View(func(txn *badger.Txn) error {
		var iterator *badger.Iterator
		iterator = txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()
		var prefixBytes []byte
		prefixBytes = []byte(prefix)
		for iterator.Seek(prefixBytes); iterator.ValidForPrefix(prefixBytes); iterator.Next() {
			var item *badger.Item
			item = iterator.Item()
			var keyBytes []byte
			keyBytes = item.Key()
			var keyString string
			keyString = string(keyBytes)
			keyBytesArray = append(keyBytesArray, keyString)
		}
		return nil
	})
	errorManagement.CheckError(err)
	return keyBytesArray
}

func FetchPrefixedKeysAndValues(db *badger.DB, prefix string) []KeyValueEntry {
	var keyValueResults []KeyValueEntry
	var keyValueEntry KeyValueEntry
	keyValueEntry = KeyValueEntry{}
	var err error
	// Abstraction to read the requested keys and values based on prefix from the database
	db.View(func(txn *badger.Txn) error {
		var iterator *badger.Iterator
		iterator = txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()
		var prefixBytes []byte
		prefixBytes = []byte(prefix)
		for iterator.Seek(prefixBytes); iterator.ValidForPrefix(prefixBytes); iterator.Next() {
			var item *badger.Item
			item = iterator.Item()
			var keyBytes []byte
			keyBytes = item.Key()
			var keyString string
			keyString = string(keyBytes)
			err = item.Value(func(valueBytes []byte) error {
				var valueString string
				valueString = string(valueBytes)
				keyValueEntry = KeyValueEntry{keyString, valueString}
				keyValueResults = append(keyValueResults, keyValueEntry)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	errorManagement.CheckError(err)
	return keyValueResults
}
