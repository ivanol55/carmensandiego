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
