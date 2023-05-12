// Sets the package name for the main script
package scannerRunner

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/helpers/configManagement"
	"carmensandiego/src/golang/functions/helpers/threadManagement"
	"carmensandiego/src/golang/functions/secrets/databaseManagement"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

// Prepare everything necessary to run a secret scan
func RunScan(profileName string, initializedDatabases []*badger.DB) {
	// Prepare the databases we've set up so we can use them
	var patternsDatabase *badger.DB
	var filesDatabase *badger.DB
	var resultsDatabase *badger.DB
	patternsDatabase = initializedDatabases[0]
	filesDatabase = initializedDatabases[1]
	resultsDatabase = initializedDatabases[2]
	// Get regex rules based on the profile filtering
	var regexDictionary []databaseManagement.KeyValueEntry
	regexDictionary = getChosenRegexRules(profileName, patternsDatabase)
	scanDatabaseFiles(regexDictionary, filesDatabase, resultsDatabase)
}

func getChosenRegexRules(profileName string, patternsDatabase *badger.DB) []databaseManagement.KeyValueEntry {
	var profileRulesets []string
	var profile configManagement.Profile
	profile = configManagement.GetProfile(profileName)
	profileRulesets = profile.Ruleset
	var regexDBEntries []databaseManagement.KeyValueEntry
	if profileRulesets[0] == "all" && len(profileRulesets) == 1 {
		regexDBEntries = databaseManagement.FetchAllKeysAndValues(patternsDatabase)
	} else {
		var regexRule string
		for _, regexRule = range profileRulesets {
			var fetchedRegexDBEntries []databaseManagement.KeyValueEntry
			var fetchedRegexDBEntry databaseManagement.KeyValueEntry
			fetchedRegexDBEntries = databaseManagement.FetchPrefixedKeysAndValues(patternsDatabase, regexRule)
			for _, fetchedRegexDBEntry = range fetchedRegexDBEntries {
				regexDBEntries = append(regexDBEntries, fetchedRegexDBEntry)
			}
		}
	}
	return regexDBEntries
}

func scanDatabaseFiles(regexDictionary []databaseManagement.KeyValueEntry, filesDatabase *badger.DB, resultsDatabase *badger.DB) {
	// Prepare as many queues as we have routine threads for the profile
	var queues [][]string
	var profile configManagement.Profile
	var threads int
	threads = profile.Threads
	queues = threadManagement.GenerateQueues()
	// Separate all targets into the queues
	var fileList []string
	fileList = databaseManagement.FetchAllKeys(filesDatabase)
	queues = threadManagement.PopulateQueues(fileList, queues)
	fmt.Println(queues)
	// Run the scan as a parallel task per queue
	var secretScanningWaitGroup sync.WaitGroup
	secretScanningWaitGroup.Add(threads)
	var thread int
	for thread = 1; thread <= threads; thread = thread + 1 {
		go scanFilesForSecrets(&secretScanningWaitGroup, queues[thread-1], filesDatabase, regexDictionary, resultsDatabase)
	}
	secretScanningWaitGroup.Wait()
}

func scanFilesForSecrets(secretScanningWaitGroup *sync.WaitGroup, queue []string, filesDatabase *badger.DB, regexDictionary []databaseManagement.KeyValueEntry, resultsDatabase *badger.DB) {
	var filenameToScan string
	for _, filenameToScan = range queue {
		fmt.Println(filenameToScan)
	}
}

func showResults(resultsDatabase *badger.DB) {

}
